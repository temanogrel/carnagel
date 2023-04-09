from central.api.site.service import AbstractService
from common.aphrodite.performer import Cam4Performer


class Cam4Service(AbstractService):

    @staticmethod
    def hydrate(data: dict) -> Cam4Performer:
        return Cam4Performer(data)

    def may_record(self, performer: Cam4Performer) -> bool:
        return False

    def dispatch_recording(self, performer: Cam4Performer) -> None:
        pass

    def bootstrap(self):
        performer_api = performer_api_factory('central')

        for performer in performer_api.get_all(online=1, service='cam4'):
            self[performer['stageName']] = performer
