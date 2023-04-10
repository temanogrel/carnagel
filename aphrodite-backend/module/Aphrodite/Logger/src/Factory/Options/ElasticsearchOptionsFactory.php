<?php
/**
 *
 * 
 */

declare(strict_types=1);

namespace Aphrodite\Logger\Factory\Options;

use Aphrodite\Logger\Options\ElasticsearchOptions;
use Psr\Container\ContainerInterface;

final class ElasticsearchOptionsFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return ElasticsearchOptions
     *
     * @throws \Psr\Container\ContainerExceptionInterface
     */
    public function __invoke(ContainerInterface $container): ElasticsearchOptions
    {
        $config = $container->get('Config')['aphrodite']['options'];
        $config = $config[ElasticsearchOptions::class] ?? [];

        return new ElasticsearchOptions($config);
    }
}
