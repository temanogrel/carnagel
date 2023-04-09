import math

from lxml import html

from central.api.site.scraper import AbstractScraper


def process_details(performer, section):
    viewers = performer.xpath('.//span[@class="viewers"]/text()')

    if len(viewers) == 0:
        viewers = 0
    else:
        viewers = int(viewers[0].replace(',', '').strip())

    return {
        'stageName': performer.xpath('.//div[@class="profileBoxTitle"]/a/text()')[0],
        'serviceId': performer.xpath('.//div[@class="profileBoxTitle"]/a/text()')[0],
        'section': section,
        'currentViewers': viewers
    }


class Cam4Scraper(AbstractScraper):
    def name(self):
        return 'cam4'

    def get_total_count(self, section):

        if section == 'female':
            uri = 'http://en.cam4.se/directoryCounts?online=true&gender=female'
        elif section == 'male':
            uri = 'http://en.cam4.se/directoryCounts?online=true&gender=male'
        elif section == 'couple':
            uri = 'http://en.cam4.se/directoryCounts?online=true&broadcastType=male_group&broadcastType=female_group&broadcastType=male_female_group'
        else:
            raise ValueError('Unsupported section: "{}"'.format(section))

        response = self.session.post(uri)
        response.raise_for_status()

        return int(response.json()['totalCount'])

    def scrape(self):

        for section in ('female',):
            total = self.get_total_count(section)
            pages = math.ceil(total / 32)

            for page in range(1, pages):
                for performer in self._scrape_page(section, page):
                    yield performer

    def _scrape_page(self, section, page):

        if section == 'female':
            uri = 'http://en.cam4.se/directoryResults?page={page}&online=true&gender=female'
        elif section == 'male':
            uri = 'http://en.cam4.se/directoryResults?page={page}&online=true&gender=male'
        elif section == 'couple':
            uri = 'http://en.cam4.se/directoryResults?page={page}&online=true&broadcastType=male_group&broadcastType=female_group&broadcastType=male_female_group'
        else:
            raise ValueError('Unsupported section: "{}"'.format(section))

        response = self.session.get(uri.format(page=page))
        response.raise_for_status()

        tree = html.fromstring(response.json()['html'])
        performers = tree.xpath('//div[@class="profileDetailBox"]')

        for performer in performers:
            yield process_details(performer, section)
