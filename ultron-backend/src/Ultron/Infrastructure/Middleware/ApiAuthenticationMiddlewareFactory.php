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

final class ApiAuthenticationMiddlewareFactory
{
    public function __invoke(ContainerInterface $container):ApiAuthenticationMiddleware
    {
        return new ApiAuthenticationMiddleware($container->get(Sites::class));
    }
}
