<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Factory\Options;

use Aphrodite\Blocktrail\Options\BlocktrailOptions;
use Psr\Container\ContainerInterface;

final class BlocktrailOptionsFactory
{
    /**
     * @param ContainerInterface $container
     * @return BlocktrailOptions
     */
    public function __invoke(ContainerInterface $container): BlocktrailOptions
    {
        /* @var array $config */
        $config = $container->get('config');

        return new BlocktrailOptions($config['aphrodite']['options'][BlocktrailOptions::class]);
    }
}
