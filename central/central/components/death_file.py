import gzip
import os
import re
from enum import Enum

from central import config
from central.components.recording import is_valid_video
from central.components.wordpress import delete_wordpress_posts, delete_ultron_posts
from common.aphrodite.death_file import DeathFile, UrlEntryState, UrlEntry, urlentry_api_factory
from common.aphrodite.recording import RecordingNotFoundError, RecordingState, recording_api_factory
from common.hermes import HermesUrlNotFound, hermes_factory


class IgnoreReason(str, Enum):
    INVALID_VIDEO = 'invalid_video'
    HERMES_URL_MISSING = 'hermes_url_missing'
    APHRODITE_RECORDING_MISSING = 'aphrodite_recording_Missing'


class ProcessUrlResult:
    SUCCESS = 0
    HERMES_URL_NOT_FOUND = 1
    APHRODITE_RECORDING_NOT_FOUND = 2

    def __init__(self, code, hermes=None, recording=None):
        self.code = code
        self.hermes = hermes
        self.recording = recording


def get_pending_urls():
    """
    :rtype: generator[UrlEntry]
    """

    entry_api = urlentry_api_factory('central')

    while True:

        result = entry_api.get_all(limit=5000, state=UrlEntryState.PENDING)

        # This can occur, IF records are removed during the process
        if len(result.items) == 0:
            break

        yield result.items


def get_url_count_from_file(death_file: DeathFile):
    abs_path = os.path.join(config.DEATH_FILE_STORAGE, death_file['location'])

    if not os.path.exists(abs_path):
        raise FileNotFoundError(abs_path)

    with gzip.open(abs_path, mode='r') as f:

        count = 0

        for _ in f:
            count += 1

    return count


def get_urls_from_file(death_file: DeathFile):
    """
    Open the compressed file and yield each url

    :param death_file:

    :return:
    """

    abs_path = os.path.join(config.DEATH_FILE_STORAGE, death_file['location'])

    if not os.path.exists(abs_path):
        raise FileNotFoundError(abs_path)

    pattern = re.compile('^(?P<url>(.*));(?P<filename>.*);(?P<uploaded>(.*));"(?P<downloads>\d+)"$')

    with gzip.open(abs_path, mode='r') as f:
        for index, url in enumerate(f):

            url = url.decode('utf-8').strip()  # type: str

            # This filters out the download count
            result = pattern.match(url)
            if result and int(result.group('downloads')) == 0:
                yield result.group('filename'), result.group('url')


def create_url_entry(url: str, file_name: str):
    result = process_url_existence(url)

    if result.code == ProcessUrlResult.SUCCESS:

        if is_valid_video(result.recording):
            return UrlEntry(url=url, state=UrlEntryState.PENDING, hermesId=result.hermes['id'],
                            recording=result.recording['id'], filename=file_name)
        else:
            return UrlEntry(url=url, state=UrlEntryState.IGNORED, hermesId=result.hermes['id'], filename=file_name,
                            recording=result.recording['id'], ignoreReason=IgnoreReason.INVALID_VIDEO)

    elif result.code == ProcessUrlResult.APHRODITE_RECORDING_NOT_FOUND:
        return UrlEntry(url=url, state=UrlEntryState.IGNORED, hermesId=result.hermes['id'], filename=file_name,
                        ignoreReason=IgnoreReason.APHRODITE_RECORDING_MISSING)

    elif result.code == ProcessUrlResult.HERMES_URL_NOT_FOUND:
        return UrlEntry(url=url, state=UrlEntryState.IGNORED, ignoreReason=IgnoreReason.HERMES_URL_MISSING,
                        filename=file_name)

    else:
        raise ValueError('Unhandled process url result code')


def process_url_existence(url: str) -> ProcessUrlResult:
    """
    Run various tests against the url and return the resulting recording if found

    :param url:
    :return:
    """

    hermes_api = hermes_factory('central')
    recording_api = recording_api_factory('central')

    try:
        hermes = hermes_api.get_by_original_url(url)
    except HermesUrlNotFound:
        return ProcessUrlResult(ProcessUrlResult.HERMES_URL_NOT_FOUND)

    try:
        recording = recording_api.get(hermes.generate_hermes_url(), identifier='video-url')
    except RecordingNotFoundError:
        return ProcessUrlResult(ProcessUrlResult.APHRODITE_RECORDING_NOT_FOUND, hermes)

    return ProcessUrlResult(ProcessUrlResult.SUCCESS, hermes, recording)


def delete_recording_by_url_entry(entry: UrlEntry):
    raise NotImplementedError("Supported for minerva is not yet present")

    if not isinstance(entry, UrlEntry):
        raise ValueError('First argument entry should be a valid UrlEntry instance')

    recording_api = recording_api_factory('central')
    recording = recording_api.get(entry['recordingId'])

    delete_wordpress_posts(recording)
    delete_ultron_posts(recording)

    recording.update(state=RecordingState.DELETED)
    recording_api.update(recording)
