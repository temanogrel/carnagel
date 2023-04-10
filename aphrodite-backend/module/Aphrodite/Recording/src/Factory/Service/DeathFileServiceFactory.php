<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Service;

use Aphrodite\Recording\Service\DeathFileService;
use Doctrine\Common\Persistence\ObjectManager;
use League\Flysystem\Filesystem;
use Rhubarb\Rhubarb;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class DeathFileServiceFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return DeathFileService
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        /* @var $objectManager ObjectManager */
        $objectManager = $serviceLocator->get('Aphrodite\ObjectManager');

        /* @var $filesystem Filesystem */
        $filesystem = $serviceLocator->get('BsbFlysystemManager')->get('default');

        /* @var $rhubarb Rhubarb */
        $rhubarb = $serviceLocator->get(Rhubarb::class);

        return new DeathFileService($objectManager, $filesystem, $rhubarb);
    }
}
