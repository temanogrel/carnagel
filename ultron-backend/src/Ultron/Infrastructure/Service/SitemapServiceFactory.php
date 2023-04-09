<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Sites;
use Zend\Expressive\Helper\UrlHelper;

final class SitemapServiceFactory
{
    public function __invoke(ContainerInterface $container): SitemapService
    {
        /* @var EntityManager $entityManager */
        $entityManager = $container->get(EntityManager::class);

        return new SitemapService(
            $entityManager->getRepository(PerformerEntity::class),
            $entityManager->getRepository(RecordingEntity::class),
            $container->get(Sites::class),
            $container->get(UrlHelper::class),
            $entityManager
        );
    }
}
