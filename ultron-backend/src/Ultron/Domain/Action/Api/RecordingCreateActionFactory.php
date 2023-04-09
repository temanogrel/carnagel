<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action\Api;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\PerformerService;
use Ultron\Domain\Service\RecordingService;

class RecordingCreateActionFactory
{
    public function __invoke(ContainerInterface $container)
    {
        return new RecordingCreateAction(
            $container->get(RecordingService::class),
            $container->get(PerformerService::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class),
            $container->get(EntityManager::class)->getRepository(PerformerEntity::class)
        );
    }
}
