import random
import re
import time
import urllib.parse
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime

import requests
from lxml import html
from requests.adapters import HTTPAdapter
from requests.exceptions import HTTPError

# Create a new request session
session = requests.Session()
session.mount('https://chaturbate.com', HTTPAdapter(max_retries=5))

session.headers.update({
    'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36',
    'Referer': 'https://chaturbate.com/auth/login/'
})

proxies = {
  "https": "https://104.203.86.39:1080",
}

class ChaturbateScraper:
    def __init__(self, api_uri: str, logging, sections: list, credentials: dict):
        self.performers = {}
        self.api_tokens = {}

        self.api_uri = api_uri
        self.logging = logging
        self.sections = sections
        self.credentials = credentials

    def init(self):
        self.scan_performers()
        self.authenticate_users()
        self.send()

    def send(self):
        response = requests.post(self.api_uri + '/cbc/models/_intersect', json=self.performers)
        response.raise_for_status()

        response = requests.post(self.api_uri + '/cbc/credentials', json=self.api_tokens)
        response.raise_for_status()

    def authenticate_users(self):
        """
        Authenticate all the credentials and use a bunch of models to retrieve the users api_token

        :return:
        """

        # reset api tokens
        self.api_tokens = {}

        for identity, credential in self.credentials.items():
            try:
                token = self._get_api_token(identity, credential)

                self.logging.info('Authenticating identity: {} received token: {}'.format(identity, token))
                self.api_tokens[identity] = token

                # Sleep so we don't trigger a 429 error
                time.sleep(20)

            except (ValueError, HTTPError, RuntimeError) as e:
                self.logging.error(str(e))

        # If none of the users successfully authenticated we run again
        if len(self.api_tokens) == 0:
            self.authenticate_users()

    def scan_performers(self, send=False):
        """
        Spin upp a thread pool executor and scan each section separately

        :return:
        """

        start = datetime.now()
        threads = len(self.sections)
        performers = []

        with ThreadPoolExecutor(max_workers=threads) as executor:

            futures = [executor.submit(self._scan_section_page, section) for section in self.sections]

            for future in futures:
                performers += future.result()

        self.logging.info('Scanned chaturbate in {} found: {} models.'.format(datetime.now() - start, len(performers)))
        self.performers = performers

        if send:
            self.send()

    def _get_api_token(self, identity: str, credential: str):
        """
        Retrieve the api token for a identity

        :param identity:
        :param credential:

        :return:
        """

        stage_name = random.choice(self.performers)['stageName']

        response = session.get('https://chaturbate.com/auth/login/', verify=False)
        response.raise_for_status()

        # parse the response
        tree = html.fromstring(response.text)

        # find the CSRF_TOKEN
        csrf_token = tree.xpath('//input[@name="csrfmiddlewaretoken"]/@value')[0]

        parameters = {
            'next': '',
            'csrfmiddlewaretoken': str(csrf_token),
            'username': identity,
            'password': credential,
            'rememberme': 'on'
        }

        # Attempt to login
        response = session.post('https://chaturbate.com/auth/login/', data=parameters)
        response.raise_for_status()

        # request a model page
        page = session.get('https://chaturbate.com/{}/'.format(stage_name), verify=False)
        page.raise_for_status()

        match = re.compile(r'EmbedViewerSwf\((?P<data>(.|\n)*)\n\s+\);').search(page.text)

        # Failed to find JS function which contains our api token.
        # This is most likely because the model has gone offline
        if not match:
            raise ValueError('Failed to find the javascript code to embed the viewer')

        lines = match.group('data').split(',')

        # Remove excess data
        lines = [line.strip() for line in lines]

        # Remove padding
        lines = [line[1:-1] for line in lines]

        # Decode html entities
        api_token = urllib.parse.unquote(lines[15])

        self.logging.info('Authenticated identity: {} received token: {}'.format(identity, api_token))

        # Return identity & flash api token
        return api_token

    def _scan_section_page(self, section: str, page=1, performers=None):
        """
        Scan a page within a given section for performers

        :param page:
        :param section:

        :return: None
        """

        if not performers:
            performers = list()

        # Create the URI
        uri = 'https://chaturbate.com/{section}-cams/?page={page}'.format(section=section, page=page, proxies=proxies)

        self.logging.info('Loading {0}'.format(uri))

        # Download the page
        response = session.get(uri, verify=False)
        response.raise_for_status()

        # Create a tree object
        tree = html.fromstring(response.text)
        rooms = tree.xpath('//ul[@class="list"]/li/div[@class="details"]')

        # Iterate
        for room_element in rooms:
            stage_name, viewers = self._parse_model_details(room_element)

            performers.append({
                'stageName': stage_name,
                'serviceId': stage_name,
                'currentViewers': viewers,
                'section': section
            })

        # If a next page exists we parse it
        next_page = self._get_next_page(tree, page)

        if next_page is not None:
            self._scan_section_page(section, next_page, performers=performers)

        return performers

    def _get_next_page(self, tree, current_page: int):
        """
        Get the next page or return none

        :param tree:
        :param current_page:

        :return:
        """

        pages = tree.xpath('//ul[@class="paging"]/li/a/text()')[1: -1]

        if len(pages) == 0:
            return None

        if int(pages.pop()) == current_page:
            return None
            
        # Sleep so we don't trigger a 429 error
                time.sleep(20)

        return current_page + 1

    def _parse_model_details(self, tree):
        """
        Extract the username and number of viewers

        :param tree:
        :return:
        """

        username = tree.xpath('.//div[@class="title"]/a/text()')[0].strip()
        cam_info = tree.xpath('.//ul[@class="sub-info"]/li[@class="cams"]/text()')[0].strip()

        match = re.compile('(?P<viewers>(\d+))\s(viewers)', flags=re.IGNORECASE).search(cam_info)

        return username, int(match.group('viewers'))
