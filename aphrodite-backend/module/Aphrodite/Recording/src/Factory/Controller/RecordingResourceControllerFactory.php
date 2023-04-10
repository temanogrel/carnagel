<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller;

use Aphrodite\Recording\Controller\RecordingCollectionController;
use Aphrodite\Recording\Controller\RecordingResourceController;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Recording\Service\RecordingService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RecordingResourceControllerFactory implements FactoryInterface
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

        /* @var $service RecordingService */
        $service = $sl->get(RecordingService::class);

        return new RecordingResourceController($repository, $service);
    }
}
