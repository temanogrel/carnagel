from multiprocessing import RLock

from cement.core.controller import CementBaseController, expose

from central import config
from central.api.site.cam4.scraper import Cam4Scraper
from central.api.site.chaturbate.scraper import ChaturbateScraper
from central.api.site.myfreecams.scraper import send_ping_at_interval, intersect_models, reload_models_at_interval
from central.api.site.myfreecams.websocket import create_server
from common.consul import get_service
from common.utils import RepeatedTimer


class DownloadingController(CementBaseController):
    class Meta:
        label = 'downloading'
        description = 'Manage the downloading/scraping servers'

        stacked_on = 'base'
        stacked_type = 'nested'

        arguments = [
            (['-i', '--interval'], dict(action='store', default=60, help='Interval between scans')),
        ]

    def _get_modelserver_address(self) -> str:
        index, instances = get_service('modelserver')

        return 'http://{}:{}'.format(instances[0]['ServiceAddress'], instances[0]['ServicePort'])

    @expose(help='Initiate the MFC scraper')
    def scrape_mfc(self):
        lock = RLock()
        server = create_server()

        self.app.log.info('Initiated the server')

        # Custom ping
        RepeatedTimer(config.MFC_PING_INTERVAL, send_ping_at_interval, server)

        # Update the model server
        RepeatedTimer(config.MFC_INTERSECT_INTERVAL, intersect_models, server, lock, self._get_modelserver_address())

        # Reload the models from mfc
        RepeatedTimer(config.MFC_RELOAD_MODELS_INTERVAL, reload_models_at_interval, server, lock)

        self.app.log.info('Running the server')

        server.run()

    @expose(help='Initiate the chaturbate scraper')
    def scrape_cbc(self):
        scraper = ChaturbateScraper(self._get_modelserver_address(), self.app.log, config.CBC_SECTIONS,
                                    config.CBC_CREDENTIALS)
        scraper.init()

        # Scrape the models
        RepeatedTimer(self.app.pargs.interval, scraper.scan_performers, send=True)

        # Re-authenticate the users
        RepeatedTimer(60 * self.app.pargs.interval, scraper.authenticate_users)

    @expose()
    def scrape_cam4(self):
        scraper = Cam4Scraper(self.app.pargs.host, self.app.log)
        scraper.send()
