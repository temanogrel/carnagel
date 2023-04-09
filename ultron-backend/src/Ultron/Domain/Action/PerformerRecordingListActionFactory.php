<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Sites;
use Zend\Expressive\Template\TemplateRendererInterface;

class PerformerRecordingListActionFactory
{
    public function __invoke(ContainerInterface $container):PerformerRecordingListAction
    {
        return new PerformerRecordingListAction(
            $container->get(Sites::class)->getCurrentSite(),
            $container->get(TemplateRendererInterface::class),
            $container->get(EntityManager::class)->getRepository(PerformerEntity::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class)
        );
    }
}
