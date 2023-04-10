<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Factory;

use Aphrodite\Blocktrail\Options\BlocktrailOptions;
use Blocktrail\SDK\BlocktrailSDK;
use Interop\Container\ContainerInterface;

final class BlocktrailSDKFactory
{
    public function __invoke(ContainerInterface $container): BlocktrailSDK
    {
        /* @var BlocktrailOptions $options */
        $options = $container->get(BlocktrailOptions::class);

        return new BlocktrailSDK(
            $options->getApiKey(),
            $options->getApiSecret(),
            'BTC',
            $options->getTestNet()
        );
    }
}
