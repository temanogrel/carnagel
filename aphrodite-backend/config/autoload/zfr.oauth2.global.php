<?php

use Zend\Authentication\AuthenticationService;
use ZfrOAuth2Module\Server\Factory\AuthenticationServiceFactory;

return [
    /**
     * Uncomment the factory if you are using a stateless REST API and want to authenticate your users
     * using the access tokens
     */
    'service_manager' => [
        'factories' => [
            AuthenticationService::class => AuthenticationServiceFactory::class
        ]
    ],

    'zfr_oauth2_server' => [

        /**
         * Doctrine object manager key
         */
        'object_manager' => 'Aphrodite\ObjectManager',

        /**
         * Various tokens TTL
         */
        'authorization_code_ttl' => 120,
        'access_token_ttl'       => 3600,
        'refresh_token_ttl'      => 86400,

        /**
         * Registered grants for this server
         */
        'grants' => [],

        /**
         * Grant plugin manager
         *
         * The configuration follows a standard service manager configuration
         */
        'grant_manager' => [
            'factories' => []
        ],
    ]
];
