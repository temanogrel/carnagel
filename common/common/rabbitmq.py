from .consul import get_kv, get_service


def rabbitmq_dsn(consul_prefix: str):
    """
    Connect to the rabbitmq server using service specific credentials 
    
    :param consul_prefix: 
    :return: 
    """

    index, nodes = get_service('rabbitmq')

    return 'amqp://{username}:{password}@{hostname}:{port}/{vhost}'.format(
        port=nodes[0]['ServicePort'],
        hostname=nodes[0]['ServiceAddress'],

        vhost=get_kv('{}/rabbitmq/vhost'.format(consul_prefix)),
        username=get_kv('{}/rabbitmq/user'.format(consul_prefix)),
        password=get_kv('{}/rabbitmq/pass'.format(consul_prefix)),
    )
