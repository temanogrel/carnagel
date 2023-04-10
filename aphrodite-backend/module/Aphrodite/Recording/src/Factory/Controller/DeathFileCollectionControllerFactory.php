<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Controller;

use Aphrodite\Recording\Controller\DeathFileCollectionController;
use Aphrodite\Recording\Entity\DeathFileEntity;
use Aphrodite\Recording\InputFilter\DeathFileUploadInputFilter;
use Aphrodite\Recording\Repository\DeathFileRepositoryInterface;
use Aphrodite\Recording\Service\DeathFileService;
use Zend\Mvc\Controller\ControllerManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class DeathFileCollectionControllerFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface|ControllerManager $controllerManager
     *
     * @return DeathFileCollectionController
     */
    public function createService(ServiceLocatorInterface $controllerManager)
    {
        $sl = $controllerManager->getServiceLocator();

        /* @var $repository DeathFileRepositoryInterface */
        $repository = $sl->get('Aphrodite\ObjectManager')->getRepository(DeathFileEntity::class);

        /* @var $service DeathFileService */
        $service = $sl->get(DeathFileService::class);

        /* @var $inputFilter DeathFileUploadInputFilter */
        $inputFilter = $sl->get('InputFilterManager')->get(DeathFileUploadInputFilter::class);

        return new DeathFileCollectionController($service, $repository, $inputFilter);
    }
}
