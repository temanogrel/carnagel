<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;

use Interop\Container\ContainerInterface;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Service\SitemapService;
use Zend\Expressive\Template\TemplateRendererInterface;

final class SiteMapIndexActionFactory
{
    public function __invoke(ContainerInterface $container): SiteMapIndexAction
    {
        return new SiteMapIndexAction(
            $container->get(Sites::class)->getCurrentSite(),
            $container->get(SitemapService::class),
            $container->get(TemplateRendererInterface::class)
        );
    }
}
