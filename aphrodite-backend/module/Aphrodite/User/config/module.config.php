<?php
/**
 *
 *
 *
 */

use Aphrodite\User\Factory\Rbac\IdentityProviderFactory;use Aphrodite\User\Rbac\IdentityProvider;

return [
    'service_manager' => [
        'factories' => [
            IdentityProvider::class => IdentityProviderFactory::class
        ]
    ],

    'view_manager' => [
        'template_path_stack' => [
            'Aphrodite\User' => __DIR__ . '/../view/'
        ]
    ],

    'doctrine' => include __DIR__ . '/doctrine.config.php',
    'router' => [
        'routes' => include __DIR__ . '/route.config.php'
    ]
];
