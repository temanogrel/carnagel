import unittest
from common.aphrodite.recording import AbstractRecording, UnassociatedRecording
from common.aphrodite.site import Site


class SiteTest(unittest.TestCase):

    def test_may_publish_recording_based_on_service(self):

        site = Site({
            'sources': {
                'cbc': {'female': True,  'couple': False}
            }
        })

        self.assertTrue(site.may_publish_recording(AbstractRecording(service='cbc', section='female')))
        self.assertFalse(site.may_publish_recording(AbstractRecording(service='cbc', section='male')))
        self.assertFalse(site.may_publish_recording(AbstractRecording(service='cbc', section='couple')))
        self.assertFalse(site.may_publish_recording(AbstractRecording(service='cbc', section=None)))
        self.assertFalse(site.may_publish_recording(AbstractRecording(service='mfc', section=None)))

    def test_my_publish_unassociated(self):
        site = Site({
            'sources': {
                'unassociated': False
            }
        })

        self.assertFalse(site.may_publish_recording(UnassociatedRecording()))

        site = Site({
            'sources': {
                'unassociated': True
            }
        })

        self.assertTrue(site.may_publish_recording(UnassociatedRecording()))
