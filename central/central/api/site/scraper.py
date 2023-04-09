import requests.adapters

# Create a new request session
session = requests.Session()
session.mount('https://chaturbate.com', requests.adapters.HTTPAdapter(max_retries=5))

session.headers.update({
    'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36',
    'Referer': 'https://chaturbate.com/auth/login/'
})


class AbstractScraper:
    def __init__(self, api_uri: str, log):
        self.log = log
        self.api_uri = api_uri
        self.session = session

    def scrape(self):
        raise NotImplementedError

    def name(self):
        raise NotImplementedError

    def send(self):
        response = self.session.post(self.api_uri + '/{}/models/_intersect'.format(self.name()), json=list(self.scrape()))
        response.raise_for_status()



