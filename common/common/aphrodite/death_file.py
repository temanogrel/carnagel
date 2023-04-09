import base64
from enum import Enum
from common.aphrodite.api import NotFoundError, AbstractApiClient, CollectionResult
from common.consul import get_kv


class UrlEntryState(str, Enum):
    IGNORED = 'ignored',
    REMOVED = 'removed',
    PENDING = 'pending',
    IN_PROGRESS = 'in-progress'


class DeathFile(dict):
    def get_file_processing_offset(self) -> int:
        """
        When starting to process a death file using the ignored/pending offset to set a
        ruff position from where to start again
        :return:
        """

        return self.get('pending', 0) + self.get('ignored', 0)


class UrlEntry(dict):
    pass


class UrlEntryNotFoundError(NotFoundError):
    pass


class DeathFileNotFoundError(NotFoundError):
    pass


class DeathFileApi(AbstractApiClient):
    @property
    def resource_endpoint(self) -> str:
        return self._base_uri + '/death-files/{file_id}'

    def get(self, file_id: int) -> DeathFile:
        """
        Retrieve a death file

        :param file_id:
        """

        response = self._session.get(self.resource_endpoint.format(file_id=file_id))

        if response.status_code == 404:
            raise DeathFileNotFoundError()

        # Raise on any other error
        response.raise_for_status()

        return DeathFile(response.json())

    def update(self, file: DeathFile):
        """
        Update a death file

        :param file:

        :return:
        """

        if not isinstance(file, DeathFile):
            raise ValueError('Expected a DeathFile instance')

        if 'id' not in file:
            raise KeyError('Missing file id, unable to update')

        response = self._session.put(self.resource_endpoint.format(file_id=file['id']), json=file)

        if response.status_code == 404:
            raise DeathFileNotFoundError()

        response.raise_for_status()

        # Update the performer with the response data
        file.update(**response.json())


class UrlEntryApi(AbstractApiClient):
    @property
    def resource_endpoint(self):
        return self._base_uri + '/urls/{url_id}'

    @property
    def collection_endpoint(self):
        return self._base_uri + '/urls'

    @property
    def restricted_collection_endpoint(self):
        return self._base_uri + '/death-files/{file_id}/urls'

    def get(self, url_id: str, identifier='id'):

        if identifier == 'url':
            if isinstance(url_id, str):
                url_id = bytes(url_id, 'utf-8')

            url_id = base64.b64encode(url_id).decode('utf-8')

        response = self._session.get(self.resource_endpoint.format(url_id=url_id), params=dict(identifier=identifier))

        if response.status_code == 404:
            raise UrlEntryNotFoundError()

        response.raise_for_status()

        return UrlEntry(response.json())

    def get_all(self, **kwargs):

        response = self._session.get(self.collection_endpoint, params=kwargs)
        response.raise_for_status()

        result = response.json()

        items = tuple(UrlEntry(entry) for entry in result['data'])
        total = result['meta']['total']
        offset = result['meta']['offset']

        return CollectionResult(items, total, offset)

    def update(self, entry: UrlEntry):

        if not isinstance(entry, UrlEntry):
            raise ValueError('Invalid entry provided')

        if 'id' not in entry:
            raise KeyError('Missing entry identifier')

        response = self._session.put(self.resource_endpoint.format(url_id=entry['id']), json=entry)

        if response.status_code == 404:
            raise UrlEntryNotFoundError()

        response.raise_for_status()

        entry.update(**response.json())

    def create(self, file: DeathFile, entry: UrlEntry):

        if not isinstance(file, DeathFile):
            raise ValueError('Invalid file provided')

        if not isinstance(entry, UrlEntry):
            raise ValueError('Invalid entry provided')

        response = self._session.post(self.restricted_collection_endpoint.format(file_id=file['id']), json=entry)
        if response.status_code == 404:
            raise DeathFileNotFoundError()

        response.raise_for_status()

        return UrlEntry(response.json())


def deathfile_api_factory(prefix: str):
    return DeathFileApi(
        base_uri=get_kv('{}/aphrodite/api'.format(prefix)),
        token=get_kv('{}/aphrodite/token'.format(prefix)),
    )


def urlentry_api_factory(prefix: str):
    return UrlEntryApi(
        base_uri=get_kv('{}/aphrodite/api'.format(prefix)),
        token=get_kv('{}/aphrodite/token'.format(prefix)),
    )
