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
use Zend\Expressive\Template\TemplateRendererInterface;

class RecordingSearchActionFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return RecordingSearchAction
     */
    public function __invoke(ContainerInterface $container):RecordingSearchAction
    {
        return new RecordingSearchAction(
            $container->get(UrlHelper::class),
            $container->get(Sites::class)->getCurrentSite(),
            $container->get(TemplateRendererInterface::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class)
        );
    }
}
