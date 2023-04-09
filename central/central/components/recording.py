import re
from datetime import datetime

import requests
requests.packages.urllib3.disable_warnings()

from central.tasks.storage import reupload_recording_image, reupload_recording_video
from common.aphrodite.recording import AbstractRecording, recording_api_factory
from common.hermes import HermesUrlNotFound, hermes_factory

hermes_api = hermes_factory('central')
recording_api = recording_api_factory('central')

IS_UPSTORE_URL = re.compile(r'http(s)?://(upsto\.re|upstore\.net)/')
EXTRACT_LARGE_URL = re.compile(r'(?P<url>http://\d+\.camgirl\.gallery/images/.*_s\.jpg)')


def get_random_published_entry(data: dict) -> list:
    """
    Get a random (due to python dict nature) published entry from the recording

    :param data:
    :return:
    """
    if 'published_on' not in data:
        return None

    for site in data['published_on']:
        if 'id' not in site:
            continue

        return site['id'], site['name']

    return None


def get_site_published_on_entry(data: dict, site: str):
    """
    Check if the data knows about a publication on a specific site, if it does then it returns that site
    else it return False

    :param data:
    :param site:

    :return: bool|dict
    """
    if 'published_on' not in data:
        return None

    for obj in data['published_on']:
        if obj['name'] == site:
            return obj

    return None


def validate_recording_links(recording: AbstractRecording, dry_run: bool = False) -> tuple:
    """
    Validate the links of a recording

    If either image or video links are valid dispatch a re-upload request to the storage server that holds it.

    :param dry_run:
    :param recording:

    :return:source
    """

    valid_image = is_valid_image(recording['imageUrls']['thumb'])

    # This is an old video, so large is actually the gallery url and we don't have the actual large image
    if valid_image and recording['imageUrls']['gallery'] is None:

        gallery = recording['imageUrls']['large']
        large = recording['imageUrls']['thumb'].replace('_s.th.jpg', '_s.jpg')

        recording['imageUrls'].update(large=large, gallery=gallery)

    valid_video = is_valid_video(recording['videoUrl'])

    if not valid_image and not dry_run:
        reupload_recording_image.apply_async(args=(recording['id'],), routing_key=recording['storageServer'])

    if not valid_video and not dry_run:
        reupload_recording_video.apply_async(args=(recording['id'],), routing_key=recording['storageServer'])

    if not dry_run:
        recording.update(
            videoUrlValid=valid_video,
            imageUrlsValid=valid_image,
            lastCheckedAt=datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        )

        recording_api.update(recording)

    return valid_image, valid_video


def is_valid_image(image: str) -> bool:
    """
    Check if the image url exists

    Since the images don't use a short url service, we can check the domain name prior to loading the url, if the url
    is valid we then run a http head request against it to see if the image is still online.

    todo: check if this works for the large image as well, since that url links to the gallery website and not the
    actual image.

    :param image:

    :return:
    """

    if not isinstance(image, str):
        return False

    # Blacklisted sites automatically fail
    if 'uploadyourimages.org' in image or 'thro.bz' in image:
        return False

    return requests.head(image, allow_redirects=True).status_code == 200


def is_valid_video(video: str) -> bool:
    """
    Check if the video url is valid

    Attempts to load the url, on a successful response we check if the domain belongs to a list of blacklisted domains.
    If that is not the case we continue to check if the response contains a string saying that the file was removed
    due to DMCA/Copyright reasons.

    :param video:
    :return:
    """

    if isinstance(video, AbstractRecording):
        video = video['videoUrl']

    if not isinstance(video, str):
        return False

    try:

        # Retrieve the original url from hermes
        url = hermes_api.get(video)

        if IS_UPSTORE_URL.search(url['originalUrl']) is None:
            return False

        response = requests.get(url['originalUrl'], verify=False)
        response.raise_for_status()

        if 'File was deleted by owner or due to a violation of service rules.' in response.text:
            return False

        if 'File not found' in response.text:
            return False

        return True

    except ValueError as e:
        return False

    except HermesUrlNotFound:
        return False

