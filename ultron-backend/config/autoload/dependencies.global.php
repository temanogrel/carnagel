<?php
use Cocur\Slugify\Slugify;
use Ultron\Domain\Factory\SitesFactory;
use Ultron\Domain\Service\CacheService;
use Ultron\Domain\Service\CacheServiceFactory;
use Ultron\Domain\Service\PerformerService;
use Ultron\Domain\Service\PerformerServiceFactory;
use Ultron\Domain\Service\RecordingService;
use Ultron\Domain\Service\RecordingServiceFactory;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Console\Command\BuildCacheCommand;
use Ultron\Infrastructure\Console\Command\BuildCacheCommandFactory;
use Ultron\Infrastructure\Console\Command\BuildPageCacheCommand;
use Ultron\Infrastructure\Console\Command\BuildPageCacheCommandFactory;
use Ultron\Infrastructure\Console\Command\GenerateSitemapCommand;
use Ultron\Infrastructure\Console\Command\GenerateSitemapCommandFactory;
use Ultron\Infrastructure\Console\Command\RebuildPerformerRecordingCountCommand;
use Ultron\Infrastructure\Console\Command\RebuildPerformerRecordingCountCommandFactory;
use Ultron\Infrastructure\DoctrineRedisCache;
use Ultron\Infrastructure\RedisFactory;
use Ultron\Infrastructure\Service\PageCacheService;
use Ultron\Infrastructure\Service\PageCacheServiceFactory;
use Ultron\Infrastructure\Service\SitemapService;
use Ultron\Infrastructure\Service\SitemapServiceFactory;
use Zend\Expressive\Application;
use Zend\Expressive\Container\ApplicationFactory;
use Zend\Expressive\Helper\ServerUrlHelper;
use Zend\Expressive\Helper\UrlHelper;
use Zend\Expressive\Helper\UrlHelperFactory;
use Zend\Expressive\Twig\TwigEnvironmentFactory;
use Zend\ServiceManager\Factory\InvokableFactory;

return [
    'dependencies' => [

        'invokables' => [
            ServerUrlHelper::class => ServerUrlHelper::class,
        ],

        'factories' => [
            // Application bootstrapping
            Application::class => ApplicationFactory::class,
            UrlHelper::class   => UrlHelperFactory::class,

            // Misc
            Sites::class              => SitesFactory::class,
            Slugify::class            => InvokableFactory::class,
            Redis::class              => RedisFactory::class,
            DoctrineRedisCache::class => RedisFactory::class,
            Twig_Environment::class => TwigEnvironmentFactory::class,

            // Services
            CacheService::class     => CacheServiceFactory::class,
            RecordingService::class => RecordingServiceFactory::class,
            PerformerService::class => PerformerServiceFactory::class,
            PageCacheService::class => PageCacheServiceFactory::class,
            SiteMapService::class   => SitemapServiceFactory::class,

            // Console
            RebuildPerformerRecordingCountCommand::class => RebuildPerformerRecordingCountCommandFactory::class,
            BuildPageCacheCommand::class                 => BuildPageCacheCommandFactory::class,
            BuildCacheCommand::class                     => BuildCacheCommandFactory::class,
            GenerateSitemapCommand::class                => GenerateSitemapCommandFactory::class,
        ],
    ],
];
