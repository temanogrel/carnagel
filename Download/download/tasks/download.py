import os
import uuid

from celery.utils.log import get_task_logger
from common.aphrodite.performer import PerformerNotFoundError, MyFreeCamsPerformer, ChaturbatePerformer, \
    performer_api_factory
from common.aphrodite.recording import RecordingState, recording_api_factory
from common.minerva import FileType, minerva_factory

from download import config
from download.celery import app
from download.const import ModelNotAvailableError
from download.sites.chaturbate import download as record_chaturbate
from download.sites.myfreecams import download as record_myfreecams
from download.tasks.transcode import transcode

logger = get_task_logger(__name__)
performer_api = performer_api_factory('downloader')
recording_api = recording_api_factory('downloader')
minerva_api = minerva_factory(config.HOSTNAME, 'downloader')


class NoStorageServerAvailableError(RuntimeError):
    pass


class AlreadyRecordingWarning(RuntimeError):
    pass


@app.task(name='download.tasks.download_myfreecams',
          throws=(PerformerNotFoundError, AlreadyRecordingWarning, ModelNotAvailableError))
def download_myfreecams(performer_id: int, session_id: int, auto_transcode=True):
    """
    Initiate a recording of a myfreecams session

    :param performer_id:
    :param session_id:

    :return:
    """

    performer = performer_api.get(performer_id)

    if not isinstance(performer, MyFreeCamsPerformer):
        raise ValueError('Invalid performer, must be an instance of MyFreeCamsPerformer')

    if performer.get('isRecording', False):
        raise AlreadyRecordingWarning()

    performer_api.set_recording_state(performer, True)

    target = os.path.join(config.DOWNLOAD_PATH, 'cbc-{}.flv'.format(uuid.uuid4()))

    try:
        record_myfreecams(target, performer, int(session_id))

        if not os.path.exists(target):
            logger.info('No file was found post recording, most likely it was unavailable')
            return

        if os.path.getsize(target) < 20 * pow(1024, 2):
            logger.info('The file was smaller than the 20mb requirements we have')
            return

        # We need to create the recording to retrieve the id before uploading to minerva
        recording = recording_api.create_associated_recording(performer,
                                                              state=RecordingState.DOWNLOADED,
                                                              size264=os.path.getsize(target))

        recording.update(videoMp4Uuid=minerva_api.upload(target, recording['id'], FileType.RECORDING, dict()))
        recording_api.update(recording)

        # Just notify we are successful
        logger.info('Successfully recorded {} on myfreecams'.format(performer['stageName']))

        # Dispatch transcoding request
        transcode.delay(recording['id'])

    except Exception as e:
        logger.exception(e)

    finally:
        if os.path.exists(target):
            # file is not uploaded, we can delete it
            os.unlink(target)

        performer_api.set_recording_state(performer, False)


@app.task(name='download.tasks.download_chaturbate',
          throws=(PerformerNotFoundError, AlreadyRecordingWarning, ModelNotAvailableError))
def download_chaturbate(performer_id: int, username: str, api_token: str, auto_transcode=True):
    """
    Initiate a recording of a chaturbate session

    :param performer_id: The aphrodite performer id
    :param username: The chaturbate username
    :param api_token: The chaturbate api token associated to the username

    :return:
    """

    performer = performer_api.get(performer_id)

    if not isinstance(performer, ChaturbatePerformer):
        raise ValueError('Invalid performer instance, expected a chaturbate performer')

    if performer.get('isRecording', False):
        raise AlreadyRecordingWarning()

    performer_api.set_recording_state(performer, True)

    target = os.path.join(config.DOWNLOAD_PATH, 'cbc-{}.flv'.format(uuid.uuid4()))

    try:
        # Get the recording
        record_chaturbate(target, performer['stageName'], username, api_token)

        if not os.path.exists(target):
            logger.info('No file was found post recording, most likely it was unavailable')
            return

        if os.path.getsize(target) < 20 * pow(1024, 2):
            logger.info('The file was smaller than the 20mb requirements we have')
            return

        # We need to create the recording to retrieve the id before uploading to minerva
        recording = recording_api.create_associated_recording(performer,
                                                              state=RecordingState.DOWNLOADED,
                                                              size264=os.path.getsize(target))

        # upload the file to minerva and get a file uuid
        video_uuid = minerva_api.upload(target, recording['id'], FileType.RECORDING, dict())

        # Lets see what happens here
        logger.info('Uploaded file to minerva, uuid: {}'.format(video_uuid))

        # Update the recording
        recording.update(videoMp4Uuid=video_uuid)
        recording_api.update(recording)

        # Just notify we are successful
        logger.info('Successfully recorded {} on chaturbate'.format(performer['stageName']))

        # Dispatch transcoding request
        transcode.delay(recording['id'])

    except Exception as e:
        logger.exception(e)

    finally:
        # Cleanup
        if os.path.exists(target):
            os.unlink(target)

        performer_api.set_recording_state(performer, False)
