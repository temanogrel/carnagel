<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller\DeathFile;

use Aphrodite\Recording\Controller\DeathFile\UrlEntryCollectionController;
use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\Entity\DeathFileEntity;
use Aphrodite\Recording\Repository\DeathFile\UrlRepositoryInterface;
use Aphrodite\Recording\Repository\DeathFileRepositoryInterface;
use Aphrodite\Recording\Service\DeathFile\UrlService;
use Aphrodite\Recording\Service\DeathFile\UrlServiceInterface;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class UrlEntryCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return UrlEntryCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $urlRepository UrlRepositoryInterface */
        $urlRepository = $sl->get('Aphrodite\ObjectManager')->getRepository(UrlEntry::class);

        /* @var $deathFileRepository DeathFileRepositoryInterface */
        $deathFileRepository = $sl->get('Aphrodite\ObjectManager')->getRepository(DeathFileEntity::class);

        /* @var $service UrlServiceInterface */
        $service = $sl->get(UrlService::class);

        return new UrlEntryCollectionController($service, $urlRepository, $deathFileRepository);
    }
}
