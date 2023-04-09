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
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\RecordingService;

class RecordingDeleteActionFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return RecordingDeleteAction
     */
    public function __invoke(ContainerInterface $container):RecordingDeleteAction
    {
        return new RecordingDeleteAction(
            $container->get(RecordingService::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class)
        );
    }
}
