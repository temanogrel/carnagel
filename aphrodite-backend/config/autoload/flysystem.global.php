<?php

return [
    'bsb_flysystem' => [
        'adapters' => [
            'local' => [
                'type'    => 'local',
                'options' => [
                    'root' => './data/storage'
                ]
            ]
        ],

        'filesystems' => [
            'default' => [
                'adapter'   => 'local',
                'cache'     => false,
                'eventable' => false,
            ]
        ]
    ]
];
