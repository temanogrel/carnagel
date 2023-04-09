import random

import consul

from central.api.site.service import AbstractService
from central.tasks.download import download_cbc_stream
from common.aphrodite.performer import ChaturbatePerformer, AbstractPerformer, performer_api_factory

performer_api = performer_api_factory('central')
kv = consul.Consul().kv


class ChaturbateService(AbstractService):
    def __init__(self, log):
        super().__init__(log)

        self.credentials = {}

    def set_credential(self, identity: str, api_token: str):
        """
        Set the api token to use with a email

        :param identity: str
        :param api_token: str

        :return: None
        """
        self.credentials[identity] = api_token

    def get_credentials(self):
        return self.credentials

    def get_required_viewer_count(self) -> int:
        index, data = kv.get('central/chaturbate/minimum_viewers')
        if data is None:
            return 15

        return int(data['Value'])

    def may_record(self, performer: AbstractPerformer) -> bool:
        if not isinstance(performer, ChaturbatePerformer):
            return False

        return True

    def dispatch_recording(self, performer: ChaturbatePerformer):

        # Get a random dict key
        username = random.sample(self.credentials.keys(), 1).pop()

        # Get the associated api token
        api_token = self.credentials[username]

        # start recording
        download_cbc_stream.delay(performer['id'], username, api_token, auto_transcode=True)

    def bootstrap(self):
        for performer in performer_api.get_all(online=1, service='cbc'):
            self[performer['stageName']] = performer

    @staticmethod
    def hydrate(data: dict) -> AbstractPerformer:
        return ChaturbatePerformer(data)
