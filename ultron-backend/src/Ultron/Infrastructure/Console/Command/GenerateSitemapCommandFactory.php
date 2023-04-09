<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Console\Command;

use Interop\Container\ContainerInterface;
use Ultron\Infrastructure\Service\SitemapService;

final class GenerateSitemapCommandFactory
{
    public function __invoke(ContainerInterface $container): GenerateSitemapCommand
    {
        return new GenerateSitemapCommand($container->get(SitemapService::class));
    }
}
