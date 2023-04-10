<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Factory\Service;

use Aphrodite\Blocktrail\Options\BlocktrailOptions;
use Aphrodite\Blocktrail\Service\BlocktrailService;
use Blocktrail\SDK\BlocktrailSDK;
use Psr\Container\ContainerInterface;

final class BlocktrailServiceFactory
{
    /**
     * @param ContainerInterface $container
     * @return BlocktrailService
     */
    public function __invoke(ContainerInterface $container): BlocktrailService
    {
        return new BlocktrailService(
            $container->get(BlocktrailSDK::class),
            $container->get(BlocktrailOptions::class)
        );
    }
}
