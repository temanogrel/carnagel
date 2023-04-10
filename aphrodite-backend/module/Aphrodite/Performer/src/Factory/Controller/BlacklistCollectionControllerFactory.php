<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Factory\Controller;

use Aphrodite\Performer\Controller\BlacklistCollectionController;
use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class BlacklistCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return BlacklistCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository PerformerRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(AbstractPerformerEntity::class);

        return new BlacklistCollectionController($repository);
    }
}
