<?php
use Doctrine\Common\Persistence\ObjectManager;

/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Factory\Service;

use Aphrodite\Site\Service\PostAssociationService;
use Doctrine\Common\Persistence\ObjectManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class PostAssociationServiceFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return PostAssociationService
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        /* @var $objectManager ObjectManager */
        $objectManager = $serviceLocator->get('Aphrodite\ObjectManager');

        return new PostAssociationService($objectManager);
    }
}
