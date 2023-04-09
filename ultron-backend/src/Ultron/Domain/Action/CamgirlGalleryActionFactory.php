<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Sites;
use Zend\Expressive\Helper\UrlHelper;

final class CamgirlGalleryActionFactory
{
    public function __invoke(ContainerInterface $container):CamgirlGalleryAction
    {
        return new CamgirlGalleryAction(
            $container->get(Sites::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class),
            $container->get(UrlHelper::class)
        );
    }
}
