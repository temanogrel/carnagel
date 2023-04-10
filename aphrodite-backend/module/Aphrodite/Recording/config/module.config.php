<?php
/**
 *
 *
 *
 */

use Aphrodite\Recording\Controller\DeathFile\UrlEntryCollectionController;
use Aphrodite\Recording\Controller\DeathFile\UrlEntryResourceController;
use Aphrodite\Recording\Controller\DeathFileCollectionController;
use Aphrodite\Recording\Controller\DeathFileResourceController;
use Aphrodite\Recording\Controller\Recording\PostAssociationCollectionController;
use Aphrodite\Recording\Controller\RecordingCollectionController;
use Aphrodite\Recording\Controller\RecordingResourceController;
use Aphrodite\Recording\Controller\StandaloneUrlEntryCollectionController;
use Aphrodite\Recording\Factory\Controller\DeathFile\UrlEntryCollectionControllerFactory;
use Aphrodite\Recording\Factory\Controller\DeathFile\UrlEntryResourceControllerFactory;
use Aphrodite\Recording\Factory\Controller\DeathFileCollectionControllerFactory;
use Aphrodite\Recording\Factory\Controller\DeathFileResourceControllerFactory;
use Aphrodite\Recording\Factory\Controller\Recording\PostAssociationCollectionControllerFactory;
use Aphrodite\Recording\Factory\Controller\RecordingCollectionControllerFactory;
use Aphrodite\Recording\Factory\Controller\RecordingResourceControllerFactory;
use Aphrodite\Recording\Factory\Controller\StandaloneUrlEntryCollectionControllerFactory;
use Aphrodite\Recording\Factory\Service\DeathFile\UrlServiceFactory;
use Aphrodite\Recording\Factory\Service\DeathFileServiceFactory;
use Aphrodite\Recording\Factory\Service\RecordingServiceFactory;
use Aphrodite\Recording\Factory\Validator\RecordingExistsFactory;
use Aphrodite\Recording\Factory\Validator\UrlDoesNotExistFactory;
use Aphrodite\Recording\Service\DeathFile\UrlService;
use Aphrodite\Recording\Service\DeathFileService;
use Aphrodite\Recording\Service\RecordingService;
use Aphrodite\Recording\Validator\RecordingExists;
use Aphrodite\Recording\Validator\UrlDoesNotExist;

return [
    'controllers' => [
        'factories' => [
            RecordingResourceController::class         => RecordingResourceControllerFactory::class,
            RecordingCollectionController::class       => RecordingCollectionControllerFactory::class,
            PostAssociationCollectionController::class => PostAssociationCollectionControllerFactory::class,

            // Death files
            DeathFileResourceController::class         => DeathFileResourceControllerFactory::class,
            DeathFileCollectionController::class       => DeathFileCollectionControllerFactory::class,

            UrlEntryResourceController::class   => UrlEntryResourceControllerFactory::class,
            UrlEntryCollectionController::class => UrlEntryCollectionControllerFactory::class,

            StandaloneUrlEntryCollectionController::class => StandaloneUrlEntryCollectionControllerFactory::class,
        ],
    ],

    'service_manager' => [
        'factories' => [
            UrlService::class       => UrlServiceFactory::class,
            RecordingService::class => RecordingServiceFactory::class,
            DeathFileService::class => DeathFileServiceFactory::class,
        ],
    ],

    'validators' => [
        'factories' => [
            RecordingExists::class => RecordingExistsFactory::class,
            UrlDoesNotExist::class => UrlDoesNotExistFactory::class,
        ],
    ],

    'view_manager' => [
        'template_path_stack' => [
            'Aphrodite\Recording' => __DIR__ . '/../view/',
        ],
    ],

    'doctrine' => include __DIR__ . '/doctrine.config.php',
    'router'   => [
        'routes' => include __DIR__ . '/route.config.php',
    ],
];
