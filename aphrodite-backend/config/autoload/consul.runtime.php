<?php
/**
 * Retrieve all configuration variables from consul or default values if they don't exist
 */

use Aphrodite\Application\Options\RedisOptions;
use Aphrodite\Blocktrail\Options\BlocktrailOptions;
use Aphrodite\Logger\Adapter\ElasticsearchAdapter;
use Aphrodite\Logger\Options\ElasticsearchOptions;
use Aphrodite\Logger\Options\LogHandlerOptions;
use SensioLabs\Consul\Exception\ClientException;
use SensioLabs\Consul\ServiceFactory;
use SensioLabs\Consul\Services\Catalog;
use SensioLabs\Consul\Services\KV;

$sf = new ServiceFactory([
    'base_uri' => 'http://consul:8500',
]);

/* @var $kv KV */
$kv = $sf->get('kv');

/* @var $catalog Catalog */
$catalog = $sf->get('catalog');

$config = [];

try {

    $keys = array_merge(
        $kv->get('aphrodite', ['recurse' => true])->json(),
        $kv->get('blocktrail', ['recurse' => true])->json()
    );

    foreach ($keys as $key) {
        $config[$key['Key']] = base64_decode($key['Value']);
    }
} catch (ClientException $e) {
    if ($e->getCode() !== 404) {
        throw $e;
    }
}

/**
 * Parse the config tree and extract the config value, provide the default if it does not exist
 *
 * @param string $name
 * @param string $defaultValue
 *
 * @return string
 */
$getConfig = function (string $name, string $defaultValue) use ($config): string {
    if (!isset($config[$name])) {
        return $defaultValue;
    }

    return $config[$name];
};

$services = $catalog->services()->json();

if (!array_key_exists('redis', $services)) {
    throw new RuntimeException('No redis service available');
}

$redisServices = $catalog->service('redis')->json();
$redisHost     = $redisServices[0]['ServiceAddress'];
$redisPort     = $redisServices[0]['ServicePort'];

$config = [
    'doctrine' => [
        'connection' => [
            'orm_default' => [
                'params' => [
                    'host'     => $getConfig('aphrodite/mysql/host', 'mysql'),
                    'user'     => $getConfig('aphrodite/mysql/user', 'aphrodite'),
                    'password' => $getConfig('aphrodite/mysql/pass', 'aphrodite'),
                    'dbname'   => $getConfig('aphrodite/mysql/name', 'aphrodite'),
                ],
            ],

        ],

        'configuration' => [
            'orm_default' => [
                'metadata_cache' => $getConfig('aphrodite/cache-driver', 'array'),
                'query_cache'    => $getConfig('aphrodite/cache-driver', 'array'),
                'result_cache'   => $getConfig('aphrodite/cache-driver', 'array'),
            ],
        ],
    ],

    'aphrodite' => [

        'serverAccessToken' => $getConfig('aphrodite/access-token', 'helloWorld'),

        'rhubarb' => [
            'broker' => [
                'type'    => 'Amqp',
                'options' => [
                    'connection' => $getConfig('aphrodite/rabbitmq-dsn', 'amqp://rtmp:rtmp@rabbitmq:5672/rtmp'),
                ],
            ],
        ],

        'options' => [
            RedisOptions::class => [
                'host' => $redisHost,
                'port' => $redisPort,
            ],

            BlocktrailOptions::class => [
                'apiKey'                                => $getConfig('blocktrail/apiKey', ''),
                'apiSecret'                             => $getConfig('blocktrail/apiSecret', ''),
                'walletId'                              => $getConfig('blocktrail/walletId', ''),
                'walletPassword'                        => $getConfig('blocktrail/walletPassword', ''),
                'numberOfConfirmationsToTriggerWebhook' => (int) $getConfig('blocktrail/numberOfConfirmationsToTriggerWebhook', ''),
                'testNet'                               => (bool) $getConfig('blocktrail/testNet', '0'),
                'testAddress'                           => $getConfig('blocktrail/testAddress', ''),
            ],

            LogHandlerOptions::class => [
                'host'            => gethostname(),
                'environment'     => 'prod',
                'debug'           => false,
                'adapters'        => [
                    ElasticsearchAdapter::class,
                ],
                'alwaysLogRoutes' => [],
            ],

            ElasticsearchOptions::class => [
                'host'     => '10g.es1.vee.bz',
                'port'     => 9200,
                'username' => '',
                'password' => '',
            ],
        ],
    ],
];

if ($getConfig('aphrodite/cache-driver', 'array') === 'redis') {
    $config['doctrine']['cache'] = [
        'redis' => [
            'namespace' => 'aphrodite',
            'instance'  => Redis::class,
        ],
    ];
}

return $config;
