import json
import os
import tempfile
from enum import Enum

import grpc
import requests
import requests.adapters

import minerva.common_pb2 as common
import shutil
from minerva.file_pb2 import RecommendDownloadRequest, RecommendStorageRequest, DeleteRequest
from minerva.file_pb2_grpc import FileStub, LoadBalancerStub

from common.consul import get_service, get_kv


class MinervaError(Exception):
    def __init__(self, status_code: int):
        self.statusCode = status_code


class FileType(int, Enum):
    RECORDING = 1
    WORDPRESS_COLLAGE = 2
    INFINITY_COLLAGE = 3
    INFINITY_SPRITE = 4
    INFINITY_IMAGE = 5


class MinervaClientApi:
    def __init__(self, hostname, address, read_token, write_token):
        self.hostname = hostname
        self.read_token = read_token
        self.write_token = write_token

        self.channel = grpc.insecure_channel(address)

        self.file_client = FileStub(self.channel)
        self.load_balancer_client = LoadBalancerStub(self.channel)

        self.http_client = requests.Session()
        self.http_client.mount('http://', requests.adapters.HTTPAdapter(max_retries=3))

    def download(self, uuid: str, target=None) -> str:
        """
        Download the file using the quickest available path dictated by minerva

        The file is then downloaded to a temporary named file that must be deleted upon usage

        :param target: Will download to the provided target if supplied
        :param uuid: 
        :return: 
        """

        download_query = RecommendDownloadRequest(
            uuid=uuid,
            countHit=False,
            originHostname=self.hostname
        )

        download_response = self.load_balancer_client.RecommendDownload(download_query)
        if download_response.status != common.Ok:
            raise MinervaError(download_response.StatusCode)

        target_host = download_response.edge if download_response.edge != "" else download_response.origin

        query_parameters = dict(host=download_response.origin, path=download_response.path)
        query_headers = dict(Authorization=self.read_token)

        response = self.http_client.get('http://{}:6000/download'.format(target_host), stream=True, params=query_parameters,
                                        headers=query_headers)
        if response.status_code != 200:
            raise MinervaError(response.status_code)

        if target is None:
            target_file = tempfile.NamedTemporaryFile(delete=False)
            target = target_file.name
        else:
            target_file = open(target, 'w+b')

            # copy the download file to the
            shutil.copyfileobj(response.raw, target_file)

        target_file.close()

        return target

    def upload(self, file: str, external_id: int, file_type: int, file_meta: dict) -> str:
        """
        Upload a file to the best minion server dictated by minerva

        :param file: 
        :param external_id: 
        :param file_type: 
        :param file_meta: 
        :return: 
        """

        if not os.path.exists(file):
            raise FileNotFoundError(file)

        upload_query = RecommendStorageRequest(size=os.path.getsize(file), originHostname=self.hostname)

        upload_response = self.load_balancer_client.RecommendStorage(upload_query)
        if upload_response.status != common.Ok:
            raise MinervaError(upload_response.status)

        data = dict(fileType=int(file_type), externalId=external_id, fileMeta=json.dumps(file_meta))
        files = dict(file=open(file, 'rb'))
        headers = dict(Authorization=self.write_token)

        response = requests.post('http://{}:6000/upload'.format(upload_response.hostname), data=data, files=files,
                                 headers=headers)
        if response.status_code != 201:
            raise MinervaError(response.status_code)

        return response.json()['uuid']

    def request_deletion(self, uuid: str) -> None:
        """
        Request a file to be deleted

        :param uuid: 
        :return: 
        """

        response = self.file_client.RequestDeletion(DeleteRequest(uuid=uuid))
        if response.status != common.Ok:
            raise MinervaError(response.status)


def minerva_factory(hostname: str, prefix: str) -> MinervaClientApi:
    index, instances = get_service('minerva')

    return MinervaClientApi(
        address='{}:{}'.format(instances[0]['ServiceAddress'], instances[0]['ServicePort']),
        hostname=hostname,
        read_token=get_kv('{}/minerva/read-token'.format(prefix)),
        write_token=get_kv('{}/minerva/write-token'.format(prefix)),
    )
