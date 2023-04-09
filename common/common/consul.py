import consul

consul_client = consul.Consul()


def get_kv(key: str) -> str:
    index, data = consul_client.kv.get(key)
    if data is None:
        raise RuntimeError('Missing key {} in consul'.format(key))

    return data['Value'].decode('UTF-8')


def get_service(name: str) -> tuple:
    index, instances = consul_client.catalog.service(name)

    if len(instances) == 0:
        raise RuntimeError('Missing service {} in consul'.format(name))

    return index, instances
