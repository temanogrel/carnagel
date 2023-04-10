<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Factory\Controller;

use Aphrodite\Site\Controller\SiteResourceController;
use Aphrodite\Site\Entity\Site;
use Aphrodite\Site\Repository\SiteRepositoryInterface;
use Aphrodite\Site\Service\SiteService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class SiteResourceControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return SiteResourceController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository SiteRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(Site::class);

        /* @var $service SiteService */
        $service = $sl->get(SiteService::class);

        return new SiteResourceController($repository, $service);
    }
}
