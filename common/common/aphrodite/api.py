import requests.adapters


class ApiError(RuntimeError):
    pass


class NotFoundError(ApiError):
    pass


class CollectionResult:
    def __init__(self, items: tuple, total: int, offset: int):
        self.items = items
        self.total = total
        self.offset = offset


class AbstractApiClient:

    def __init__(self, base_uri, token: str):
        self._base_uri = base_uri

        session = requests.Session()
        session.mount('http://', requests.adapters.HTTPAdapter(max_retries=3))
        session.headers.update({'Authorization': 'server {}'.format(token)})

        self._session = session
