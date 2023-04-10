<?php
/**
 *
 *
 *
 */


use Aphrodite\Site\Controller\PostAssociationResourceController;
use Aphrodite\Site\Controller\SiteCollectionController;
use Aphrodite\Site\Controller\SiteResourceController;

return [

    'sites' => [
        'type'    => 'literal',
        'options' => [
            'route'    => '/sites',
            'defaults' => [
                'controller' => SiteCollectionController::class,
            ],
        ],

        'may_terminate' => true,
        'child_routes'  => [
            
            'resource' => [
                'type'    => 'segment',
                'options' => [
                    'route'       => '/:siteId',
                    'defaults'    => [
                        'controller' => SiteResourceController::class,
                    ],
                    'constraints' => [
                        'siteId' => '\d+',
                    ],
                ],
            ],
        ],
    ],

    'posts' => [
        'type'    => 'segment',
        'options' => [

            'route'    => '/posts/:id',
            'defaults' => [
                'controller' => PostAssociationResourceController::class,
            ],
        ],
    ],
];
