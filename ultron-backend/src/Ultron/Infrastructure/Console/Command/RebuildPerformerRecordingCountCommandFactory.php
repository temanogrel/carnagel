<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Console\Command;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Entity\RecordingEntity;

class RebuildPerformerRecordingCountCommandFactory
{
    public function __invoke(ContainerInterface $container): RebuildPerformerRecordingCountCommand
    {
        /* @var EntityManager $entityManager */
        $entityManager = $container->get(EntityManager::class);

        return new RebuildPerformerRecordingCountCommand(
            $entityManager->getRepository(RecordingEntity::class),
            $entityManager->getRepository(PerformerEntity::class),
            $entityManager
        );
    }
}
