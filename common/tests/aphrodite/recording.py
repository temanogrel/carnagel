import json
import unittest

from common.aphrodite.recording import AbstractRecording, RecordingState
from common.aphrodite.site import Site


class RecordingTest(unittest.TestCase):
    def test_setting_invalid_state(self):
        recording = AbstractRecording()

        with self.assertRaises(ValueError):
            recording['state'] = 'test'

        with self.assertRaises(ValueError):
            recording.update(state='test')

    def test_setting_valid_state(self):
        recording = AbstractRecording()
        recording['state'] = RecordingState.DOWNLOADED

    def test_convert_to_json(self):

        expected = json.dumps(dict(state=RecordingState.DOWNLOADED.value))
        recording = AbstractRecording(state=RecordingState.DOWNLOADED)

        self.assertEqual(expected, json.dumps(recording))

    def test_published_on_with_empty_list(self):

        site = Site({
            'sources': {
                'cbc': ['female', 'couple']
            }
        })

        recording = AbstractRecording(service='cbc', publishedOn=())

        self.assertFalse(recording.is_published_on(site))

    def test_published_on_with_existing_association(self):

        site = Site({
            'id': 1,
            'sources': {
                'cbc': ['female', 'couple']
            }
        })

        recording = AbstractRecording(service='cbc', publishedOn=({
            'id': 1,
            'site': 1,
            'postId': 1337
        },))

        self.assertTrue(recording.is_published_on(site))

