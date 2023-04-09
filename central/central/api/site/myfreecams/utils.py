import random
import re

import requests

EXTRACT_UID_PATTERN = re.compile(r'nProfileUserId:\s(?P<id>\d+),')


def get_server_config() -> dict:
    """
    Get any of the available servers

    :return:
    """
    response = requests.get('http://www.myfreecams.com/mfc2/data/serverconfig.js')
    response.raise_for_status()

    return response.json()


def select_server_uri() -> tuple:
    """
    Retrieve a server uri for any current server that implements the rfc6455

    :return:
    """

    config = get_server_config()
    servers = []

    for server, implementation in config['websocket_servers'].items():
        if implementation == 'rfc6455':
            servers.append(server)

    domain = servers[random.randint(0, len(servers) - 1)]

    return 'ws://{0}.myfreecams.com:8080/fcsl'.format(domain), domain

def username_to_model_id(stage_name: str) -> str:
    """
    Try and associate a username with a model id

    :param stage_name:

    :return:
    """

    response = requests.get('http://profiles.myfreecams.com/{}'.format(stage_name))

    match = EXTRACT_UID_PATTERN.search(response.text)

    if match:
        return match.group('id')

    return None
