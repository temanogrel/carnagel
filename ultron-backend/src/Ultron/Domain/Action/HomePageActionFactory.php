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
use Ultron\Infrastructure\Service\PageCacheService;
use Zend\Expressive\Template\TemplateRendererInterface;

class HomePageActionFactory
{
    public function __invoke(ContainerInterface $container)
    {
        $site       = $container->get(Sites::class)->getCurrentSite();
        $template   = $container->get(TemplateRendererInterface::class);
        $repository = $container->get(EntityManager::class)->getRepository(RecordingEntity::class);
        $pageCache  = $container->get(PageCacheService::class);

        return new HomePageAction($site, $template, $repository, $pageCache);
    }
}
