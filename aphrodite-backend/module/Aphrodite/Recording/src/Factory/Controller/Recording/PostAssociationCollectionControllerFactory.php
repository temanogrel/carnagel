<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller\Recording;

use Aphrodite\Recording\Controller\Recording\PostAssociationCollectionController;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Site\Entity\Site;
use Aphrodite\Site\Repository\SiteRepositoryInterface;
use Aphrodite\Site\Service\PostAssociationService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class PostAssociationCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return PostAssociationCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository RecordingRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(RecordingEntity::class);

        /* @var $wpRepository SiteRepositoryInterface */
        $wpRepository = $sl->get('Aphrodite\ObjectManager')->getRepository(Site::class);

        /* @var $service PostAssociationService */
        $service = $sl->get(PostAssociationService::class);

        return new PostAssociationCollectionController($repository, $wpRepository, $service);
    }
}
