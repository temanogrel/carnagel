import base64
from enum import Enum

import common.wordpress
from common.aphrodite.api import AbstractApiClient, ApiError, CollectionResult
from common.aphrodite.performer import AbstractPerformer, PerformerNotFoundError
from common.aphrodite.site import Site
from common.consul import get_kv


class RecordingNotFoundError(ApiError):
    pass


class RecordingOrphanedError(RuntimeError):
    def __init__(self, recording_id: int):
        self.recording_id = recording_id


class MultipleRecordingsFoundError(ApiError):
    def __init__(self, identifiers=()):
        self.identifiers = identifiers


class RecordingState(str, Enum):
    # Downloading server
    DOWNLOADED = 'downloaded'

    # Transcoding server
    TRANSCODING = 'transcoding'
    TRANSCODED = 'transcoded'
    TRANSCODING_FAILED = 'transcoding_failed'
    AWAITING_TRANSCODE = 'awaiting_transcode'

    # Storage server (Uploading to upstore.net)
    UPLOADING = 'uploading'
    UPLOADED = 'uploaded'
    UPLOADING_FAILED = 'uploading_failed'
    AWAITING_UPLOAD = 'awaiting_upload'

    # Central server
    PUBLISHING = 'publishing'
    PUBLISHED = 'published'
    PUBLISHING_FAILED = 'publishing_failed'
    AWAITING_PUBLISH = 'awaiting_publish'

    DELETED = 'deleted'


class AbstractRecording(dict):
    def is_published_on(self, site: Site) -> bool:
        """
        Check if the recording has been published on the given site.

        :param site:
        :return:
        """

        if not isinstance(site, Site):
            raise ValueError('Expected an instance of of Site')

        published_on = self.get('publishedOn')

        if published_on is None or len(published_on) == 0:
            return False

        for post in published_on:
            if post['site'] == site['id']:
                return True

        return False

    def __setitem__(self, key, value):

        if key == 'state':
            if isinstance(value, str):
                value = RecordingState(value)

            if not isinstance(value, RecordingState):
                raise ValueError('Invalid state value, expected a RecordingState')

        super().__setitem__(key, value)

    def update(self, e=None, **f):
        if 'state' in f:
            value = f['state']

            if isinstance(value, str):
                value = f['state'] = RecordingState(value)

            if not isinstance(value, RecordingState):
                raise ValueError('Invalid state value, expected a RecordingState')

        super().update(E=e, **f)


class AssociatedRecording(AbstractRecording):
    """
    Recording that has an association to a performer
    """
    pass


class UnassociatedRecording(AbstractRecording):
    """
    Recording that does not have an association to a performer
    """
    pass


class RecordingImage(dict):
    """
    Recording image, represents the images used by infinity
    """
    pass


