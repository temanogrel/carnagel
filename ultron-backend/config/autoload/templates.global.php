<?php

use Ultron\Infrastructure\View\Extension\UltronTwigExtension;
use Ultron\Infrastructure\View\Extension\UltronTwigExtensionFactory;
use Zend\Expressive\Container\TemplatedErrorHandlerFactory;
use Zend\Expressive\Template\TemplateRendererInterface;
use Zend\Expressive\Twig\TwigRendererFactory;

return [
    'dependencies' => [
        'factories' => [
            TemplateRendererInterface::class => TwigRendererFactory::class,
            UltronTwigExtension::class       => UltronTwigExtensionFactory::class,
        ],
    ],

    'templates' => [
        'extension' => 'twig',
        'paths'     => [
            'app'    => ['templates/app'],
            'layout' => ['templates/layout'],
            'error'  => ['templates/error'],
        ],
    ],

    'twig' => [
        'cache_dir'      => 'data/cache/twig',
        'assets_url'     => '/',
        'assets_version' => null,
        'extensions'     => [
            UltronTwigExtension::class,
        ],
    ],
];
