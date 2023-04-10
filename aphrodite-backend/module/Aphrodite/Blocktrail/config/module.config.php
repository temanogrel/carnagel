<?php
/**
 *
 *
 */

use Aphrodite\Blocktrail\Controller\RpcController;
use Aphrodite\Blocktrail\Factory\BlocktrailSDKFactory;
use Aphrodite\Blocktrail\Factory\Controller\RpcControllerFactory;
use Aphrodite\Blocktrail\Factory\Options\BlocktrailOptionsFactory;
use Aphrodite\Blocktrail\Factory\Service\BlocktrailServiceFactory;
use Aphrodite\Blocktrail\Options\BlocktrailOptions;
use Aphrodite\Blocktrail\Service\BlocktrailService;
use Blocktrail\SDK\BlocktrailSDK;

return [
    'service_manager' => [
        'factories' => [
            BlocktrailSDK::class => BlocktrailSDKFactory::class,

            BlocktrailOptions::class => BlocktrailOptionsFactory::class,

            BlocktrailService::class => BlocktrailServiceFactory::class,
        ],
    ],

    'controllers' => [
        'factories' => [
            RpcController::class => RpcControllerFactory::class,
        ],
    ],

    'router' => [
        'routes' => include __DIR__ . '/route.config.php',
    ],
];
