import os
import requests.adapters

from common.aphrodite.performer import AbstractPerformer
from common.aphrodite.recording import AbstractRecording
from common.consul import get_kv

session = requests.session()
session.mount('http://', requests.adapters.HTTPAdapter(max_retries=3))


class UltronApiError(RuntimeError):
    pass


class UltronAuthorizationError(UltronApiError):
    pass


class UltronRecordingNotFoundError(UltronApiError):
    pass


class UltronRecordingAlreadyPublishedError(UltronApiError):
    pass


class UltronApi:
    """
    Api to communicate with the ultron system
    """

    def __init__(self, base_uri: str, token: str):
        if token is None:
            raise ValueError('Token is not provided')

        if base_uri is None:
            raise ValueError('Base uri is not provided')

        self.token = token
        self.base_uri = base_uri

    @property
    def collection_endpoint(self):
        return self.base_uri + '/api/recordings'

    @property
    def resource_endpoint(self):
        return self.base_uri + '/api/recordings/{recording_id}'

    def _get_headers(self):
        return {
            'Authorization': self.token
        }

    def create(self, performer: AbstractPerformer, recording: AbstractRecording):
        """
        Create the recording on ultron

        :param performer:
        :param recording:
        :return:
        """

        data = dict(performer=performer, recording=recording)

        response = session.post(self.collection_endpoint, json=data, headers=self._get_headers())

        if response.status_code == 401:
            raise UltronAuthorizationError()

        elif response.status_code == 409:
            raise UltronRecordingAlreadyPublishedError

    def delete(self, recording: AbstractRecording):
        """
        Remove a recording from ultron

        :param recording:
        :return:
        """

        url = self.resource_endpoint.format(recording_id=recording['id'])

        response = session.delete(url, headers=self._get_headers())

        if response.status_code == 401:
            raise UltronAuthorizationError()

        elif response.status_code == 404:
            raise UltronRecordingNotFoundError

        response.raise_for_status()


def ultron_factory(prefix: str) -> UltronApi:
    return UltronApi(
        token=get_kv('{}/ultron/token'.format(prefix)),
        base_uri=get_kv('{}/ultron/api'.format(prefix)),
    )
