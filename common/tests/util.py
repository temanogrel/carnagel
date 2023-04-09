import unittest
from common.aphrodite.recording import UnassociatedRecording
from common.utils import create_published_file_name


class TestCreatePublishedFileName(unittest.TestCase):

    def chaturbate(self):

        recording = UnassociatedRecording({
            'createdAt': '2015-01-01T14:30:00+0000',
            'service': 'cbc',
            'section': 'couple',
            'stageName': 'helloWorld'
        })

        file_name = create_published_file_name(recording, '.mp4')

        self.assertEqual('helloWorld_010115_1430_couple_Chaturbate.mp4', file_name)

