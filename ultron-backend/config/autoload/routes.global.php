<?php

use Ultron\Domain\Action\Api\RecordingCreateAction;
use Ultron\Domain\Action\Api\RecordingCreateActionFactory;
use Ultron\Domain\Action\Api\RecordingDeleteAction;
use Ultron\Domain\Action\Api\RecordingDeleteActionFactory;
use Ultron\Domain\Action\CamgirlGalleryAction;
use Ultron\Domain\Action\CamgirlGalleryActionFactory;
use Ultron\Domain\Action\HomePageAction;
use Ultron\Domain\Action\HomePageActionFactory;
use Ultron\Domain\Action\PerformerRecordingListAction;
use Ultron\Domain\Action\PerformerRecordingListActionFactory;
use Ultron\Domain\Action\RecordingDetailsAction;
use Ultron\Domain\Action\RecordingDetailsActionFactory;
use Ultron\Domain\Action\RecordingSearchAction;
use Ultron\Domain\Action\RecordingSearchActionFactory;
use Ultron\Domain\Action\SiteMapIndexAction;
use Ultron\Domain\Action\SiteMapIndexActionFactory;
use Zend\Expressive\Router\FastRouteRouter;
use Zend\Expressive\Router\RouterInterface;

return [
    'dependencies' => [
        'invokables' => [
            RouterInterface::class => FastRouteRouter::class,
        ],
        'factories'  => [
            HomePageAction::class        => HomePageActionFactory::class,

            // Api actions
            RecordingCreateAction::class => RecordingCreateActionFactory::class,
            RecordingDeleteAction::class => RecordingDeleteActionFactory::class,

            // Misc
            SiteMapIndexAction::class    => SiteMapIndexActionFactory::class,
            CamgirlGalleryAction::class  => CamgirlGalleryActionFactory::class,

            // Domain
            RecordingSearchAction::class        => RecordingSearchActionFactory::class,
            RecordingDetailsAction::class       => RecordingDetailsActionFactory::class,
            PerformerRecordingListAction::class => PerformerRecordingListActionFactory::class,
        ],
    ],

    'routes' => [
        [
            'name'            => 'home',
            'path'            => '/',
            'middleware'      => HomePageAction::class,
            'allowed_methods' => ['GET'],
        ],

        [
            'name'            => 'gallery',
            'path'            => '/camgirl-gallery',
            'middleware'      => CamgirlGalleryAction::class,
            'allowed_methods' => ['GET'],
        ],

        [
            'name'            => 'sitemap',
            'path'            => '/sitemap.xml',
            'middleware'      => SiteMapIndexAction::class,
            'allowed_methods' => ['GET'],
        ],


        [
            'name'            => 'recording.search',
            'path'            => '/search[/{query}]',
            'middleware'      => RecordingSearchAction::class,
            'allowed_methods' => ['GET', 'POST'],
        ],

        [
            'name'            => 'recording.details',
            'path'            => '/{prefix}/r/{slug}',
            'middleware'      => RecordingDetailsAction::class,
            'allowed_methods' => ['GET'],
        ],

        [
            'name'            => 'performer.list-recordings',
            'path'            => '/{prefix}/p/{slug}',
            'middleware'      => PerformerRecordingListAction::class,
            'allowed_methods' => ['GET'],
        ],

        [
            'name'            => 'api.recording.create',
            'path'            => '/api/recordings',
            'middleware'      => RecordingCreateAction::class,
            'allowed_methods' => ['POST'],
        ],

        [
            'name'            => 'api.recording.delete',
            'path'            => '/api/recordings/{id:\d+}',
            'middleware'      => RecordingDeleteAction::class,
            'allowed_methods' => ['DELETE'],
        ],
    ],
];
