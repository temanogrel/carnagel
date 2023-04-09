<?php
use Ultron\Infrastructure\Logging\ErrorListener;
use Ultron\Infrastructure\Logging\ErrorListenerDelegatorFactory;
use Ultron\Infrastructure\Logging\ErrorListenerFactory;
use Ultron\Infrastructure\Middleware\ApiAuthenticationMiddleware;
use Ultron\Infrastructure\Middleware\ApiAuthenticationMiddlewareFactory;
use Ultron\Infrastructure\Middleware\SiteSelectionMiddleware;
use Ultron\Infrastructure\Middleware\SiteSelectionMiddlewareFactory;
use Zend\Expressive\Container\ApplicationFactory;
use Zend\Expressive\Container\ErrorHandlerFactory;
use Zend\Expressive\Helper;
use Zend\Expressive\Helper\BodyParams\BodyParamsMiddleware;
use Zend\Expressive\Helper\ServerUrlMiddleware;
use Zend\Expressive\Helper\ServerUrlMiddlewareFactory;
use Zend\Expressive\Helper\UrlHelperMiddleware;
use Zend\Expressive\Helper\UrlHelperMiddlewareFactory;
use Zend\Stratigility\Middleware\ErrorHandler;

return [
    'dependencies' => [
        'factories' => [
            ServerUrlMiddleware::class => ServerUrlMiddlewareFactory::class,
            UrlHelperMiddleware::class => UrlHelperMiddlewareFactory::class,

            SiteSelectionMiddleware::class     => SiteSelectionMiddlewareFactory::class,
            ApiAuthenticationMiddleware::class => ApiAuthenticationMiddlewareFactory::class,

            ErrorHandler::class => ErrorHandlerFactory::class,
            ErrorListener::class => ErrorListenerFactory::class,
        ],

        'delegators' => [
            ErrorHandler::class => [
                ErrorListenerDelegatorFactory::class,
            ]
        ]
    ],

    'middleware_pipeline' => [
        'always' => [
            'middleware' => [
                ServerUrlMiddleware::class,
                SiteSelectionMiddleware::class,
            ],

            'priority'   => 10000,
        ],

        'routing' => [
            'middleware' => [
                ErrorHandler::class,

                ApplicationFactory::ROUTING_MIDDLEWARE,
                UrlHelperMiddleware::class,
                ApiAuthenticationMiddleware::class,
                BodyParamsMiddleware::class => BodyParamsMiddleware::class,
                ApplicationFactory::DISPATCH_MIDDLEWARE,
            ],
            'priority'   => 1,
        ],

        'error' => [

            'middleware' => [
                // Add error middleware here.
            ],

            'error'    => true,
            'priority' => -10000,
        ],
    ],
];
