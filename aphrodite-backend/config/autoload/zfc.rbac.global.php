<?php

use Aphrodite\Blocktrail\BlocktrailPermissions;
use Aphrodite\Performer\Service\PerformerService;
use Aphrodite\Recording\DeathFilePermissions;
use Aphrodite\Recording\Service\RecordingService;
use Aphrodite\Site\PostAssociationPermissions;
use Aphrodite\Site\Service\SiteService;
use Aphrodite\User\Rbac\IdentityProvider;
use ZfcRbac\Role\InMemoryRoleProvider;

return [
    'zfc_rbac' => [
        'identity_provider' => IdentityProvider::class,

        'guest_role'    => 'guest',
        'assertion_map' => [],
        'role_provider' => [
            InMemoryRoleProvider::class => [

                /**
                 * Role used by all the servers when communicating with the api
                 */
                'server' => [
                    'permissions' => [
                        BlocktrailPermissions::CREATE_ADDRESS,
                    ],

                    'children' => ['admin'],
                ],

                /**
                 * Role used by ultron when syncing
                 */
                'ultron' => [
                    'permissions' => [
                        PerformerService::PERMISSION_READ,
                        RecordingService::PERMISSION_READ,
                    ],
                ],

                /**
                 * Role of a logged in user in aphrodite UI
                 */
                'admin'  => [
                    'permissions' => [
                        'dashboard',

                        DeathFilePermissions::VIEW_DEATH_FILE,
                        DeathFilePermissions::LIST_DEATH_FILES,
                        DeathFilePermissions::UPLOAD_DEATH_FILE,
                        DeathFilePermissions::UPDATE_DEATH_FILE,
                        DeathFilePermissions::DELETE_DEATH_FILE,

                        DeathFilePermissions::ADD_ENTRY,
                        DeathFilePermissions::VIEW_ENTRY,
                        DeathFilePermissions::LIST_ENTRIES,
                        DeathFilePermissions::REMOVE_ENTRY,
                        DeathFilePermissions::UPDATE_ENTRY,

                        PostAssociationPermissions::CREATE,
                        PostAssociationPermissions::DELETE,

                        SiteService::PERMISSION_READ,
                        SiteService::PERMISSION_UPDATE,
                        SiteService::PERMISSION_CREATE,
                        SiteService::PERMISSION_DELETE,

                        PerformerService::PERMISSION_UPDATE,
                        PerformerService::PERMISSION_READ,
                        PerformerService::PERMISSION_INTERSECT,

                        RecordingService::PERMISSION_CREATE,
                        RecordingService::PERMISSION_DELETE,
                        RecordingService::PERMISSION_UPDATE,
                        RecordingService::PERMISSION_READ,
                    ],
                ],

                'guest' => [],
            ],
        ],
    ],
];
