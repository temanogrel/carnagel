from common.aphrodite.api import AbstractApiClient, NotFoundError
from common.consul import get_kv


class AbstractSite(dict):
    def __str__(self):
        return self.get('name')


class Site(dict):
    def may_publish_recording(self, recording) -> bool:
        """
        Check if the recording matches the requirements to post on the given site

        :param recording:
        :return:
        """

        # This is ugly as fuck, but i don't know how else to resolve it.
        from common.aphrodite.recording import UnassociatedRecording, AbstractRecording

        if not isinstance(recording, AbstractRecording):
            raise ValueError('Expected a recording instance')

        sources = self['sources']

        if isinstance(recording, UnassociatedRecording):
            return sources.get('unassociated', False)

        service = recording.get('service')

        if service not in sources:
            return False

        # Can either be bool or a dict
        tmp = sources[service]

        # if it's a bool it means this service does not support sections
        # so we can short circuit with it's value
        if isinstance(tmp, bool):
            return tmp

        section = recording.get('section')

        if section not in tmp or not tmp[section]:
            return False

        return True


class SiteNotFoundError(NotFoundError):
    pass


class SiteApi(AbstractApiClient):
    @property
    def resource_endpoint(self) -> str:
        return self._base_uri + '/sites/{id}'

    @property
    def collection_endpoint(self) -> str:
        return self._base_uri + '/sites'

    def get_all(self, **kwargs):
        response = self._session.get(self.collection_endpoint, params=kwargs)
        response.raise_for_status()

        for site in response.json()['data']:
            yield Site(**site)


def site_api_factory(prefix: str):
    return SiteApi(
        base_uri=get_kv('{}/aphrodite/api'.format(prefix)),
        token=get_kv('{}/aphrodite/token'.format(prefix)),
    )
