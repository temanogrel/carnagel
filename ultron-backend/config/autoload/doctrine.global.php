<?php
use ContainerInteropDoctrine\EntityManagerFactory;
use Doctrine\Common\Cache\RedisCache;
use Doctrine\DBAL\Driver\PDOMySql\Driver;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\Mapping\Driver\AnnotationDriver;
use Ultron\Infrastructure\DoctrineRedisCache;
use Ultron\Infrastructure\RepositoryFactory;
use Ultron\Infrastructure\RepositoryFactoryFactory;

return [
    'dependencies' => [
        'factories' => [
            EntityManager::class     => EntityManagerFactory::class,
            RepositoryFactory::class => RepositoryFactoryFactory::class,
        ],
    ],

    'doctrine' => [
        'configuration' => [
            'orm_default' => [
                'driver'                        => 'orm_default',
                'auto_generate_proxy_classes'   => false,
                'proxy_dir'                     => 'data/cache/DoctrineEntityProxy',
                'proxy_namespace'               => 'DoctrineEntityProxy',
                'entity_namespaces'             => [],
                'datetime_functions'            => [],
                'string_functions'              => [],
                'numeric_functions'             => [],
                'filters'                       => [],
                'named_queries'                 => [],
                'named_native_queries'          => [],
                'custom_hydration_modes'        => [],
                'naming_strategy'               => null,
                'default_repository_class_name' => null,
                'repository_factory'            => RepositoryFactory::class,
                'class_metadata_factory_name'   => null,
                'entity_listener_resolver'      => null,
                'second_level_cache'            => [
                    'enabled'                    => false,
                    'default_lifetime'           => 3600,
                    'default_lock_lifetime'      => 60,
                    'file_lock_region_directory' => '',
                    'regtions'                   => [],
                ],
                'sql_logger'                    => null,
            ],
        ],

        'connection' => [
            'orm_default' => [
                'driver_class'             => Driver::class,
                'wrapper_class'            => null,
                'pdo'                      => null,
                'configuration'            => 'orm_default',
                'event_manager'            => 'orm_default',
                'doctrine_mapping_types'   => [],
                'doctrine_commented_types' => [],
            ],
        ],

        'entity_manager' => [
            'orm_default' => [
                'connection'    => 'orm_default', // Actually defaults to the entity manager config key, not hard-coded
                'configuration' => 'orm_default', // Actually defaults to the entity manager config key, not hard-coded
            ],
        ],

        'event_manager' => [
            'orm_default' => [
                'subscribers' => [],
            ],
        ],

        'driver' => [
            'orm_default' => [
                'class' => AnnotationDriver::class,
                'cache' => 'array',
                'paths' => [
                    'src/Ultron/Domain/Entity',
                ],
            ],
        ],

        'cache' => [
            'redis' => [
                'class'     => RedisCache::class,
                'instance'  => DoctrineRedisCache::class,
                'namespace' => 'ultron',
            ],

        ],
        'types' => [],
    ],
];
