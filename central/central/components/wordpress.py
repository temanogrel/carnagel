import re

from common.aphrodite.recording import AbstractRecording
from common.aphrodite.site import site_api_factory
from common.ultron import UltronRecordingNotFoundError, ultron_factory
from common.wordpress import WPApi, WPPost, PostNotFoundError

SIZE_DETECTION_PATTERN = re.compile('size:\d+\.\d{2}(mb|gib)', re.IGNORECASE)
SIZE_EXTRACTION_PATTERN = re.compile('size: (?P<size>\d+) bytes', re.IGNORECASE)

EXTRACT_POST_NAME = re.compile(r'/(?P<post_name>([a-z0-9]|-|_)+)/$', re.IGNORECASE)
EXTRACT_DURATION = re.compile(r'duration: (?P<duration>\d{2}:\d{2}:\d{2})', re.IGNORECASE)
EXTRACT_IMAGE_THUMB_URL = re.compile(r'src=("|\')(?P<uri>(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*))("|\')', re.IGNORECASE)
EXTRACT_IMAGE_LARGE_URL = re.compile(r'href=("|\')(?P<uri>(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*))("|\')', re.IGNORECASE)
EXTRACT_VIDEO_URL = re.compile(r'href=("|\')(?P<uri>http://(pip|cur)\.bz/([a-zA-Z0-9])+)("|\')', re.IGNORECASE)


def detect_bbcode(content: str) -> bool:
    """
    Check if the post contains bbcode

    :param content:
    :return:
    """
    return '[/url' in content.lower()


def detected_size_formatted(content: str) -> bool:
    """
    Checks if the size of the video is denoted in a formatted version

    :param content:

    :return:
    """
    return SIZE_DETECTION_PATTERN.search(content) is not None


def get_post_size(content: str) -> str:
    """
    Get the size in bytes of the current video

    :param content:
    :return:
    """
    matches = SIZE_EXTRACTION_PATTERN.search(content)

    if not matches:
        return None

    return matches.group('size')


def extract_video_url(content: str) -> str:
    """
    Extract the video uri from the post

    :param content:
    :return:
    """
    matches = EXTRACT_VIDEO_URL.search(content)

    if not matches:
        return None

    return matches.group('uri')


def extract_image_url(content: str) -> dict:
    """
    Extract the image uri from the post

    :param content:
    :return:
    """
    thumb_matches = EXTRACT_IMAGE_THUMB_URL.search(content)
    large_matches = EXTRACT_IMAGE_LARGE_URL.search(content)

    thumb = thumb_matches.group('uri') if thumb_matches else None
    large = large_matches.group('uri') if large_matches else None

    if 'http://pip.bz/a5f' in (large, thumb):
        raise ValueError('Parsed wrong uri')

    return dict(large=large, thumb=thumb)


def extract_duration(content: str) -> int:
    """
    Extract the duration from the post

    :param content:
    :return:
    """
    matches = EXTRACT_DURATION.search(content)

    if not matches:
        return None

    return sum(int(x) * 60 ** i for i, x in enumerate(reversed(matches.group('duration').split(':'))))


def extract_post_name(url: str) -> str:
    """
    Retrieve the post name from the database

    :param url:
    :return:
    """

    matches = EXTRACT_POST_NAME.search(url)

    if not matches:
        return None

    return matches.group('post_name')


def delete_wordpress_posts(recording: AbstractRecording):
    """
    Remove all the wordpress posts related to a recording

    :param recording:
    :return:
    """

    if not isinstance(recording, AbstractRecording):
        raise ValueError('Invalid recording provided')

    if 'publishedOn' not in recording:
        return

    site_api = site_api_factory('central')
    sites = {site['id']: site for site in site_api.get_all(enabled=1)}

    for post in recording['publishedOn']:

        # The site is probably disabled or removed
        if post['site'] not in sites:
            continue

        site_config = sites[post['site']]

        try:
            wp_api = WPApi(site_config['apiUri'], site_config['username'], site_config['password'])
            wp_api.delete(WPPost(ID=post['postId']))
        except PostNotFoundError:
            pass


def delete_ultron_posts(recording: AbstractRecording):
    try:
        ultron_api = ultron_factory('central')
        ultron_api.delete(recording)
    except UltronRecordingNotFoundError:
        pass
