import json
import logging
from random import randint
from urllib.parse import unquote as unquote_entities

import requests
import websocket

from central.api.site.myfreecams.utils import select_server_uri
from central.utils import deep_update

COMMAND_LOGIN = 1
COMMAND_METRIC = 69
COMMAND_USERNAME_LOOKUP = 10
COMMAND_UPDATE_SESSION = 20
COMMAND_UPDATE_ROOM = 44
COMMAND_ROOMDATA = 44
COMMAND_TAGS = 64

RESPONSE_ERROR = 1

logger = logging.getLogger()


def create_server(kill_after_load=False):
    uri, domain = select_server_uri()

    ws = websocket.WebSocketApp(uri)
    server = Server(ws, domain, kill_after_load=kill_after_load)
    protocol = WebSocketProtocol(server)

    ws.on_open = protocol.on_open
    ws.on_error = protocol.on_error
    ws.on_message = protocol.on_message

    return server


class Server():
    def __init__(self, ws: websocket.WebSocketApp, domain: str, kill_after_load=False):
        """
        :param domain: The domain used to connect

        :return:
        """
        self.user_id = None
        self.username = None
        self.session_id = None
        self.kill_after_load = kill_after_load

        self.ws = ws
        self.domain = domain
        self.models = {}
        self.lookups = {}

    def run(self):
        self.ws.run_forever(ping_interval=5)

    def ping(self):
        """
        Apart from the web-socket ping we also need to implement a custom ping

        :return:
        """
        self.ws.send('0 0 0 1 0\n\0')

    def handle(self, message: str):
        args = message.split(' ', maxsplit=5)

        command = int(args[0])
        from_id = int(args[1])
        to_id = int(args[2])
        arg1 = int(args[3])
        arg2 = int(args[4])

        if len(args) == 6:
            payload = args[5]
        else:
            payload = None

        if command == COMMAND_LOGIN:
            return self.on_login(from_id, to_id, arg1, arg2, payload)
        elif command == COMMAND_METRIC:
            return self.load_models(to_id, arg1, arg2, payload)
        elif command == COMMAND_USERNAME_LOOKUP:
            return self.on_username_lookup(arg1, arg2)
        elif command == COMMAND_UPDATE_ROOM:
            return self.update_room_viewer_count(arg1, arg2, payload)
        elif command == COMMAND_UPDATE_SESSION:
            return self.update_session(arg1, arg2, payload)
        elif command == COMMAND_TAGS:
            return self.update_tags(from_id, payload)
        else:
            logger.debug(
                'Ignored command: {} from: {}, to: {}, arg1: {}, arg2: {}, payload: {}'.format(command, from_id, to_id,
                                                                                               arg1, arg2, payload))

    def query_username(self, username: str):

        request_id = None

        while request_id is None or request_id in self.lookups:
            request_id = randint(1000000, 10000000000)

        self.lookups[request_id] = username
        self.ws.send('{} 0 0 {} 0 {}'.format(COMMAND_USERNAME_LOOKUP, request_id, username))

    def on_username_lookup(self, arg1, arg2):
        if arg2 == RESPONSE_ERROR:
            logger.warn('An error response occurred in the username lookup')

    def on_login(self, from_id: int, to_id: int, arg1: int, arg2: int, payload: str):
        if arg1 == 0:

            self.user_id = arg2
            self.username = payload
            self.session_id = to_id

            self.ws.send('{} {} 0 1 0\n\0'.format(COMMAND_ROOMDATA, self.session_id),
                         opcode=websocket.ABNF.OPCODE_BINARY)

    def load_models(self, to_id: int, arg1: int, arg2: int, payload: str):

        if to_id != COMMAND_UPDATE_SESSION:
            logger.warn('Ignoring ADD friend')
            return

        logger.debug('Load metrics, to_id: {}, arg1: {}, arg2: {}, payload: {}'.format(to_id, arg1, arg2, payload))

        if type(payload) is str:

            # Unquote the last argument since thats what contains the
            json_string = unquote_entities(payload)

            # Decode
            metrics = json.loads(json_string)

            uri = 'http://www.myfreecams.com/mfc2/php/mobj.php?f={0}&s={1}'.format(metrics['fileno'], self.domain)

            response = requests.get(uri)
            response.raise_for_status()

            # The json is prefixed by "var g_hModelData = " and suffixed by "LoadModelsFromObject(g_hModelData);"
            # so we need to remove that
            json_string = response.text[19:-38]

            # Removed control characters that a user had in there name
            json_string = json_string.replace('\x1f', '')

            # parse the string
            models = json.loads(json_string)

            # iterate over the models and let the downloading begin
            for uid, data in models.items():

                # If the data is a list it means we are iterating over a models tags
                if isinstance(data, list):
                    continue

                self.models[int(uid)] = data

            logger.info('Loaded models: {}'.format(len(self.models)))

        if self.kill_after_load:
            self.ws.close()

    def update_room_viewer_count(self, arg1: int, arg2: int, payload: str):
        """
        Handles updating the model viewer count

        The data received from the web-socket is a long array structured in a stupid way.
        It's structured in a way that it's [model_id, viewers, model_id, viewers]


        :param arg1: int
        :param arg2: int
        :param payload: string

        :return:
        """

        # This check is done in the javascript so we do it as well
        # doing it allows us to be safe about execution
        if arg1 == 0 and arg2 == 0 and len(payload) > 0:

            counts = json.loads(unquote_entities(payload))

            for index in range(0, len(counts), 2):
                user_id = int(counts[index])
                viewers = int(counts[index + 1])

                if user_id in self.models:
                    self.models[user_id]['viewer_count'] = viewers

    def update_session(self, arg1: int, arg2: int, payload: str):
        """
        Handle adding/removing users

        :param arg1: int
        :param arg2: int

        :param payload: str

        :return:
        """

        logger.debug('Session update, arg1: {}, arg2: {}, payload: {}'.format(arg1, arg2, payload))

        if arg1 == 127:
            if arg2 in self.models:
                logger.info('Removing the model: {}'.format(self.models[arg2]['nm']))

                del self.models[arg2]
            else:
                logger.debug('Trying to remove a model that does not exist')

        else:

            data = json.loads(unquote_entities(payload))

            if arg2 not in self.models:
                if 'nm' in data:
                    self.models[arg2] = data
                    logger.info('Added a new model: {}'.format(data['nm']))
            else:
                self.models[arg2] = deep_update(self.models[arg2], data)

    def update_tags(self, from_id, payload):
        """
        We don't actually care about the models tags, so we just simply ignore implementing this

        :param from_id:
        :param payload:

        :return:
        """
        pass


