<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Factory\Controller;

use Aphrodite\Performer\Controller\PerformerCollectionController;
use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Performer\Service\IntersectionService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class PerformerCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return PerformerCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository PerformerRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(AbstractPerformerEntity::class);

        /* @var $service IntersectionService */
        $service = $sl->get(IntersectionService::class);

        return new PerformerCollectionController($repository, $service);
    }
}
