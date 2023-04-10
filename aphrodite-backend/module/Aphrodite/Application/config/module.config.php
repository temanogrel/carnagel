<?php
/**
 * Zend Framework (http://framework.zend.com/)
 *
 * @link      http://github.com/zendframework/ZendSkeletonApplication for the canonical source repository
 * @copyright Copyright (c) 2005-2015 Zend Technologies USA Inc. (http://www.zend.com)
 * @license   http://framework.zend.com/license/new-bsd New BSD License
 */

use Aphrodite\Application\Controller\PrometheusController;
use Aphrodite\Application\Factory\Controller\PrometheusControllerFactory;
use Aphrodite\Application\Factory\Options\RedisOptionsFactory;
use Aphrodite\Application\Factory\RedisFactory;
use Aphrodite\Application\Factory\RhubarbFactory;
use Aphrodite\Application\Options\RedisOptions;
use Rhubarb\Rhubarb;
use Zend\Mvc\Router\Http\Literal;

return [
    'service_manager' => [
        'factories' => [
            Redis::class        => RedisFactory::class,
            RedisOptions::class => RedisOptionsFactory::class,
            Rhubarb::class      => RhubarbFactory::class,
        ],
    ],

    'controllers' => [
        'factories' => [
            PrometheusController::class => PrometheusControllerFactory::class,
        ]
    ],

    'view_manager' => [
        'display_not_found_reason' => true,
        'display_exceptions'       => true,
        'doctype'                  => 'HTML5',
        'not_found_template'       => 'error/404',
        'exception_template'       => 'error/index',
        'template_map'             => [
            'layout/layout'           => __DIR__ . '/../view/layout/layout.phtml',
            'application/index/index' => __DIR__ . '/../view/application/index/index.phtml',
            'error/404'               => __DIR__ . '/../view/error/404.phtml',
            'error/index'             => __DIR__ . '/../view/error/index.phtml',
        ],
        'template_path_stack'      => [
            __DIR__ . '/../view',
        ],

        'strategies' => [
            'ViewJsonStrategy',
        ],
    ],

    'router' => [
        'routes' => [
            'metrics' => [
                'type'    => Literal::class,
                'options' => [
                    'route'    => '/metrics',
                    'defaults' => [
                        'controller' => PrometheusController::class,
                        'action'     => 'metrics',
                    ],
                ],
            ],
        ],
    ],
];
