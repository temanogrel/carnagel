<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure;

use Interop\Container\ContainerInterface;
use Ultron\Domain\Service\CacheService;

class RepositoryFactoryFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return RepositoryFactory
     */
    public function __invoke(ContainerInterface $container):RepositoryFactory
    {
        return new RepositoryFactory($container->get(CacheService::class));
    }
}
