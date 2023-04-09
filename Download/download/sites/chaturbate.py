import re

import requests
from celery.utils.log import get_task_logger
from librtmp.amf import AMFObject

from download.const import ModelNotAvailableError, RecordingTerminatedError, TERMINATION_STATUSES
from download.rtmp import RTMPClient

logger = get_task_logger(__name__)


def get_vars_from_page(username: str) -> list:
    """
    Parse the flash vars from the page

    This is very likely to fail as soon as they modify the page, but not much we can do about that.

    @todo: Add a check for the availability of the model.

    :param username:

    :return:
    """

    page = requests.get('https://chaturbate.com/{username}'.format(username=username))
    page.raise_for_status()

    match = re.compile(r'EmbedViewerSwf\((?P<data>(.|\n)*)\n\s+\);').search(page.text)

    if not match:
        if re.compile(r'EmbedViewerSwf').search(page.text):
            logger.warning("Found 'EmbedViewerSwf' in page content but was unable to match it against regex")

        raise ModelNotAvailableError('Failed to find the javascript code to embed the viewer')

    lines = match.group('data').split('",')

    # Remove excess data
    lines = [line.strip() for line in lines]

    # Remove padding
    lines = [line[1:] for line in lines]

    logger.info('Retrieved chaturbate SWF parameters', extra=dict(parameters=lines))

    return lines


def client_factory(model_name: str, username: str, api_token: str) -> RTMPClient:
    """
    Create the RTMP client and configure it with all the required callbacks

    :param model_name:
    :param username:
    :param api_token:

    :return:
    """

    # EmbedViewerSwf(swfname,modelname,fms_address,mute,sa,uid,ws,pw,rp,tv,cv,headless,auth,flash_debug)
    (
        swf,
        name,
        fms_address,
        mute,
        sa,
        uid,
        ws,
        pw,
        rp,
        tv,
        cv,
        headless,
        auth,
        flash_debug
    ) = get_vars_from_page(model_name)
    
    auth_token = auth.replace("\\u0022", "\"")

    data = (
        'AnonymousUser',
        name,
        '2.690',
        pw,
        rp,
        auth_token
    )

    logger.info("Creating rtmp client with connect data: {}".format(data))
    client = RTMPClient('rtmp://{}'.format(fms_address), app='live-edge', connect_data=data)

    @client.invoke_handler('cpsQuality')
    def cpsQuality(success: bool, code, stream_id: str, server: str, quality: str) -> None:
        """
        This method gets executed when the client receive a cps method.
        What this method does is tell the client what server and stream to user

        :param success:
        :param code:
        :param stream_id:
        :param server:

        :return:
        """
        log_args = dict(
            success=success,
            code=code,
            stream_id=stream_id.replace('\n', ' '),
            server=server.replace('\n', ' ')
        )

        logger.info('CheckPublicStatus, success: {}, code: {}, stream: {}, server: {}'.format(
            *log_args.values()), extra=log_args)

        if code == 'hidden':
            raise ModelNotAvailableError('Performer is in a hidden show')

        if not success:
            raise ModelNotAvailableError('CheckPublicStatus check returned false')

        client.set_option('playpath', 'mp4:rtmp://{}/live-origin/{}'.format(server, stream_id))

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


def download(target: str, stage_name: str, account_username: str, account_token: str):
    """
    Download a video from the chaturbate service

    Will return a string with the folder & file name if the recording is successful, else it will return None

    :type target: The path we want to download to
    :type stage_name: String the stage name of the performer
    :type account_username: String the username associated to the api token
    :type account_token: String api token to allow access
    """

    client = client_factory(stage_name, account_username, account_token)
    client.connect()
    client.process_packets(transaction_id=1.0, debug=True)

    if not client.connected:
        raise RecordingTerminatedError('Failed to connect')

    # This will tell the stream that we are awaiting to view the stream
    client.call('CheckPublicStatus')

    # pause until the cps has been invoked because that will tell
    # the client what play-path to use
    client.process_packets(invoked_method='cpsQuality', debug=True)

    # This will initiate the remote server to start sending us the data
    stream = client.call('createStream')

    # block until we get the response and validate it
    if stream.result() != 1.0:
        raise RecordingTerminatedError('Failed to create a valid video stream')

    # Start recording
    client.record(target)
