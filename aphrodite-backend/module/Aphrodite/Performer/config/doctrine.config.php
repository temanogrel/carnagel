<?php
/**
 *
 *
 *
 */

use Doctrine\ORM\Mapping\Driver\AnnotationDriver;

return [
    'driver' => [
        'aphrodite_performer_annotation_driver' => [
            'class'     => AnnotationDriver::class,
            'paths'     => [
                'default' => __DIR__ . '/../src/Entity/',
            ]
        ],

        'orm_default' => [
            'drivers' => [
                'Aphrodite\Performer\Entity' => 'aphrodite_performer_annotation_driver'
            ]
        ]
    ]
];
