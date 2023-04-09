from logging import getLogger
from threading import RLock

import requests.adapters
from requests.packages.urllib3.exceptions import HTTPError

from central.api.site.myfreecams.websocket import Server, create_server
from central.utils import deep_update

logger = getLogger()

session = requests.Session()
session.mount('http://', requests.adapters.HTTPAdapter(max_retries=5))


def intersect_models(server: Server, lock: RLock, api_uri: str):
    """
    Restructure the server data and send it

    :param server: Server
    :param lock: Rlock
    :param api_uri: str

    :return:
    """

    if server.ws.sock is None:
        logger.info('Main server is no longer online, terminating...')
        return

    with lock:

        logger.info('Transmitting {} models'.format(len(server.models)))

        data = []

        for uid, performer in server.models.items():
            try:
                data.append({
                    'serviceId': str(uid),
                    'stageName': performer['nm'],
                    'currentViewers': int(performer['viewer_count']),

                    'videoState': performer['vs'],
                    'accessLevel': performer['lv'],

                    'camServer': performer['u']['camserv'],
                    'camSCore': performer['m']['camscore'],

                    'missMfcRank': performer['m']['rank']
                })
            except KeyError as e:
                logger.debug('Missing key in mfc scraper {}'.format(e.args[0]))

    try:
        response = session.post(api_uri + '/mfc/session_id', json=dict(session_id=server.session_id))
        response.raise_for_status()

        response = session.post(api_uri + '/mfc/models/_intersect', json=data)
        response.raise_for_status()

    except HTTPError as e:
        logger.exception(e)


def reload_models_at_interval(server: Server, lock: RLock) -> None:
    logger.info('Initiating reload models')

    if server.ws.sock is None:
        logger.info('Main server is no longer online, terminating...')
        return

    s2 = create_server(kill_after_load=True)
    s2.run()

    with lock:
        for uid, data in s2.models.items():

            uid = int(uid)

            if uid not in server.models:
                server.models[uid] = data
            else:
                orig_data = server.models[uid]
                new_data = deep_update(server.models[uid], data)

                if orig_data != new_data:
                    logger.info('Updated model: {} from scan'.format(uid))

                server.models[uid] = new_data


def send_ping_at_interval(server: Server):
    logger.info('Sending a custom ping to the web-socket server')

    if server.ws.sock is None:
        logger.info('Main server is no longer online, terminating...')
        return

    server.ping()
