<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Application\Factory;

use Rhubarb\Rhubarb;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RhubarbFactory implements FactoryInterface
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
        return new Rhubarb($serviceLocator->get('config')['aphrodite']['rhubarb']);
    }
}
