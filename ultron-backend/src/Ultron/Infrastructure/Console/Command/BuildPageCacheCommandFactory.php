<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Console\Command;

use Interop\Container\ContainerInterface;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Service\PageCacheService;

final class BuildPageCacheCommandFactory
{
    public function __invoke(ContainerInterface $container): BuildPageCacheCommand
    {
        return new BuildPageCacheCommand(
            $container->get(PageCacheService::class),
            $container->get(Sites::class)
        );
    }
}
