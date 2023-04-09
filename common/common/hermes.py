import base64
import re

import requests.adapters

from common.consul import get_kv

session = requests.Session()
session.mount('http://', requests.adapters.HTTPAdapter(max_retries=3))

KEY_DETECTION = re.compile(r'^http://(?P<host>(.*))/(?P<key>[a-zA-Z0-9\-]+)$')


class HermesUrl(dict):
    """
    Data object class presenting the hermes urls
    """

    def generate_hermes_url(self) -> str:
        return 'http://{}/{}'.format(self['hostname'], self['key'])


class HermesUrlNotFound(RuntimeError):
    pass


class HermesClient:
    """
    A client for our short url service
    """

    def __init__(self, token: str, hostname: str):
        self.host = hostname
        self.api_token = token

    @property
    def collection_endpoint(self):
        return self.host + '/api/url'

    @property
    def resource_endpoint(self):
        return self.host + '/api/url/{key}'

    def _get_headers(self) -> dict:
        """
        Retrieve the headers

        :return: dict
        """

        return {
            'Authorization': self.api_token
        }

    def create(self, url: str) -> str:
        """
        Create a new short url and return it

        :param url:

        :return:
        """

        response = session.post(self.collection_endpoint, data=dict(url=url), headers=self._get_headers())
        response.raise_for_status()

        data = response.json()

        return data['shortUrl']

    def get(self, url: str):
        """
        :param urL:
        :return:
        """

        matches = KEY_DETECTION.search(url)

        if not matches or matches.group('host') not in ('cur.bz', 'pip.bz'):
            raise ValueError('Unsupported video url, expected a cur.bz or pip.bz host origin')

        # Build url based upon the host
        uri = 'http://{host}/api/url/{key}'.format(host=matches.group('host'), key=matches.group('key'))

        response = session.get(uri, headers=self._get_headers())
        if response.status_code == 404:
            raise HermesUrlNotFound()

        response.raise_for_status()

        return HermesUrl(response.json())

    def get_by_original_url(self, url: str):
        """
        Retrieve a url by it's original url instead of the short key

        :param url:

        :return:
        """

        if isinstance(url, str):
            url = bytes(url, 'utf-8')

        key = base64.b64encode(url).decode('utf-8')

        response = session.get(
            self.resource_endpoint.format(key=key),
            headers=self._get_headers(),
            params=dict(identifier='originalUrl')
        )

        if response.status_code == 404:
            raise HermesUrlNotFound()

        response.raise_for_status()

        return HermesUrl(response.json())

    def replace_url(self, old_url: str, new_url: str):
        """
        Update the original url for a existing key.

        Works by passing out the last string of the old_uri and using that as a key when updating

        :param old_url:
        :param new_url:

        :return:
        """

        matches = KEY_DETECTION.search(old_url)

        if not matches or matches.group('host') not in ('cur.bz', 'pip.bz'):
            raise ValueError('Unsupported video url, expected a cur.bz or pip.bz host origin')

        data = dict(url=new_url)

        # Build the uri based on the used host
        uri = 'http://{host}/api/url/{key}'.format(host=matches.group('host'), key=matches.group('key'))

        response = session.patch(uri, data=data, headers=self._get_headers())
        response.raise_for_status()


def hermes_factory(prefix: str):
    return HermesClient(
        token=get_kv('{}/hermes/token'.format(prefix)),
        hostname=get_kv('{}/hermes/api'.format(prefix)),
    )
