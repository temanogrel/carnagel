<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Redis;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\CacheService;
use Ultron\Domain\Sites;

final class PageCacheServiceFactory
{
    public function __invoke(ContainerInterface $container): PageCacheService
    {
        return new PageCacheService(
            $container->get(Redis::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class),
            $container->get(CacheService::class),
            $container->get(Sites::class)
        );
    }
}
