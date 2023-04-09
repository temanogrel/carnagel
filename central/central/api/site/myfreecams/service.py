import consul

from central.api.site.myfreecams.const import VideoStates, AccessLevel
from central.api.site.service import AbstractService
from central.tasks.download import download_mfc_stream
from common.aphrodite.performer import MyFreeCamsPerformer, AbstractPerformer, performer_api_factory

# Unsupported states
IGNORED_STATES = [
    VideoStates.TX_PVT,
    VideoStates.TX_GRP,
    VideoStates.TX_AWAY,
    VideoStates.TX_PVT,
    VideoStates.RX_IDLE,
    VideoStates.OFFLINE
]

# Aphrodite api's
performer_api = performer_api_factory('central')
kv = consul.Consul().kv


class MyFreeCamsService(AbstractService):

    def __init__(self, log):
        super().__init__(log)
        self._session_id = None

    @property
    def session_id(self):
        return self._session_id

    def get_required_viewer_count(self) -> int:
        index, data = kv.get('central/myfreecams/minimum_viewers')
        if data is None:
            return 15

        return int(data['Value'])

    def may_record(self, performer) -> bool:
        if not isinstance(performer, MyFreeCamsPerformer):
            return False

        if not self.session_id:
            return False

        if performer['videoState'] in IGNORED_STATES:
            return False

        if performer['accessLevel'] != AccessLevel.MODEL.value:
            return False

        if performer['camServer'] == 0:
            return False

        return True

    def dispatch_recording(self, performer: MyFreeCamsPerformer) -> None:
        if not isinstance(performer, MyFreeCamsPerformer):
            raise ValueError('Expected an instance of myfreecams')

        download_mfc_stream.delay(int(performer['id']), self.session_id, auto_transcode=True)

    def bootstrap(self):
        for performer in performer_api.get_all(online=1, service='mfc'):
            self[performer['stageName']] = performer

    @staticmethod
    def hydrate(data: dict) -> AbstractPerformer:
        return MyFreeCamsPerformer(data)
