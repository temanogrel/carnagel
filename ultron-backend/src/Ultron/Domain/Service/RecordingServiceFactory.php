<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Service;

use Cocur\Slugify\Slugify;
use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Infrastructure\Service\PageCacheService;

class RecordingServiceFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return RecordingService
     */
    public function __invoke(ContainerInterface $container):RecordingService
    {
        return new RecordingService(
            $container->get(CacheService::class),
            $container->get(EntityManager::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class),
            $container->get(Slugify::class),
            $container->get(PageCacheService::class)
        );
    }
}
