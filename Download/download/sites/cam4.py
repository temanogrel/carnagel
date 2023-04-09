import logging

import requests
from celery.utils.log import get_task_logger
from librtmp.amf import AMFObject

from download.const import ModelNotAvailableError, RecordingTerminatedError, TERMINATION_STATUSES, RecordingState
from download.rtmp import RTMPClient
from download.utils import create_normalized_filename, is_valid_recording

logging.basicConfig(level=logging.INFO)
logger = get_task_logger(__name__)

def get_rtmp_server(stage_name: str) -> list:
    """
    Parse the flash vars from the page

    This is very likely to fail as soon as they modify the page, but not much we can do about that.

    @todo: Add a check for the availability of the model.

    :param stage_name:

    :return:
    """

    settings = requests.get('http://webchat.cam4.com:8080/requestRoomInformation?roomname={stage_name}&failedVideoUrl='.format(stage_name=stage_name))
    settings.raise_for_status()

    data = settings.json()

    if data['status'] != 'success':
        raise ModelNotAvailableError('Performer not available, error: {}'.format(data['status']))

    return data['rtmpUrl']


def client_factory(server, stage_name: str) -> RTMPClient:
    """
    Create the RTMP client and configure it with all the required callbacks

    :param model_name:
    :param stage_name:
    :param api_token:

    :return:
    """

    data = (stage_name, 'guest', '')

    client = RTMPClient(server, connect_data=data)

    @client.invoke_handler('updateChallenge')
    def set_status(status: AMFObject):
        pass

    @client.invoke_handler('onStatus')
    def on_status(status_data: AMFObject) -> None:
        """
        This method gets notified about events that are relevant to the stream

        :param status_data:

        :return:
        """

        logger.info('Status update: {}'.format(status_data['description']), extra=status_data)

        if status_data['code'] in TERMINATION_STATUSES:
            client.close()

    return client


def download(stage_name: str, section: str) -> str:
    """
    Download a video from the chaturbate service

    Will return a string with the folder & file name if the recording is successful, else it will return None

    :param stage_name: string
    :param section: String
    :param stage_name: String the stage_name associated to the api token
    :param api_token: String api token to allow access

    :return:
    """

    try:
        server = get_rtmp_server(stage_name)

        client = client_factory(server, stage_name)
        client.connect()
        client.process_packets(transaction_id=1.0)

        if not client.connected:
            return logger.warning('Failed to connect', extra=dict(state=RecordingState.ERROR.value))

        # This will tell the stream that we are awaiting to view the stream
        result = client.call('receiveRTMPResponse')

        print(result.result())

        # This will initiate the remote server to start sending us the data
        stream = client.call('createStream')

        # block until we get the response and validate it
        if stream.result() != 1.0:
            raise RecordingTerminatedError('Failed to create a valid video stream')

        # Generate a normalized filename for the chaturbate service
        target = create_normalized_filename(stage_name, 'cam4', section)

        try:
            # Start recording
            client.record(target)

        except Exception as e:
            logger.exception(e)

        finally:
            return target if is_valid_recording(target) else None

    except Exception as e:
        print(e)
        logger.exception(e)

    finally:
        # Close only if the client is available
        if 'client' in locals():
            client.close()