class RecordingApi(AbstractApiClient):
    @property
    def performer_collection_endpoint(self) -> str:
        """
        Performer restricted collection endpoint api

        :return:
        """
        return self._base_uri + '/performers/{performer_id}/recordings'

    @property
    def collection_endpoint(self) -> str:
        """
        Global api recordings endpoint

        :return:
        """
        return self._base_uri + '/recordings'

    @property
    def resource_endpoint(self) -> str:
        """
        API resource endpoint uri

        :return:
        """
        return self._base_uri + '/recordings/{recording_id}'

    @property
    def association_endpoint(self) -> str:
        """
        Endpoint for creating associations to wordpress posts

        :return:
        """
        return self._base_uri + '/recordings/{recording_id}/posts'

    @staticmethod
    def _hydrate(data: dict):

        type = data.get('type')

        if type == 'associated':
            return AssociatedRecording(data)
        elif type == 'unassociated':
            return UnassociatedRecording(data)

        raise ValueError('Unsupported recording with type "{]"'.format(type))

    def get(self, recording_id, identifier=None) -> AbstractRecording:

        if identifier in ('path', 'location', 'gallery-url', 'video-url'):
            if isinstance(recording_id, str):
                recording_id = bytes(recording_id, 'utf-8')

            recording_id = base64.b64encode(recording_id).decode('utf-8')

        response = self._session.get(self.resource_endpoint.format(recording_id=recording_id),
                                     params=dict(identifier=identifier))

        if response.status_code == 404:
            raise RecordingNotFoundError()

        elif response.status_code == 409:
            raise MultipleRecordingsFoundError(identifiers=response.json()['errors']['identifiers'])

        response.raise_for_status()

        return self._hydrate(response.json())

    def get_all(self, limit=50, offset=0, **kwargs) -> CollectionResult:

        # Update kwargs with the other parameters
        kwargs.update(limit=limit, offset=offset)

        response = self._session.get(self.collection_endpoint, params=kwargs)
        response.raise_for_status()

        result = response.json()

        items = tuple(self._hydrate(recording) for recording in result['data'])
        total = result['meta']['total']
        offset = result['meta']['offset']

        return CollectionResult(items, total, offset)

    def delete(self, recording: AbstractRecording):
        """
        Delete a recording

        :param recording:
        :return:
        """

        if not isinstance(recording, AbstractRecording):
            raise ValueError('Expected a recording instance')

        if 'id' not in recording:
            raise KeyError('Recording is missing an identifier')

        response = self._session.delete(self.resource_endpoint.format(recording_id=recording['id']))

        if response.raise_for_status() == 404:
            raise PerformerNotFoundError()

        response.raise_for_status()

    def create_unassociated_recording(self, **kwargs) -> AbstractRecording:
        """
        Create a new unassociated recording

        :param kwargs:
        :return:
        """

        response = self._session.post(self.collection_endpoint, json=kwargs)
        response.raise_for_status()

        return AbstractRecording(response.json())

    def create_associated_recording(self, performer: AbstractPerformer, **kwargs) -> AbstractRecording:
        """
        Create a recording associated to the given performer

        :param performer:
        :param kwargs:
        :return:
        """

        if not isinstance(performer, AbstractPerformer):
            raise ValueError('Expected a AbstractPerformer instance')

        response = self._session.post(self.performer_collection_endpoint.format(performer_id=performer['id']),
                                      json=kwargs)

        if response.status_code == 404:
            raise PerformerNotFoundError()

        response.raise_for_status()

        return AbstractRecording(response.json())

    def update(self, recording: AbstractRecording):
        """
        Update a recording

        :param recording:

        :return:
        """

        if not isinstance(recording, AbstractRecording):
            raise ValueError('Expected a AbstractRecording instance')

        if 'id' not in recording:
            raise KeyError('Missing recording identifier')

        if 'state' not in recording:
            raise ValueError('Missing a valid recording state for {}'.format(recording['id']))

        response = self._session.put(self.resource_endpoint.format(recording_id=recording['id']), json=recording)

        if response.status_code == 404:
            raise RecordingNotFoundError()

        response.raise_for_status()

        # Update the recording
        recording.update(**response.json())

    def set_recording_state(self, recording: AbstractRecording, state: RecordingState):
        """
        Change the recording state of a performer, simple utility function.

        :param recording:
        :param state:

        :return:
        """

        if not isinstance(recording, AbstractRecording):
            raise ValueError('Expected a AbstractRecording instance')

        if not isinstance(state, RecordingState):
            raise ValueError('Expected a RecordingState instance')

        # Update the performer
        recording.update(state=state)

        # Do a regular update
        self.update(recording)

    def add_post_association(self, recording: AbstractRecording, site: Site, post: common.wordpress.WPPost):
        """
        Add a wordpress post association

        :param recording:
        :param post:
        :return:
        """

        if not isinstance(recording, AbstractRecording):
            raise ValueError('Expected a AbstractRecording instance')

        if not isinstance(site, Site):
            raise ValueError('Expected a Site instance')

        if not isinstance(post, common.wordpress.WPPost):
            raise ValueError('Expected a WPPost instance')

        if 'id' not in recording:
            raise KeyError('AbstractRecording is missing it\'s identifier "id"')

        data = {
            'site': site['id'],
            'post': post['id'],
        }

        response = self._session.post(self.association_endpoint.format(recording_id=recording['id']), json=data)
        response.raise_for_status()

    def get_recordings_matching(self, max_results=None, **criteria):
        """
        Retrieve all recordings, or until we hit the max results of recordings matching the given criteria

        :param max_results:
        :param criteria:

        :return:
        """

        limit = 500
        offset = 0
        yielded = 0

        while True:

            result = self.get_all(limit=limit, offset=offset, **criteria)

            for recording in result.items:
                yield recording

                yielded += 1

                if max_results is not None and max_results <= yielded:
                    break

            # Short circuit
            if max_results is not None and max_results <= yielded:
                break

            # No more recordings left
            if offset > result.total:
                break

            offset += limit


def recording_api_factory(prefix: str):
    return RecordingApi(
        base_uri=get_kv('{}/aphrodite/api'.format(prefix)),
        token=get_kv('{}/aphrodite/token'.format(prefix)),
    )
