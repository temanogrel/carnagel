<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Factory\Controller\Performer;

use Aphrodite\Performer\Controller\Performer\RecordingCollectionController;
use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Recording\Service\RecordingService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RecordingCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return RecordingCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $performerRepository PerformerRepositoryInterface */
        $performerRepository = $sl->get('Aphrodite\ObjectManager')->getRepository(AbstractPerformerEntity::class);

        /* @var $recordingRepository RecordingRepositoryInterface */
        $recordingRepository = $sl->get('Aphrodite\ObjectManager')->getRepository(RecordingEntity::class);

        /* @var $recordingService RecordingService */
        $recordingService = $sl->get(RecordingService::class);

        return new RecordingCollectionController($recordingService, $performerRepository, $recordingRepository);
    }
}
