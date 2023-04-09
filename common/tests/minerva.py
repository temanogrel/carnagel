import os
import unittest

from common.minerva import MinervaClientApi
import minerva.common_pb2 as common


class MinervaClientTest(unittest.TestCase):
    def test_upload(self):
        image_path = os.path.dirname(os.path.realpath(__file__)) + "/assets/image.jpg"

        client = MinervaClientApi()
        uuid = client.upload(image_path, 1, common.Infinity_Image, dict(capturedAt=10.15))
        file = client.download(uuid)

        self.assertTrue(os.path.exists(file))
        os.unlink(file)

        self.assertEqual(client.download(uuid, target="/tmp/image.jpg"), "/tmp/image.jpg")
        self.assertTrue(os.path.exists("/tmp/image.jpg"))
        os.unlink("/tmp/image.jpg")

        client.request_deletion(uuid)
