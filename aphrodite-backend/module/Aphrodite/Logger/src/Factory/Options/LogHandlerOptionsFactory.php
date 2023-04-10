<?php
/**
 *
 * 
 */

declare(strict_types=1);

namespace Aphrodite\Logger\Factory\Options;

use Aphrodite\Logger\Options\LogHandlerOptions;
use Psr\Container\ContainerInterface;

final class LogHandlerOptionsFactory
{
    /**
     * @param ContainerInterface $container
     * @return LogHandlerOptions
     */
    public function __invoke(ContainerInterface $container): LogHandlerOptions
    {
        $config = $container->get('Config')['aphrodite']['options'];
        $config = $config[LogHandlerOptions::class] ?? [];

        return new LogHandlerOptions($config);
    }
}
