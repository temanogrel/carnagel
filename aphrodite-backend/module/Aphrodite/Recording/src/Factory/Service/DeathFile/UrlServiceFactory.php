<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Service\DeathFile;

use Aphrodite\Recording\Service\DeathFile\UrlService;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class UrlServiceFactory implements FactoryInterface
{

    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return UrlService
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        return new UrlService($serviceLocator->get('Aphrodite\ObjectManager'));
    }
}
