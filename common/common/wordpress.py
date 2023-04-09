import base64

import iso8601
import requests.adapters

from common.utils import get_service_full_name

session = requests.session()
session.mount('http://', requests.adapters.HTTPAdapter(max_retries=3))


def create_post_title(recording: dict) -> str:
    """
    Create the post tile for word press posts

    :param recording:
    :return:
    """

    from common.aphrodite.recording import AbstractRecording

    if not isinstance(recording, AbstractRecording):
        raise ValueError('Expected a proper AbstractRecording instance')

    created_at = iso8601.parse_date(recording['createdAt'])

    date = created_at.strftime('%d%m%y')
    time = created_at.strftime('%H%M')

    parts = [recording['stageName'], date, time]

    if recording.get('section') is not None:
        parts.append(recording['section'])

    parts.append(get_service_full_name(recording['service']))

    return ' '.join(filter(None, parts))


def create_post_content(recording: dict) -> str:
    """
    Create the post content for word press sites

    :param recording:
    :return:
    """

    from common.aphrodite.recording import AbstractRecording

    if not isinstance(recording, AbstractRecording):
        raise ValueError('Expected a proper AbstractRecording instance')

    structure = '{title}\n' \
                '<a href="{gallery_url}" target="_blank"><img src="{image_thumb}" alt="{stage_name}"/></a>\n' \
                '{description}\n' \
                '<a href="{download_url}"><img src="http://pip.bz/a5f"/></a>'

    return structure.format(
        title=create_post_title(recording),
        stage_name=recording['stageName'],
        image_thumb=recording['imageUrls']['thumb'].replace('camgirl.gallery', 'oop.bz'),
        gallery_url=recording['imageUrls']['gallery'],
        description=recording['description'],
        download_url=recording['videoUrl']
    )


class PostNotFoundError(RuntimeError):
    """
    Raised when the wordpress post is not found

    """
    pass


class WPPost(dict):
    def __getitem__(self, item):
        item = 'ID' if item == 'id' else 'ID'

        return super().__getitem__(item)


class WPApi:
    def __init__(self, api_uri: str, username: str, password: str):
        self.api_uri = api_uri
        self.username = username
        self.password = password

    @property
    def collection_endpoint(self):
        return self.api_uri + '/posts'

    @property
    def resource_endpoint(self):
        return self.api_uri + '/posts/{post_id}'

    def _create_headers(self) -> dict:

        credentials = '{}:{}'.format(self.username, self.password)
        encoded_credentials = base64.standard_b64encode(credentials.encode('utf-8')).decode('utf-8')

        auth = 'Basic {}'.format(encoded_credentials)

        return {
            'Authorization': auth
        }

    def get(self, post_id: int):
        """
        Retrieve a post from the site

        :param post_id:
        :return:
        """

        response = session.get(self.resource_endpoint.format(post_id=post_id), headers=self._create_headers())

        if response.status_code == 404:
            raise PostNotFoundError()

        response.raise_for_status()

        return WPPost(**response.json())

    def update(self, post: WPPost):
        """
        Update a existing post

        :param post:
        :return:
        """

        if not isinstance(post, WPPost):
            raise ValueError('Expected a instance of WPPost')

        if 'ID' not in post:
            raise ValueError('Post has no ID')

        response = session.put(self.resource_endpoint.format(post_id=post['ID']), json=post,
                               headers=self._create_headers())

        if response.status_code == 404:
            raise PostNotFoundError()

        response.raise_for_status()

        return post.update(**response.json())

    def create(self, post: WPPost):
        """
        Create a new post

        :param post:
        :return:
        """

        if not isinstance(post, WPPost):
            raise ValueError('Expected a instance of WPPost')

        response = session.post(self.collection_endpoint, json=post, headers=self._create_headers())
        response.raise_for_status()

        return post.update(**response.json())

    def delete(self, post: WPPost):
        """
        Remove a wordpress post

        :param post:
        :return:
        """

        if not isinstance(post, WPPost):
            raise ValueError('Expected a instance of WPPost')

        if 'ID' not in post:
            raise ValueError('Post has no ID')

        response = session.delete(self.resource_endpoint.format(post_id=post['ID']), headers=self._create_headers())

        if response.status_code == 404:
            raise PostNotFoundError()

        response.raise_for_status()
