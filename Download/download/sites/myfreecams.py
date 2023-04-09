import math

import execjs
from celery.utils.log import get_task_logger
from common.aphrodite.performer import MyFreeCamsPerformer
from librtmp import RTMPError
from librtmp.amf import AMFObject
from consul import Consul
from random import choice

from download.const import *
from download.rtmp import RTMPClient

FCTYPE_LOGIN = 1
FCTYPE_METRICS = 69

# RTMP Amf data
MODE_DOWNLOAD = 'DOWNLOAD'
MODE_CAM2CAM = 'CAM2CAM'

# Model stats
MODEL_AWAY = 2
MODEL_OFFLINE = 1

SERVER_MAP = {
    104: 64,
    109: 69,
    114: 74,
    149: 112,
    153: 116,
    157: 120,
    161: 124,
    165: 128,
    169: 132,
    181: 144,
    185: 148,
    199: 162,
    207: 170,
    211: 174,
    215: 178,
    237: 200,
    241: 204,
    245: 208,
    249: 212,
    253: 216,
    257: 220
}

logger = get_task_logger(__name__)


def get_real_server(server: int) -> int:
    """
    The server number provided by the websocket server is not 100% accurate and this fix
    is actually copied from their javascript source code.

    :param server:
    :return:
    """

    server = int(server)

    if server == 0:
        raise ModelNotAvailableError('Invalid cam server 0')

    if server > 500:
        server -= 500

    if server in SERVER_MAP:
        return SERVER_MAP[server]

    return server


def create_rtmp_client(server: int, model_id: int, session_id: int) -> RTMPClient:
    """
    Create the rtmp client and configure it with all the required callback for mfc to function properly

    :param server:
    :param model_id:
    :param session_id:

    :return:
    """

    # Get proxies from consul
    index = None
    index, data = Consul().kv.get('proxy-servers/ips', index=index)
    proxies = data['Value'].decode('utf-8').split(',')
    proxy = choice(proxies)
    proxy = proxy.split(':')[0]
    proxy = proxy + ':1080'

    # Define variables
    room_id = model_id + 100000000
    password = ''
    play_path = 'mp4:mfc_{}.f4v'.format(room_id)
    flash_version = 'MAC 15,0,0,242'
    server = 'http://video{}.myfreecams.com:1935/NxServer'.format(get_real_server(server))
    params = (
        session_id,
        password,
        room_id,
        MODE_DOWNLOAD,
        model_id
    )

    logger.info('Connecting to RTMP server: {} via proxy: {} to download: {}'.format(server, proxy, play_path))

    # Create the RTMP client but don't connect
    client = RTMPClient(
        server,
        app='NxServer',
        playpath=play_path,
        flashver=flash_version,
        connect_data=params,
        # from librtmp docs:
        # These options define how to connect to the media server.
        # socks=host:port
        # Use the specified SOCKS4 proxy.
        socks=proxy
    )

    @client.invoke_handler('loginResult')
    def login_challenge(data: AMFObject) -> str:
        """
        The remote flash server will send us a challenge that is basically a javascript function
        that we need to execute and return the results of

        :return
        """

        if type(data['challenge']) is str:
            response = execjs.eval(data['challenge'])
        else:
            response = str(math.floor(math.sqrt(data['challenge']) * 4))

        return response

    @client.invoke_handler('UpdateSession')
    def update_session(data: AMFObject):
        """
        Get notifications about session updates

        """
        logger.info('Session update: {}'.format(data))

        if int(data['cmdarg2']) != model_id:
            raise ModelNotAvailableError('Model id miss-match')

        current_cam_mode = int(data['cmdarg3'])

        if current_cam_mode == 13:
            raise WeAreBannedError()

        if current_cam_mode in [MODEL_AWAY, MODEL_OFFLINE]:
            raise ModelNotAvailableError('Model is currently offline or away')

    @client.invoke_handler('onStatus')
    def on_status(status_data: AMFObject):
        """
        This method gets notified about events that are relevant to the stream

        :param data:

        :return:
        """

        logger.info('Status update: {}'.format(status_data['description']), extra=status_data)

        if status_data['code'] in TERMINATION_STATUSES:
            client.close()

    @client.invoke_handler('authFailure')
    def authentication_failure(data: AMFObject):
        raise AuthenticationFailureError()

    return client


EXPECTED_EXCEPTIONS = (
    ModelNotAvailableError, AuthenticationFailureError, WeAreBannedError, RTMPError,
    RecordingTimeoutError, RecordingTerminatedError
)


def download(target: str, performer: MyFreeCamsPerformer, session_id: int):
    """
    Download a stream from a MFC server

    :param target: 
    :param performer: MyFreeCamsPerformer
    :param session_id: int

    :return:
    """

    server = int(performer['camServer'])
    model_id = int(performer['serviceId'])

    # Setup the client
    client = create_rtmp_client(server, model_id, session_id)
    client.connect()

    if not client.connected:
        raise RecordingTerminatedError('Failed to connect')

    # If we successfully connected then we need to wait for the update session to be called
    # so that we know everything is ok.
    client.process_packets(invoked_method='UpdateSession')

    # This will create a new stream as well as execute the play command
    # The play command value is the name as the playpath in the librtmp.RTMP constructor
    stream = client.call('createStream')

    # Failed to create the stream
    if stream.result() != 1.0:
        raise ModelNotAvailableError('Failed to create a stream')

    # Non-standard calls to start the stream but are sent by their flash application
    # so we might as well try and mimic it as much as possible.
    client.call('startDownload')
    client.call('receiveAudio')

    # Start recording
    client.record(target)
