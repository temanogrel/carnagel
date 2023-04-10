<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Factory\Controller;

use Aphrodite\Performer\Controller\PerformerCollectionController;
use Aphrodite\Performer\Controller\PerformerResourceController;
use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Performer\Service\PerformerService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class PerformerResourceControllerFactory implements FactoryInterface
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

        /* @var $service PerformerService */
        $service = $sl->get(PerformerService::class);

        return new PerformerResourceController($repository, $service);
    }
}
