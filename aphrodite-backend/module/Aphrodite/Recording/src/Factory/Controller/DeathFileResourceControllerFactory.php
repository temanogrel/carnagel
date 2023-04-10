<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller;

use Aphrodite\Recording\Controller\DeathFileResourceController;
use Aphrodite\Recording\Entity\DeathFileEntity;
use Aphrodite\Recording\Repository\DeathFileRepositoryInterface;
use Aphrodite\Recording\Service\DeathFileService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class DeathFileResourceControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return DeathFileResourceController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository DeathFileRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(DeathFileEntity::class);

        /* @var $service DeathFileService */
        $service = $sl->get(DeathFileService::class);

        return new DeathFileResourceController($service, $repository);
    }
}
