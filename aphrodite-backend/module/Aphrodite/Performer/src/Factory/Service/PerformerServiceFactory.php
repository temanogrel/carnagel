<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Factory\Service;

use Aphrodite\Performer\Service\PerformerService;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class PerformerServiceFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return PerformerService
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        return new PerformerService($serviceLocator->get('Aphrodite\ObjectManager'));
    }
}
