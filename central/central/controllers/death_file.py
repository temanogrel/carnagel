from datetime import datetime
from pprint import pprint

from cement.core.controller import CementBaseController, expose

from central import process_death_file, config, ThreadPoolExecutor
from central.components.death_file import get_pending_urls, delete_recording_by_url_entry
from common.aphrodite.death_file import UrlEntryState, UrlEntryApi


class DeathFileController(CementBaseController):
    class Meta:
        label = 'death-file'
        description = 'Process death files'

        stacked_on = 'base'
        stacked_type = 'nested'

        arguments = [
            (['-f', '--file'], dict(action='store', type=int, help='File id to process')),
        ]

    @expose()
    def parse(self):
        if not self.app.pargs.file:
            return self.app.log.error('Missing required argument file')

        for entry in process_death_file(self.app.pargs.file):
            pprint(entry)

    @expose()
    def delete_all(self):

        entry_api = UrlEntryApi(config.APHRODITE_API)

        # Retrieves all urls with the state pending
        for batch in get_pending_urls():
            with ThreadPoolExecutor(max_workers=60) as exc:
                futures = [exc.submit(self._delete_pending_url, entry, entry_api) for entry in batch]

                for future in futures:
                    future.result()

    def _delete_pending_url(self, entry, entry_api):

        try:
            started_at = datetime.now()

            # Update the state prior to passing it to the worker pool
            # else the polling with get_pending_urls will get fucked up
            entry.update(state=UrlEntryState.IN_PROGRESS)
            entry_api.update(entry)

            delete_recording_by_url_entry(entry)

            entry.update(state=UrlEntryState.REMOVED)
            entry_api.update(entry)

            duration = ((datetime.now() - started_at).microseconds / 1000)

            self.app.log.info('Processed {} successfully in {}ms'.format(entry['recordingId'], duration))

        except Exception as e:
            self.app.log.warn('Failed to process {}'.format(entry['recordingId']))
            pprint(e)
