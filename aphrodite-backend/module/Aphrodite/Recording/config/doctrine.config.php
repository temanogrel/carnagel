<?php
/**
 *
 *
 *
 */

use Doctrine\ORM\Mapping\Driver\AnnotationDriver;

return [
    'driver' => [
        'aphrodite_recording_annotation_driver' => [
            'class'     => AnnotationDriver::class,
            'paths'     => [
                'default' => __DIR__ . '/../src/Entity/',
            ]
        ],

        'orm_default' => [
            'drivers' => [
                'Aphrodite\Recording\Entity' => 'aphrodite_recording_annotation_driver'
            ]
        ]
    ]
];
