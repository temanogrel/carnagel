<?php
/**
 *
 *
 *
 */

use Aphrodite\Performer\Controller\BlacklistCollectionController;use Aphrodite\Performer\Controller\Performer\RecordingCollectionController;use Aphrodite\Performer\Controller\PerformerCollectionController;use Aphrodite\Performer\Controller\PerformerResourceController;use Aphrodite\Performer\Factory\Controller\BlacklistCollectionControllerFactory;use Aphrodite\Performer\Factory\Controller\Performer\RecordingCollectionControllerFactory;use Aphrodite\Performer\Factory\Controller\PerformerCollectionControllerFactory;use Aphrodite\Performer\Factory\Controller\PerformerResourceControllerFactory;use Aphrodite\Performer\Factory\Service\IntersectionServiceFactory;use Aphrodite\Performer\Factory\Service\PerformerServiceFactory;use Aphrodite\Performer\Service\IntersectionService;use Aphrodite\Performer\Service\PerformerService;

return [
    'controllers'     => [
        'factories' => [
            BlacklistCollectionController::class => BlacklistCollectionControllerFactory::class,
            PerformerResourceController::class   => PerformerResourceControllerFactory::class,
            PerformerCollectionController::class => PerformerCollectionControllerFactory::class,
            RecordingCollectionController::class => RecordingCollectionControllerFactory::class
        ]
    ],
    'service_manager' => [
        'factories' => [
            IntersectionService::class => IntersectionServiceFactory::class,
            PerformerService::class    => PerformerServiceFactory::class
        ]
    ],
    'view_manager'    => [
        'template_path_stack' => [
            'Aphrodite\Performer' => __DIR__ . '/../view/'
        ]
    ],
    'router'          => [
        'routes' => include __DIR__ . '/route.config.php'
    ],
    'doctrine'        => include __DIR__ . '/doctrine.config.php'
];
