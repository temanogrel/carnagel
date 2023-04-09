<?php
/**
 * Retrieve all configuration variables from consul or default values if they don't exist
 */

use Hermes\Options\UpstoreOptions;
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

    $keys = $kv->get('hermes', ['recurse' => true])->json();

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

if (!array_key_exists('elasticsearch', $services)) {
    throw new RuntimeException('No elasticsearch service available');
}

$redisServices = $catalog->service('redis')->json();
$redisHost = $redisServices[0]['ServiceAddress'];
$redisPort = $redisServices[0]['ServicePort'];

$elasticSearchServers = [];
foreach ($catalog->service('elasticsearch')->json() as $instance) {
    $elasticSearchServers[] = [
        'host' => $instance['ServiceAddress'],
        'port' => $instance['ServicePort'],
    ];
}

return [
    'debug' => $getConfig('hermes/debug', 'no') === 'yes',

    'doctrine' => [
        'host'     => $getConfig('hermes/mysql/host', 'mysql'),
        'user'     => $getConfig('hermes/mysql/user', 'hermes'),
        'password' => $getConfig('hermes/mysql/pass', 'hermes'),
        'dbname'   => $getConfig('hermes/mysql/name', 'hermes'),

        'proxy_path' => sys_get_temp_dir(),
    ],

    'redis' => [
        'host' => $redisHost,
        'port' => $redisPort,
    ],

    'api_token' => 'hp1RqDu6EcA8dJCv7MN6',

    UpstoreOptions::class => [
        'apiKey' => 'wSUTFXw2BHdvM',
        'apiUri' => 'https://upstore.net/api/alt/',
    ],
];
