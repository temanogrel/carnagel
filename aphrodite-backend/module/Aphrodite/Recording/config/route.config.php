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

return [

    'recordings' => [
        'type'          => 'literal',
        'options'       => [
            'route'    => '/recordings',
            'defaults' => [
                'controller' => RecordingCollectionController::class,
            ],
        ],
        'may_terminate' => true,
        'child_routes'  => [
            'resource' => [
                'type'    => 'segment',
                'options' => [
                    'route'    => '/:recordingId',
                    'defaults' => [
                        'controller' => RecordingResourceController::class,
                    ],
                ],

                'may_terminate' => true,
                'child_routes'  => [
                    'posts' => [
                        'type'    => 'literal',
                        'options' => [
                            'route'    => '/posts',
                            'defaults' => [
                                'controller' => PostAssociationCollectionController::class,
                            ],
                        ],
                    ],
                ],
            ],
        ],
    ],


    'death-files' => [
        'type'          => 'literal',
        'options'       => [
            'route'    => '/death-files',
            'defaults' => [
                'controller' => DeathFileCollectionController::class,
            ],
        ],
        'may_terminate' => true,
        'child_routes'  => [
            'resource' => [
                'type'    => 'segment',
                'options' => [
                    'route'    => '/:id',
                    'defaults' => [
                        'controller' => DeathFileResourceController::class,
                    ],
                ],

                'may_terminate' => true,
                'child_routes'  => [
                    'urls' => [
                        'type'    => 'literal',
                        'options' => [
                            'route'    => '/urls',
                            'defaults' => [
                                'controller' => UrlEntryCollectionController::class,
                            ],
                        ],
                    ],
                ],
            ],
        ],
    ],

    'url-entries' => [
        'type'    => 'literal',
        'options' => [
            'route'    => '/urls',
            'defaults' => [
                'controller' => StandaloneUrlEntryCollectionController::class,
            ],
        ],

        'may_terminate' => true,
        'child_routes'  => [
            'resource' => [
                'type'    => 'segment',
                'options' => [
                    'route'    => '/:id',
                    'defaults' => [
                        'controller' => UrlEntryResourceController::class,
                    ],
                ],
            ],
        ],
    ],
];
