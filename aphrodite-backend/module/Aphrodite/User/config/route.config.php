<?php
/**
 *
 *
 *
 */


use Aphrodite\User\Controller\UserCollectionController;use Aphrodite\User\Controller\UserResourceController;

return [
    'users' => [
        'type'          => 'literal',
        'options'       => [
            'route'    => '/users',
            'defaults' => [
                'controller' => UserCollectionController::class
            ]
        ],
        'may_terminate' => true,
        'child_routes'  => [
            'resource' => [
                'type'          => 'segment',
                'options'       => [
                    'route'       => '/:userId',
                    'defaults'    => [
                        'controller' => UserResourceController::class,
                    ],
                    'constraints' => [
                        'userId' => '\d+'
                    ]
                ]
            ]
        ]
    ]
];
