<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Service;

use Interop\Container\ContainerInterface;
use Redis;
use Ultron\Domain\Sites;

class CacheServiceFactory
{
    public function __invoke(ContainerInterface $container)
    {
        return new CacheService(
            $container->get(Redis::class),
            $container->get(Sites::class)
        );
    }
}
