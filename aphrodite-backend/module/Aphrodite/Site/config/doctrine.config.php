<?php
/**
 *
 *
 *
 */

use Doctrine\ORM\Mapping\Driver\AnnotationDriver;

return [
    'driver' => [
        'aphrodite_site_annotation_driver' => [
            'class'     => AnnotationDriver::class,
            'paths'     => [
                'default' => __DIR__ . '/../src/Entity/',
            ]
        ],

        'orm_default' => [
            'drivers' => [
                'Aphrodite\Site\Entity' => 'aphrodite_site_annotation_driver'
            ]
        ]
    ]
];
