from celery.utils.log import get_task_logger

from central.celery import app
from central.tasks.storage import reupload_recording_image, reupload_recording_video
from common.aphrodite.performer import performer_api_factory
from common.aphrodite.recording import RecordingNotFoundError, RecordingState, AbstractRecording, \
    recording_api_factory
from common.aphrodite.site import site_api_factory
from common.ultron import UltronRecordingAlreadyPublishedError, ultron_factory
from common.wordpress import WPPost, create_post_content, create_post_title, WPApi

logger = get_task_logger(__name__)

site_api = site_api_factory('central')
ultron_api = ultron_factory('central')
performer_api = performer_api_factory('central')
recording_api = recording_api_factory('central')


@app.task(name='central.tasks.publish_recording', throws=(RecordingNotFoundError,))
def publish_recording(recording_id: int, wordpress=True, ultron=True, infinity=True):
    recording = recording_api.get(recording_id)
    recording_api.set_recording_state(recording, RecordingState.PUBLISHING)

    try:
        if any(image is None for image in recording['imageUrls'].values()):
            reupload_recording_image.apply_async(args=(recording['id'],), routing_key=recording['storageServer'])
            return logger.warning('Missing uploaded images for recording: {}'.format(recording['id']))

        if recording['videoUrl'] is None:
            reupload_recording_video.apply_async(args=(recording['id'],), routing_key=recording['storageServer'])
            return logger.warning('Missing uploaded video url for recording: {}'.format(recording['id']))

        if wordpress:
            publish_to_wordpress(recording)

        if ultron:
            publish_to_ultron(recording)

        if infinity:
            publish_to_infinity(recording)

        recording_api.set_recording_state(recording, RecordingState.PUBLISHED)

    except Exception as e:

        # Update state
        recording_api.set_recording_state(recording, RecordingState.PUBLISHING_FAILED)

        # re-raise exception, makes celery notify the task as failed.
        raise e


def publish_to_wordpress(recording: AbstractRecording):
    """
    Publishes to all wordpress sites tht are currently enabled

    :param recording:

    :return:
    """

    sites = site_api.get_all(enabled=1, service=recording['service'])

    for site in sites:

        if recording.is_published_on(site):
            continue

        if not site.may_publish_recording(recording):
            continue

        post = WPPost({
            'title': create_post_title(recording),
            'content_raw': create_post_content(recording),
            'status': 'publish',
            'comment_status': 'closed'
        })

        wp_api = WPApi(site['apiUri'], site['username'], site['password'])
        wp_api.create(post)

        recording_api.add_post_association(recording, site, post)


def publish_to_ultron(recording: AbstractRecording):
    """
    Publish a recording to ultron

    :param recording:
    :return:
    """

    performer = performer_api.get(recording['performerId'])

    try:
        ultron_api.create(performer, recording)
    except UltronRecordingAlreadyPublishedError:
        logger.warning('Recording "{}" has already been published'.format(recording['id']))


def publish_to_infinity(recording: AbstractRecording):
    """
    Publish the recording to the infinity system

    :param recording:
    :return:
    """
    #
    # performer = performer_api.get(recording['performerId'])
    #
    # request = CreateRecordingRequest(
    #     performer=InfinityPerformer(
    #         id=performer['id'],
    #         serviceId=performer['serviceId'],
    #         stageName=performer['stageName'],
    #         service=performer['service'],
    #         aliases=performer['aliases'],
    #     ),
    #     recording=InfinityRecording(
    #         id=recording['id'],
    #         performerID=performer['id'],
    #         type=recording['type'],
    #         stageName=recording['stageName'],
    #         section=recording['section'],
    #         service=recording['service'],
    #         duration=recording['duration'],
    #         size=recording['size264'],
    #         description=recording['description'],
    #         location=FileLocation(hostname=recording['storageServer'], path=recording['storagePath']),
    #         storagePathCollage=recording['storagePathCollage'],
    #         bitRate=recording['bitRate'],
    #         audio=recording['audio'],
    #         video=recording['video'],
    #         images=images,
    #         createdAt=recording['createdAt'],
    #         updatedAt=recording['updatedAt']
    #     )
    # )
    #
    # try:
    #
    #     # Setup infinity grpc connection
    #     channel = grpc.insecure_channel(config.INFINITY_API)
    #
    #     infinity = RecordingServiceStub(channel=channel)
    #     infinity.Create(request)
    #
    #     # Not sure if this will properly close it
    #     del channel
    #
    # except Exception as e:
    #     logger.info('Failed to published to infinity')
    #     logger.exception(e)
    pass

