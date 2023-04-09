import os
import re
import threading
from datetime import datetime

import iso8601

from common.aphrodite.performer import AbstractPerformer, ChaturbatePerformer

OLD_STRUCTURE_PATTERN = re.compile(r'(.*)/(?P<name>(.*))_(?P<date>(([0-9]{4})201([1-5])))_(?P<time>[0-9]{3,4})_(.*)')


class RepeatedTimer(object):
    def __init__(self, interval, function, *args, **kwargs):

        if not isinstance(interval, int):
            raise ValueError('Interval must be a int')

        if not callable(function):
            raise ValueError('Function must be callable')

        self._timer = None
        self.interval = interval
        self.function = function
        self.args = args
        self.kwargs = kwargs
        self.is_running = False
        self.start()

    def _run(self):
        self.is_running = False
        self.start()
        self.function(*self.args, **self.kwargs)

    def start(self):
        if not self.is_running:
            self._timer = threading.Timer(self.interval, self._run)
            self._timer.start()
            self.is_running = True

    def stop(self):
        self._timer.cancel()
        self.is_running = False


class TemporarySymlink():
    def __init__(self, source, target):
        self._source = source
        self._target = target

    @property
    def target(self):
        return self._target

    @property
    def source(self):
        return self._source

    def __enter__(self):
        if os.path.exists(self.target):
            os.unlink(self.target)

        os.symlink(self.source, self.target)
        os.sync()

        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        os.unlink(self.target)
        os.sync()


def create_published_file_name(recording: dict, ext: str):
    """
    Create the published file name by using properties of the recording

    :param recording:
    :param ext:

    :return:
    """

    created_at = iso8601.parse_date(recording['createdAt'])

    parts = [
        recording['stageName'],
        created_at.strftime('%d%m%y'),
        created_at.strftime('%H%M'),
    ]

    if recording.get('section') is not None:
        parts.append(recording['section'])

    parts.append(get_service_full_name(recording['service']))

    return '_'.join(filter(None, parts)) + ext


def get_service_full_name(service: str) -> str:
    if service == 'mfc':
        return 'MyFreeCams'

    elif service == 'cbc':
        return 'Chaturbate'

    elif service == 'cam4':
        return 'Cam4'


def extract_data_from_recording(file: str) -> tuple:
    if file.startswith('/opt/transcoded-legacy'):
        return extract_data_from_old_structure(file)

    return extra_data_from_new_structure(file[16:])


def extract_data_from_old_structure(file):
    """
    :param file:
    :return:
    """

    misc = dict()

    if 'myfreecams' in file.lower():
        service = 'mfc'
    elif 'chaturbate' in file.lower():
        service = 'chaturbate'

        if 'Couples' in file:
            misc.update(section='Couples')
        else:
            misc.update(section='Female')

    else:
        raise ValueError('Could not determinate service')

    matches = OLD_STRUCTURE_PATTERN.search(file)

    if not matches:
        raise ValueError('Could not determinate recording date')

    username = matches.group('name')

    while '__' in username:
        username = username.replace('__', '_')

    date = str(matches.group('date'))
    time = str(matches.group('time'))

    # If time is only three digits it means the it's lacking a zero in the beginning
    if len(time) == 3:
        hour = time[-2]
        minute = time[1:]
    else:
        hour = time[:-2]
        minute = time[2:]

    created_at = datetime(
        # Year
        int(date[4:]),

        # Month
        int(date[2:-4]),

        # day
        int(date[0:-6]),

        # hour
        int(hour),

        # minute
        int(minute)
    )

    return service, created_at, username.lower(), misc


def extra_data_from_new_structure(file):
    """
    The new file structure is the following

    /opt/transcoded/<service>_yyyy-mm-dd/<stage_name>_mmss(_section)?.mp4

    A few examples
    /opt/transcoded/mfc_2013-12-20/sweet_nichole_0806.mp4

    :param file:
    :return:
    """

    path, extension = os.path.splitext(file)
    service_and_date, model_and_time = path.split('/')
    service, date = service_and_date.split('_')
    username, timestamp = model_and_time.rsplit('_', 1)

    created_at = datetime.strptime('{} {}:{}'.format(date, timestamp[:2], timestamp[2:]), '%Y-%m-%d %H:%M')

    if service == 'chaturbate':
        username, section = username.rsplit('_', 1)

        return service, created_at, username, dict(section=section)

    return service, created_at, username, dict()


def format_file_size(num, suffix='B'):
    """
    :param num:
    :param suffix:
    :return:
    """
    for unit in ['','Ki','Mi','Gi','Ti','Pi','Ei','Zi']:
        if abs(num) < 1024.0:
            return "%3.1f%s%s" % (num, unit, suffix)
        num /= 1024.0

    return "%.1f%s%s" % (num, 'Yi', suffix)


def created_normalized_filename(performer: AbstractPerformer, created_at: datetime, ext: str) -> tuple:
    """
    Create a normalized file name for placing videos in the storage server

    Returns a tuple of folder and file name
    :param performer:
    :param created_at:
    :param ext:
    :return:
    """

    # Chaturbate is stored as cbc but presented as chaturbate in the storage system
    if isinstance(performer, ChaturbatePerformer):
        folder = 'chaturbate_' + '_' + created_at.strftime('%Y-%m-%d')
    else:
        folder = performer['service'] + '_' + created_at.strftime('%Y-%m-%d')

    name = performer['stageName'].lower()
    name += '_{}.{}.{}'.format(created_at.strftime('%H%M'), performer['id'], ext)

    while '__' in name:
        name = name.replace('__', '_')

    return folder, name