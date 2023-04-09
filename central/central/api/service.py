import json
import os
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime

from tornado.ioloop import PeriodicCallback

from central import config
from central.api.site.chaturbate.service import ChaturbateService
from central.api.site.myfreecams.service import MyFreeCamsService
from central.api.site.service import AbstractService
from common.aphrodite.blacklist import blacklist_api_factory
from common.aphrodite.performer import PerformerServices, performer_api_factory
from common.utils import RepeatedTimer


class ServiceContainer:
    """
    Collection of services
    """

    def __init__(self, log):
        self.log = log
        self.services = dict()
        self.executor = ThreadPoolExecutor(max_workers=2)

        self._is_sending = False
        self._is_dispatching = False

    def add_service(self, name: str, service: AbstractService) -> None:
        if name in self.services:
            raise ValueError('A service with the name "{}" already exists.'.format(name))

        if not isinstance(service, AbstractService):
            raise ValueError('Service must be an instance of AbstractService.')

        self.services[name] = service

    def get_service(self, name) -> AbstractService:
        return self.services[name]

    def get_meta(self):
        for name, service in self.services.items():
            (viewers, models, recording, pending_recording) = service.get_meta()

            yield name, {
                'name': PerformerServices[name].value,
                'models': models,
                'viewers': viewers,
                'recording': recording,
                'pending_recording': pending_recording
            }

    def dispatch_recordings(self) -> None:
        """
        Iterates of all the registered services and tells it to process it's performers.
        If the performer matches the required criteria for that service, then it will dispatch a request
        to record the performer using celery
        """

        if self._is_dispatching:
            return self.log.error('Already dispatching')

        self._send_data()

        self.log.info('Running a dispatch check')
        self._is_dispatching = True

        started_at = datetime.now()

        for name, service in self.services.items():
            try:
                result = service.process_performers()

                # Show some relevant information
                self.log.info('Processed {}, online: {}, recording: {}, pending recording: {} '
                              'dispatched: {}, blacklisted: {}, not synced: {}'.format(name, *result))

            except Exception as e:
                self.log.error('An error occurred when dispatch requests for for {}'.format(name))
                self.log.error(str(e))

        self.log.info('Check dispatch finished in {} ms'.format((datetime.now() - started_at).microseconds / 1000))
        self._is_dispatching = False

    def _send_data(self):
        """
        Send the current performers to the aphrodite api

        :return void
        """

        self.log.info('Sending data to aphrodite')

        started_at = datetime.now()
        performer_api = performer_api_factory('central')

        for name, service in self.services.items():

            # Ignore empty components
            if len(service.performers) == 0:
                continue

            if os.environ.get('DEBUG'):
                with open('service-{}.json'.format(name), 'w') as f:
                    json.dump(service.performers, f)

            try:
                performer_api.intersect_online_performers(name, service.performers)
            except Exception as e:
                self.log.error('An error occurred when sending data for {}'.format(name))
                self.log.error(str(e))

        duration = (datetime.now() - started_at)

        self.log.info('Finished sending data in {}s'.format(duration))

    def bootstrap(self):
        """
        Load the server with the last synchronized state
        """

        for name, service in self.services.items():
            service.bootstrap()

            self.log.info('Bootstrapped {} with {} performers'.format(name, len(service.performers)))

        self.sync_blacklist()

    def sync_blacklist(self):
        """
        Get a updated list from the api
        """

        self.log.info('Syncing blacklist')

        # Retrieve the current blacklist
        blacklist = blacklist_api_factory('central').get()

        # Nuke the blacklist
        for service in self.services.values():
            service.blacklist = blacklist


def service_container_factory(log, io_loop, bootstrap=True) -> ServiceContainer:
    """
    Create and configure the service container

    :param ioloop:

    :return:
    """

    container = ServiceContainer(log)
    container.add_service('cbc', ChaturbateService(log))
    container.add_service('mfc', MyFreeCamsService(log))

    if bootstrap:
        log.info('Running the bootstrap')
        container.bootstrap()

    else:
        log.info('Skipping bootstrap')

    # Update the blacklist every minute
    RepeatedTimer(60, container.sync_blacklist)

    # Scan the components every x milliseconds
    callback = PeriodicCallback(container.dispatch_recordings, 20 * 1000, io_loop=io_loop)
    callback.start()

    return container
