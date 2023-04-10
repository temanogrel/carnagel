<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Factory\Service;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Repository\PerformerRepository;
use Aphrodite\Performer\Service\IntersectionService;
use Doctrine\Common\Persistence\ObjectManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class IntersectionServiceFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return IntersectionService
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        /* @var $objectManager ObjectManager */
        $objectManager = $serviceLocator->get('Aphrodite\ObjectManager');

        /* @var $repository PerformerRepository */
        $repository = $objectManager->getRepository(AbstractPerformerEntity::class);

        return new IntersectionService($objectManager, $repository);
    }
}
