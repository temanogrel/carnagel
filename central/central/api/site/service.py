import random

from common.aphrodite.blacklist import Blacklist
from common.aphrodite.performer import AbstractPerformer


class AbstractService(object):
    def __init__(self, log):
        self.log = log
        self.performers = dict()
        self.blacklist = Blacklist()

    def __contains__(self, item):
        return str(item) in self.performers

    def __getitem__(self, item):
        return self.performers[str(item)]

    def __setitem__(self, key, value):
        self.performers[str(key)] = value

    def __delitem__(self, key):
        del self.performers[str(key)]

    def dispatch_recording(self, performer: AbstractPerformer) -> None:
        """
        Dispatch the recording request

        :param performer: AbstractPerformer

        :return: None
        """
        raise NotImplementedError

    def may_record(self, performer: AbstractPerformer) -> bool:
        """
        Check if the performer passes the criteria for the given service

        :param performer: ChaturbatePerformer

        :return: bool
        """

        raise NotImplementedError

    @staticmethod
    def hydrate(data: dict) -> AbstractPerformer:
        """
        Convert the basic data dict to a performer instance

        :param data:

        :return:
        """
        raise NotImplementedError

    def bootstrap(self):
        """
        Call the aphrodite api can load in all users that are currently marked as online.

        :return:
        """
        raise NotImplementedError

    def get_required_viewer_count(self) -> int:
        raise NotImplementedError

    def get_meta(self) -> tuple:
        """
        Get meta data about the service

        :return: tuple
        """

        viewers = 0
        recording = 0
        pending_recording = 0

        for uid, performer in self.performers.items():
            viewers += performer.current_viewers

            if performer.is_recording:
                recording += 1

            if performer.is_pending_recording:
                pending_recording += 1

        return viewers, len(self.performers), recording, pending_recording

    def get_performer_count(self) -> int:
        """
        Get the current number of performers

        :return: int
        """

        return len(self.performers)

    def intersect_performers(self, data: dict) -> None:
        """
        Compare the current list of performers with a new list

        will either add, update or delete performers based on the new data
        """

        stage_names = []

        for raw in data:

            # Ignore performers with no stageName
            if 'stageName' not in raw:
                self.log.warning('Missing property stageName, not added')
                continue

            # Append the stage name to a list, so we can track which stage names are missing
            stage_names.append(raw['stageName'])

            # if the stage name already exists update it
            if raw['stageName'] in self:
                performer = self[raw['stageName']]
                performer.update(raw)

            else:
                self[raw['stageName']] = self.hydrate(raw)

        # Compare the updated list with all the stage names from the current intersection
        # If the stage name is not in the stage_name list it means we should delete it
        to_delete = [performer for performer in self.performers if performer not in stage_names]

        for performer in to_delete:
            del self[performer]

    def get_random_performer(self) -> AbstractPerformer:
        """
        Get a random performer

        :return
        """
        return self.performers[random.sample(self.performers.keys(), 1).pop()]

    def process_performers(self):
        """
        Process performers determining which performers to record

        :return
        """

        blacklisted, not_synced, dispatched, recording, pending = 0, 0, 0, 0, 0

        # Cache the required viewer count so we don't rape the api
        required_viewer_count = self.get_required_viewer_count()

        for uid, performer in self.performers.items():

            # The performer has yet to be synchronized to aphrodite
            if 'id' not in performer:
                not_synced += 1
                continue

            # Ignore models that are already recording
            if performer.get('isRecording', False):
                recording += 1
                continue

            # Same goes for those pending
            if performer.get('isPendingRecording', False):
                pending += 1
                continue

            if int(performer.get('currentViewers', 0)) < required_viewer_count:
                continue

            # Check configuration for the minimum required number of viewers
            if not self.may_record(performer):
                continue

            if self.blacklist.is_blacklisted(performer):
                blacklisted += 1
                continue

            # each service requires different parameters
            self.dispatch_recording(performer)

            # Increment the number of dispatched performers
            dispatched += 1

            # Update the state of the performer
            performer.update(dict(isRecording=False, isPendingRecording=True))

        return len(self.performers), recording, pending, dispatched, blacklisted, not_synced
