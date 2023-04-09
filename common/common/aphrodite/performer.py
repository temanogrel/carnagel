from enum import Enum

from common.aphrodite.api import NotFoundError, AbstractApiClient
from common.consul import get_kv


class AbstractPerformer(dict):
    """
    Generic abstract performer
    """

    def __str__(self):
        return self.get('stageName')


class ChaturbatePerformer(AbstractPerformer):
    """
    Performer from the chaturbate service
    """
    pass


class MyFreeCamsPerformer(AbstractPerformer):
    """
    Performer from the myfreecams service
    """
    pass


class Cam4Performer(AbstractPerformer):
    """
    Performer from the Cam4 service
    """
    pass


class PerformerNotFoundError(NotFoundError):
    """
    Error for the response code 404 when trying to do something with a resource
    """
    pass


class PerformerServices(Enum):
    """
    List of services available
    """

    cbc = 'chaturbate'
    mfc = 'myfreecams'
    cam4 = 'cam4'


class PerformerApi(AbstractApiClient):
    """
    Performer api client
    """

    @property
    def resource_endpoint(self) -> str:
        return self._base_uri + '/performers/{performer_id}'

    @property
    def collection_endpoint(self) -> str:
        return self._base_uri + '/performers'

    @property
    def intersect_endpoint(self):
        return self._base_uri + '/performers/intersect'

    def _hydrate(self, performer: dict) -> AbstractPerformer:
        """
        Convert a python dict to a the proper performer object

        :param performer:
        :return:
        """
        if not isinstance(performer, dict):
            raise ValueError('Invalid performer object specified')

        if 'service' not in performer:
            raise ValueError('No service property was available to determine service')

        service = performer['service']

        if service == 'cbc':
            return ChaturbatePerformer(performer)
        elif service == 'mfc':
            return MyFreeCamsPerformer(performer)
        elif service == 'cam4':
            return Cam4Performer(performer)

        raise ValueError('Unsupported service performer service "{}" found.'.format(service))

    def intersect_online_performers(self, service: str, performers: dict):
        """
        Send all the online performers for the given service to the api

        :param service:
        :param performers:
        :return:
        """

        if not isinstance(performers, dict):
            raise ValueError('Performers should be a dict')

        structured = {}

        for performer in performers.values():
            structured[performer['serviceId']] = performer

        query = {
            'online': 1,
            'service': service
        }

        # Don't overwrite the recording / pending recording state from aphrodite in the intersect
        for performer in self.get_all(**query):

            # A "now" offline performer that no longer exists in the modelserver
            if performer['serviceId'] not in structured:
                continue

            structured[performer['serviceId']].update(
                isRecording=performer.get('isRecording'),
                isPendingRecording=performer.get('isPendingRecording')
            )

        response = self._session.post(self.intersect_endpoint, json=structured, params=dict(service=service))
        response.raise_for_status()

        for performer in response.json()['data']:
            structured[performer['serviceId']].update(**performer)

    def get_all(self, **kwargs):
        """
        Retrieve all performers matching the given criteria

        :param online:
        :param service:

        :return:
        """

        response = self._session.get(self.collection_endpoint, params=kwargs)
        response.raise_for_status()

        for performer in response.json()['data']:
            yield self._hydrate(performer)

    def get(self, performer_id, identifier=None) -> AbstractPerformer:
        """
        Retrieve a performer
        """

        response = self._session.get(self.resource_endpoint.format(performer_id=performer_id),
                                     params=dict(identifier=identifier))

        if response.status_code == 404:
            raise PerformerNotFoundError()

        # Raise on any other error
        response.raise_for_status()

        return self._hydrate(response.json())

    def update(self, performer: AbstractPerformer):
        """
        Update a performer

        :param performer:

        :return:
        """

        if not isinstance(performer, AbstractPerformer):
            raise ValueError('Expected a performer instance')

        if 'id' not in performer:
            raise KeyError('Missing performer id, unable to update')

        response = self._session.put(self.resource_endpoint.format(performer_id=performer['id']), json=performer)

        if response.status_code == 404:
            raise PerformerNotFoundError()

        response.raise_for_status()

        # Update the performer with the response data
        performer.update(**response.json())

    def set_recording_state(self, performer: AbstractPerformer, state: bool):
        """
        Change the recording state of a performer, simple utility function.

        :param performer:
        :param state:

        :return:
        """

        if not isinstance(performer, AbstractPerformer):
            raise ValueError('Expected a performer instance')

        # Update the performer
        performer.update(isRecording=state)

        # Do a regular update
        self.update(performer)


def performer_api_factory(prefix: str):
    return PerformerApi(
        base_uri=get_kv('{}/aphrodite/api'.format(prefix)),
        token=get_kv('{}/aphrodite/token'.format(prefix)),
    )
