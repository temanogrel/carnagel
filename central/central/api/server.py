import tornado.escape
import tornado.ioloop
import tornado.web

from .request_handlers import *
from .service import ServiceContainer


def create_http_server(service_container: ServiceContainer) -> tornado.web.Application:
    """
    Create the http server

    :param service_container: ServiceContainer

    :return:
    """

    handler_kwargs = {
        'service_container': service_container
    }

    return tornado.web.Application([
        # Misc
        tornado.web.url(r'/ping', PingHandler),

        # Services
        tornado.web.url(r'/services', ServiceCollectionHandler, kwargs=handler_kwargs),

        # Models
        tornado.web.url(r'^/(?P<service>[a-z0-9]+)/models/_intersect', ModelIntersectionHandler, kwargs=handler_kwargs),

        # Service specific crap
        tornado.web.url(r'^/(?P<service>[a-z]+)/credentials', CredentialCollectionHandler, kwargs=handler_kwargs),
        tornado.web.url(r'^/(?P<service>[a-z]+)/session_id', SessionIdResourceHandler, kwargs=handler_kwargs),
    ])
