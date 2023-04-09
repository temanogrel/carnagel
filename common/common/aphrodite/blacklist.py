import re

from common.aphrodite.api import AbstractApiClient
from common.aphrodite.performer import AbstractPerformer
from common.consul import get_kv

IS_INT = re.compile(r'^\d+$')


class Blacklist:
    def __init__(self, stage_names=tuple(), service_ids=tuple()):
        self.service_ids = service_ids
        self.stage_names = stage_names

    def is_blacklisted(self, performer: AbstractPerformer):

        if not isinstance(performer, AbstractPerformer):
            raise ValueError('Performer is not a valid AbstractPerformer object')

        for alias in performer.get('aliases', []):
            if alias in self.stage_names:
                return True

        service_id = performer.get('serviceId', None)

        if service_id is not None and IS_INT.match(service_id):
            if int(service_id) in self.service_ids:
                return True

        return False


class BlacklistApi(AbstractApiClient):
    def get(self) -> Blacklist:
        """
        Get all the blacklisted service id's and aliases

        :return:
        """

        response = self._session.get(self._base_uri + '/blacklist')
        response.raise_for_status()

        stage_names = []
        service_ids = []

        for performer in response.json()['data']:

            if IS_INT.match(performer['serviceId']) and int(performer['serviceId']) not in service_ids:
                service_ids.append(int(performer['serviceId']))

            for alias in performer['aliases']:
                if alias not in stage_names:
                    stage_names.append(alias)

        return Blacklist(tuple(stage_names), tuple(service_ids))


def blacklist_api_factory(prefix: str):
    return BlacklistApi(
        base_uri=get_kv('{}/aphrodite/api'.format(prefix)),
        token=get_kv('{}/aphrodite/token'.format(prefix)),
    )
