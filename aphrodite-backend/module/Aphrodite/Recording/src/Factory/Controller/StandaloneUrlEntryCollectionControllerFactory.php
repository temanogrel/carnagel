<?php

namespace Aphrodite\Recording\Factory\Controller;

use Aphrodite\Recording\Controller\StandaloneUrlEntryCollectionController;
use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\Repository\DeathFile\UrlRepository;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class StandaloneUrlEntryCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ControllerManager|ServiceLocatorInterface $controllerManager
     *
     * @return StandaloneUrlEntryCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository UrlRepository */
        $repository = $sl
                        ->get('Aphrodite\ObjectManager')
                        ->getRepository(UrlEntry::class);

        return new StandaloneUrlEntryCollectionController($repository);
    }
}
