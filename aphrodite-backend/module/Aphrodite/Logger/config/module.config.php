<?php

use Elastica\Client;
use Aphrodite\Logger\Adapter\ElasticsearchAdapter;
use Aphrodite\Logger\Factory\Adapter\ElasticsearchAdapterFactory;
use Aphrodite\Logger\Factory\Client\ElasticaClientFactory;
use Aphrodite\Logger\Factory\Listener\ErrorListenerFactory;
use Aphrodite\Logger\Factory\Listener\RequestResponseDataListenerFactory;
use Aphrodite\Logger\Factory\Options\ElasticsearchOptionsFactory;
use Aphrodite\Logger\Factory\Options\LogHandlerOptionsFactory;
use Aphrodite\Logger\Factory\Service\LogHandlerServiceFactory;
use Aphrodite\Logger\Listener\ErrorListener;
use Aphrodite\Logger\Listener\RequestResponseDataListener;
use Aphrodite\Logger\Options\ElasticsearchOptions;
use Aphrodite\Logger\Options\LogHandlerOptions;
use Aphrodite\Logger\Service\LogHandlerService;

return [
    'service_manager' => [
        'factories' => [
            RequestResponseDataListener::class => RequestResponseDataListenerFactory::class,
            ErrorListener::class               => ErrorListenerFactory::class,

            ElasticsearchAdapter::class => ElasticsearchAdapterFactory::class,

            LogHandlerService::class => LogHandlerServiceFactory::class,

            Client::class => ElasticaClientFactory::class,

            LogHandlerOptions::class    => LogHandlerOptionsFactory::class,
            ElasticsearchOptions::class => ElasticsearchOptionsFactory::class,
        ],
    ],
];
