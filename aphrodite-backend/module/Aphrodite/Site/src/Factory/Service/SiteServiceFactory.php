<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Factory\Service;

use Aphrodite\Site\Service\SiteService;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class SiteServiceFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return mixed
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        return new SiteService($serviceLocator->get('Aphrodite\ObjectManager'));
    }
}
