<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Aphrodite\Site\Factory\Controller;

use Aphrodite\Site\Controller\PostAssociationResourceController;
use Aphrodite\Site\Entity\PostAssociation;
use Aphrodite\Site\Service\PostAssociationService;
use Doctrine\ORM\EntityManager;
use Zend\Mvc\Controller\ControllerManager;

final class PostAssociationResourceControllerFactory
{
    public function __invoke(ControllerManager $manager): PostAssociationResourceController
    {
        $container = $manager->getServiceLocator();

        return new PostAssociationResourceController(
            $container->get(EntityManager::class)->getRepository(PostAssociation::class),
            $container->get(PostAssociationService::class)
        );
    }
}
