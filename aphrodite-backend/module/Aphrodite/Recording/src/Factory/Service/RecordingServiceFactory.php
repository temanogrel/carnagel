<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Service;

use Aphrodite\Recording\Service\RecordingService;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RecordingServiceFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return RecordingService
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        return new RecordingService($serviceLocator->get('Aphrodite\ObjectManager'));
    }
}
