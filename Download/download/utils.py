import os
from datetime import datetime
from download import config


def is_valid_recording(file: str):
    path = os.path.join(config.DOWNLOAD_PATH, file)

    if not os.path.exists(path):
        return False

    # It should be at least a megabyte in size
    return os.path.getsize(path) > 1024


def create_normalized_filename(performer_id: int, performer_name: str, service: str, *args) -> str:
    folder = '{}_{}'.format(service, datetime.now().strftime('%Y-%m-%d'))
    name = performer_name.lower()

    for arg in args:
        name += '_' + arg

    name += '_{}.{}.flv'.format(datetime.now().strftime('%H%M'), performer_id)

    while '__' in name:
        name = name.replace('__', '_')

    return os.path.join(folder, name)
