<?php
use Aphrodite\Recording\Service\RecordingServiceInterface;

/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller;

use Aphrodite\Recording\Controller\RecordingCollectionController;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Recording\Service\IntersectionService;
use Aphrodite\Recording\Service\RecordingService;
use Aphrodite\Recording\Service\RecordingServiceInterface;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RecordingCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return RecordingCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository RecordingRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(RecordingEntity::class);

        /* @var $service RecordingServiceInterface */
        $service = $sl->get(RecordingService::class);

        return new RecordingCollectionController($repository, $service);
    }
}