class WebSocketProtocol():
    def __init__(self, server: Server):
        self.server = server
        self.message_queue = ''

    def on_error(self, ws, error):
        logger.exception(error)

    def on_open(self, ws) -> None:
        """
        After the websocket connection has been completed

        :return:
        """

        ws.send('hello fcserver\n\0')
        ws.send('1 0 0 20071025 0 guest:guest\n\0', opcode=websocket.ABNF.OPCODE_BINARY)

    def on_message(self, ws, message):
        pos = 0
        packets = 0

        if len(self.message_queue) > 0:
            self.message_queue += message
        else:
            self.message_queue = message

        while pos + 4 < len(self.message_queue):

            # get length of data chunk
            n_strlen = int(self.message_queue[pos: pos + 4])

            # get data chunk
            if 0 < n_strlen < 512000:

                if pos + 4 + n_strlen > len(self.message_queue):
                    # We can't process this packet yet, lets queue it until we get more data
                    break

                chunk = self.message_queue[pos + 4: pos + 4 + n_strlen]  # n_strlen includes additional 2?

                if chunk != -1:
                    # Log('Processed chunk ' + packets + ', len: ' + chunk.length + ' == ' + n_strlen + '?')
                    self.server.handle(chunk)
                    packets += 1
                else:
                    print('warning')
                    break

                pos += 4 + n_strlen

            else:
                self.message_queue = ''
                return

        # whatever is left from pos -> queue.length is queued
        if pos != len(self.message_queue):

            if pos > len(self.message_queue):
                self.message_queue = ''
                return

            self.message_queue = self.message_queue[pos:]

        else:
            self.message_queue = ''
