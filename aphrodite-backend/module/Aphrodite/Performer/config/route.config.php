<?php
/**
 *
 *
 *
 */


use Aphrodite\Performer\Controller\BlacklistCollectionController;use Aphrodite\Performer\Controller\Performer\RecordingCollectionController;use Aphrodite\Performer\Controller\PerformerCollectionController;use Aphrodite\Performer\Controller\PerformerResourceController;

return [
    'blacklist' => [
        'type'          => 'literal',
        'options'       => [
            'route'    => '/blacklist',
            'defaults' => [
                'controller' => BlacklistCollectionController::class
            ]
        ],
    ],

    'performers' => [
        'type'          => 'literal',
        'options'       => [
            'route'    => '/performers',
            'defaults' => [
                'controller' => PerformerCollectionController::class
            ]
        ],
        'may_terminate' => true,
        'child_routes'  => [
            'actions'  => [
                'type'    => 'segment',
                'options' => [
                    'route'       => '/:action',
                    'defaults'    => [
                        'controller' => PerformerCollectionController::class,
                    ],
                    'constraints' => [
                        'action' => 'intersect'
                    ]
                ]
            ],
            'resource' => [
                'type'          => 'segment',
                'options'       => [
                    'route'       => '/:performerId',
                    'defaults'    => [
                        'controller' => PerformerResourceController::class,
                    ],

                    'constraints' => [
                        'performerId' => '\d+|[a-z]+:(.*)'
                    ]
                ],

                'may_terminate' => true,
                'child_routes'  => [
                    'recordings' => [
                        'type'    => 'literal',
                        'options' => [
                            'route'    => '/recordings',
                            'defaults' => [
                                'controller' => RecordingCollectionController::class
                            ]
                        ]
                    ]
                ]
            ]
        ]
    ]
];
