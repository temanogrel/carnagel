<?php

use Aphrodite\Stdlib\DateFormat;
use Aphrodite\User\Entity\UserEntity;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\Mapping\UnderscoreNamingStrategy;
use ZfrOAuth2\Server\Entity\TokenOwnerInterface;

return [
    'service_manager' => [
        'aliases'    => [
            'Aphrodite\ObjectManager'          => EntityManager::class,
            'Roave\NonceUtility\ObjectManager' => EntityManager::class,
        ],
        'invokables' => [
            UnderscoreNamingStrategy::class => UnderscoreNamingStrategy::class
        ]
    ],
    'doctrine'        => [
        'connection'      => [
            'orm_default' => [
                'driverClass' => 'Doctrine\DBAL\Driver\PDOMySql\Driver',
                'params'      => [
                    'charset' => 'utf8',
                ],
            ]
        ],
        'configuration'   => [
            'orm_default' => [
                'string_functions' => [
                    'DATE_FORMAT' => DateFormat::class
                ],

                'generate_proxies' => false,
                'proxy_dir'        => 'data/cache/doctrine-proxies',
                'proxy_namespace'  => 'Aphrodite\DoctrineProxies',

                'naming_strategy'  => UnderscoreNamingStrategy::class,
            ]
        ],

        'migrations_configuration' => [
            'orm_default' => [
                'name'      => 'Database migrations',
                'namespace' => 'Aphrodite\DbMigrations',
                'directory' => 'data/migrations',
                'table'     => 'doctrine_migrations',
            ],
        ],

        'cache' => [
            'redis' => [
                'namespace' => 'Aphrodite',
                'instance'  => Redis::class
            ],
        ],

        'entity_resolver' => [
            'orm_default' => [
                'resolvers' => [
                    TokenOwnerInterface::class => UserEntity::Class
                ]
            ]
        ]
    ],
];
