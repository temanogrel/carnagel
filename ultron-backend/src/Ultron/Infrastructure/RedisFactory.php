<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure;

use Interop\Container\ContainerInterface;
use Redis;

class RedisFactory
{
    public function __invoke(ContainerInterface $container, string $requestedName): Redis
    {
        $redis = new Redis();
        $redis->connect(
            $container->get('config')['ultron']['redis']['host'],
            $container->get('config')['ultron']['redis']['port']
        );

        if ($requestedName !== DoctrineRedisCache::class) {
            $redis->setOption(Redis::OPT_SERIALIZER, (string) Redis::SERIALIZER_NONE);
            $redis->setOption(Redis::OPT_PREFIX, 'doctrine');
        }

        return $redis;
    }
}
