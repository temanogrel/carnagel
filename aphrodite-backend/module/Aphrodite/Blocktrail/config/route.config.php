<?php
/**
 *
 *
 */

use Aphrodite\Blocktrail\Controller\RpcController;

return [
    'blocktrail-rpc-create-address' => [
        'type'    => 'literal',
        'options' => [
            'route'    => '/rpc/blocktrail.create-address',
            'defaults' => [
                'controller' => RpcController::class,
                'action'     => 'create-address',
            ],
        ],
    ],

    'blocktrail-rpc-send-bitcoin' => [
        'type'    => 'literal',
        'options' => [
            'route'    => '/rpc/blocktrail.send-bitcoin',
            'defaults' => [
                'controller' => RpcController::class,
                'action'     => 'send-bitcoin',
            ],
        ],
    ]
];
