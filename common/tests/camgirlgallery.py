import os
from unittest import TestCase
from common.aphrodite.recording import UnassociatedRecording
from common.camgirlgallery import CamGirlGalleryClient


class CamgirlGalleryClientTest(TestCase):

    def setUp(self):
        self.client = CamGirlGalleryClient('http://camgirl.gallery/myapi/1/upload/', '14df3e1f7035e4084ad4e31b3fb25de4')

    def test_get_cdn_number(self):
        self.assertEqual(1, self.client.cdn_number)
        self.assertEqual(2, self.client.cdn_number)
        self.assertEqual(3, self.client.cdn_number)
        self.assertEqual(4, self.client.cdn_number)
        self.assertEqual(5, self.client.cdn_number)
        self.assertEqual(6, self.client.cdn_number)
        self.assertEqual(7, self.client.cdn_number)
        self.assertEqual(1, self.client.cdn_number)

    def test_upload(self):
        recording = UnassociatedRecording({
            'createdAt': '2015-01-01T14:30:00+0000',
            'service': 'cbc',
            'section': 'couple',
            'stageName': 'helloWorld',

            'storagePathThumb':  os.path.dirname(os.path.realpath(__file__)) + '/assets/image.jpg'
        })

        images = self.client.upload(recording)

        self.assertRegex(images['gallery'], r'^http://camgirl.gallery/image/([a-zA-Z0-9]+)$')
        self.assertRegex(images['large'], r'^http://\d.camgirl.gallery/images/(.*)$')
        self.assertRegex(images['thumb'], r'^http://\d.camgirl.gallery/images/(.*)$')
