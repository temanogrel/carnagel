<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Application\Factory\Controller;

use Aphrodite\Application\Controller\PrometheusController;
use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Recording\Entity\RecordingEntity;
use Doctrine\ORM\EntityManager;
use Zend\Mvc\Controller\ControllerManager;

final class PrometheusControllerFactory
{
    public function __invoke(ControllerManager $manager): PrometheusController
    {
        $container = $manager->getServiceLocator();

        return new PrometheusController(
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class),
            $container->get(EntityManager::class)->getRepository(AbstractPerformerEntity::class)
        );
    }
}
