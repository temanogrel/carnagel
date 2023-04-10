<?php
/**
 *
 *
 *
 */

use Doctrine\ORM\Mapping\Driver\AnnotationDriver;

return [
    'driver' => [
        'aphrodite_user_annotation_driver' => [
            'class'     => AnnotationDriver::class,
            'paths'     => [
                'default' => __DIR__ . '/../src/Entity/',
            ]
        ],

        'orm_default' => [
            'drivers' => [
                'Aphrodite\User\Entity' => 'aphrodite_user_annotation_driver'
            ]
        ]
    ]
];
