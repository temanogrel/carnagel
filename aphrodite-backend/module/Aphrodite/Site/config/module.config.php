<?php
/**
 *
 *
 *
 */

use Aphrodite\Site\Controller\PostAssociationResourceController;
use Aphrodite\Site\Controller\SiteCollectionController;
use Aphrodite\Site\Controller\SiteResourceController;
use Aphrodite\Site\Factory\Controller\PostAssociationResourceControllerFactory;
use Aphrodite\Site\Factory\Controller\SiteCollectionControllerFactory;
use Aphrodite\Site\Factory\Controller\SiteResourceControllerFactory;
use Aphrodite\Site\Factory\Service\PostAssociationServiceFactory;
use Aphrodite\Site\Factory\Service\SiteServiceFactory;
use Aphrodite\Site\Service\PostAssociationService;
use Aphrodite\Site\Service\SiteService;

return [
    'controllers' => [
        'factories' => [
            SiteResourceController::class            => SiteResourceControllerFactory::class,
            SiteCollectionController::class          => SiteCollectionControllerFactory::class,
            PostAssociationResourceController::class => PostAssociationResourceControllerFactory::class,
        ],
    ],

    'service_manager' => [
        'factories' => [
            SiteService::class            => SiteServiceFactory::class,
            PostAssociationService::class => PostAssociationServiceFactory::class,

        ],
    ],

    'view_manager' => [
        'template_path_stack' => [
            'Aphrodite\Site' => __DIR__ . '/../view/',
        ],
    ],

    'doctrine' => include __DIR__ . '/doctrine.config.php',
    'router'   => [
        'routes' => include __DIR__ . '/route.config.php',
    ],
];
