<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Middleware;

use Interop\Container\ContainerInterface;
use Ultron\Domain\Sites;

class SiteSelectionMiddlewareFactory
{
    public function __invoke(ContainerInterface $container)
    {
        return new SiteSelectionMiddleware($container->get(Sites::class));
    }
}
