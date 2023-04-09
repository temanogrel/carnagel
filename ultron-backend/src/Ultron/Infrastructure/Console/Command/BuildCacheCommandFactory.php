<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Console\Command;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\CacheService;

final class BuildCacheCommandFactory
{
    public function __invoke(ContainerInterface $container): BuildCacheCommand
    {
        return new BuildCacheCommand(
            $container->get(CacheService::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class)
        );
    }
}
