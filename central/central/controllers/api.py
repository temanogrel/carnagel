import logging
import os

import tornado.ioloop
from cement.core.controller import CementBaseController, expose

from central import config
from central.api.server import create_http_server
from central.api.service import service_container_factory


class ApiController(CementBaseController):

    class Meta:
        label = 'api'
        description = 'Run the api'

        stacked_on = 'base'
        stacked_type = 'nested'

        arguments = [
            (['-p', '--port'], dict(action='store', type=int, default=80, help='Port to the run the api on')),
        ]

    @expose(hide=True, aliases=['run'])
    def default(self):

        port = os.environ.get('NOMAD_PORT_http', self.app.pargs.port)

        if not self.app.pargs.debug:
            nh = logging.NullHandler()

            for p in ('general', 'application', 'access'):
                logger = logging.getLogger('tornado.{}'.format(p))
                logger.propagate = False
                logger.addHandler(nh)

        loop = tornado.ioloop.IOLoop.instance()

        service_container = service_container_factory(self.app.log, loop)

        self.app.log.info('Configured the services')
        self.app.log.info('RabbitMQ broker: {}'.format(config.BROKER_URL))

        # Create and configure the central server
        http_server = create_http_server(service_container)
        http_server.listen(port)

        self.app.log.info('Created the http server listening on {}'.format(port))

        try:
            loop.start()
        except KeyboardInterrupt:
            pass
