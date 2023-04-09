import re
import requests.adapters

from common.aphrodite.recording import AbstractRecording
from common.consul import get_kv
from common.utils import create_published_file_name, TemporarySymlink

session = requests.Session()
session.mount('http://', requests.adapters.HTTPAdapter(max_retries=3))


class CamGirlGalleryError(RuntimeError):
    def __init__(self, message: str, code: int):
        self._code = code
        self._message = message

    def __str__(self):
        return 'Camgirl.gallery upload failed with code: {}\n{}'.format(self._code, self._message)


class CamGirlGalleryClient:
    CDN_NUMBER = 0
    CDN_MAX_NUMBER = 7

    def __init__(self, api_uri: str, api_key: str):
        self.api_uri = api_uri
        self.api_key = api_key

    @property
    def cdn_number(self) -> int:

        if self.CDN_NUMBER >= self.CDN_MAX_NUMBER:
            self.CDN_NUMBER = 0

        self.CDN_NUMBER += 1

        return self.CDN_NUMBER

    def upload(self, recording: AbstractRecording):

        target_file = create_published_file_name(recording, '_s.jpg')

        with TemporarySymlink(recording['storagePathThumb'], target_file) as symlink:
            params = {
                'key': self.api_key,
            }

            files = {
                'source': open(symlink.target, 'rb')
            }

            response = session.post(self.api_uri, files=files, params=params)

            if response.status_code != 200:
                raise CamGirlGalleryError(response.text, response.status_code)

            image = response.json()['image']

            # Replace the http server to use a round robin cdn server
            pattern = re.compile(r'http://1')

            thumbnail = image['thumb']['url']
            thumbnail = pattern.sub('http://{}'.format(self.cdn_number), thumbnail)

            return dict(thumb=thumbnail, large=image['url'], gallery=image['url_viewer'])


def camgirlgallery_factory(prefix: str) -> CamGirlGalleryClient:
    return CamGirlGalleryClient(
        api_uri=get_kv('{}/camgirlgallery/api'.format(prefix)),
        api_key=get_kv('{}/camgirlgallery/token'.format(prefix)),
    )
