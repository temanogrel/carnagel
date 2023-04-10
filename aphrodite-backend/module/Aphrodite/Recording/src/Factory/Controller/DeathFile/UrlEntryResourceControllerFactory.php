<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller\DeathFile;

use Aphrodite\Recording\Controller\DeathFile\UrlEntryResourceController;
use Aphrodite\Recording\Controller\DeathFileResourceController;
use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\Repository\DeathFile\UrlRepositoryInterface;
use Aphrodite\Recording\Service\DeathFile\UrlService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class UrlEntryResourceControllerFactory implements FactoryInterface
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

        /* @var $repository UrlRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(UrlEntry::class);

        /* @var $service UrlService */
        $service = $sl->get(UrlService::class);

        return new UrlEntryResourceController($service, $repository);
    }
}
