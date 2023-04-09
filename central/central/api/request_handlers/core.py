import tornado.escape
import tornado.web

from central.api.service import ServiceContainer
from central.api.site.service import AbstractService


class JsonRequestHandler(tornado.web.RequestHandler):
    """
    Our generic json request handler
    """
    def __init__(self, *args, **kwargs):
        self.service_container = None

        super().__init__(*args, **kwargs)

    def data_received(self, chunk):
        raise NotImplemented

    def initialize(self, service_container: ServiceContainer):
        self.service_container = service_container

    def get_service(self, service: str) -> AbstractService:
        return self.service_container.get_service(service)

    def prepare(self):
        if self.request.body:
            try:
                self.request.json = tornado.escape.json_decode(self.request.body)
            except ValueError:
                self.send_error(400, message='Malformed json')

    def set_default_headers(self):
        self.set_header('Content-Type', 'application/json')

    def write_error(self, status_code, **kwargs):
        self.set_status(status_code)

        if 'exc_info' in kwargs:
            return super().write_error(status_code, **kwargs)

        if 'message' not in kwargs:
            if status_code == 405:
                kwargs['message'] = 'Method Not Allowed'
            else:
                kwargs['message'] = 'Unknown error occurred'

        self.write(kwargs)


class PingHandler(tornado.web.RequestHandler):
    def get(self, *args, **kwargs):
        self.set_status(204)
