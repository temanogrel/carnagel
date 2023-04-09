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
use Ultron\Domain\Entity\PerformerEntity;

class PerformerServiceFactory
{
    public function __invoke(ContainerInterface $container):PerformerService
    {
        return new PerformerService(
            $container->get(EntityManager::class),
            $container->get(EntityManager::class)->getRepository(PerformerEntity::class),
            $container->get(Slugify::class)
        );
    }
}
