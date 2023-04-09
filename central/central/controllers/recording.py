import time
import traceback
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime, timedelta
from itertools import chain

from cement.core.controller import CementBaseController, expose

from central import publish_recording, transcode, upload_media, recording_api_factory
from central.components.recording import validate_recording_links
from central.utils import get_recordings_matching
from common.aphrodite.recording import RecordingState


class RecordingController(CementBaseController):
    class Meta:
        label = 'recording'
        description = 'Recording related utilities'

        stacked_on = 'base'
        stacked_type = 'nested'

        arguments = [
            (['-r', '--recording'], dict(action='store', help='Recording id')),
            (['-s', '--state'], dict(action='store', help='Do something with the given state')),
            (['-l', '--limit'], dict(action='store', type=int, default=None, help='Maximum of recordings to process')),
            (['-d', '--dry-run'], dict(action='store', type=bool, default=False, help='Only check, dont do anything')),
            (['-t', '--threads'], dict(action='store', type=int, default=4, help='Number of threads to use')),
            (['-i', '--interval'], dict(action='store', type=int, default=1, help='Number of seconds to sleep between dispatches'))
        ]

    @expose(help='Publish a recording')
    def publish(self):
        if not self.app.pargs.recording:
            return self.app.log.info('No recording id was provided')

        publish_recording(self.app.pargs.recording, wordpress=False, ultron=False)

    @expose(help='Upload all videos that have given the state, defaults to uploading_failed')
    def upload(self):

        state = self.app.pargs.state if self.app.pargs.state else RecordingState.UPLOADING_FAILED

        for recording in get_recordings_matching(state=state):
            upload_media.apply_async(args=(recording['id'],), routing_key=recording['storageServer'],
                                     exchange='storage')

    @expose(help='Transcode all videos that have given the state, defaults to downloaded')
    def transcode(self):

        state = self.app.pargs.state if self.app.pargs.state else RecordingState.DOWNLOADED

        for recording in get_recordings_matching(state=state):
            transcode.delay(recording['id'])

    @expose(help='Fix all recordings that are stuck')
    def fix_stuck(self):

        # Re-queue all failed published recording
        for recording in get_recordings_matching(state=RecordingState.PUBLISHING_FAILED, orphaned=0):

            self.app.log.info('Dispatch publish for {}'.format(recording['id']))

            publish_recording.delay(recording['id'])

        # Only apply the method to recording that have not been updated in the last six hours
        timestamp = (datetime.now() - timedelta(hours=6)).timestamp()

        for recording in get_recordings_matching(state=RecordingState.TRANSCODING, after=timestamp, orphaned=0):

            self.app.log.info('Dispatch transcode for {}'.format(recording['id']))

            transcode.delay(recording['id'])

        for recording in get_recordings_matching(state=RecordingState.PUBLISHING, after=timestamp, orphaned=0):

            self.app.log.info('Dispatch publish for {}'.format(recording['id']))

            publish_recording.delay(recording['id'])

        # Only apply the fix to stuck recordings not updated with in the last 48 hours
        timestamp = (datetime.now() - timedelta(hours=48)).timestamp()

        # chain all sources
        source_chain = (

            # It can take up to 48 hours to upload a video if the queue grows insanely
            get_recordings_matching(state=RecordingState.TRANSCODED, after=timestamp, orphaned=0),

            # Just get all failed recording
            get_recordings_matching(state=RecordingState.UPLOADING_FAILED, orphaned=0)
        )

        for recording in chain(*source_chain):
            self.app.log.info('Dispatch upload for {} to {}'.format(recording['id'], recording['storageServer']))

            upload_media.apply_async(args=(recording['id'],), routing_key=recording['storageServer'])

            time.sleep(self.app.pargs.interval)

    @expose(help='Validate a single recording')
    def validate(self):
        if self.app.pargs.recording is None:
            raise ValueError('Recording argument is missing')

        recording_api = recording_api_factory('central')
        recording = recording_api.get(self.app.pargs.recording)

        image, video = validate_recording_links(recording, self.app.pargs.dry_run)

        self.app.log.info('Image: {}, video: {}'.format(image, video))

    @expose(help='Ensure that the links are valid')
    def validate_all(self):
        """
        Utility to validate the links of the recordings

        Takes a sizable chunk of videos that have a last checked at date older than 10 days and validates that the links
        are up to date and still functional.

        :return:
        """

        started_at = datetime.now()
        last_checked_at = (datetime.now() - timedelta(days=10)).timestamp()

        with ThreadPoolExecutor(max_workers=self.app.pargs.threads) as executor:

            futures = []

            try:
                for recording in get_recordings_matching(max_results=self.app.pargs.limit, checkedAt=last_checked_at,
                                                         state=RecordingState.PUBLISHED, orphaned=0):

                    self.app.log.debug('Dispatching recording: {}'.format(recording['id']))

                    # Submit the task to the thread pool
                    future = executor.submit(validate_recording_links, recording, self.app.pargs.dry_run)

                    # queue the future with the recording id
                    futures.append((future, recording))

            except Exception as e:
                self.app.log.error('Exception during retrieval: {}'.format(str(e)))
                traceback.print_exc()

            # Notify when retrieving has finished
            self.app.log.info('Finished retrieving recordings, starting process')

            total = len(futures)
            videos = 0
            images = 0

            for future, recording in futures:

                try:
                    image, video = future.result()
                except Exception as e:
                    self.app.log.error('Exception: {}'.format(str(e)))
                    traceback.print_exc()
                    continue

                # increment values
                images += int(image)
                videos += int(video)

                args = recording['id'], image, video

                self.app.log.info(
                    'Recording {}, image status: {}, video status: {}'.format(*args)
                )

        duration = (datetime.now() - started_at).seconds

        self.app.log.info(
            'Processed {} recordings, {} valid videos and {} valid images, duration: {}s'.format(
                total, videos, images,
                duration)
        )
