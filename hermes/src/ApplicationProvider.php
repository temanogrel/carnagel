<?php
/**
 *
 *
 *
 */

namespace Hermes;

use Doctrine\Common\Cache\ArrayCache;
use Doctrine\Common\Cache\RedisCache;
use Doctrine\DBAL\DriverManager;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\Tools\Setup;
use Hermes\Controller\ApiController;
use Hermes\Controller\IndexController;
use Hermes\Entity\UrlEntity;
use Hermes\Http\Response\FailedAuthenticationResponse;
use Hermes\Hydrator\UrlHydrator;
use Hermes\Options\UpstoreOptions;
use Hermes\Repository\UrlRepositoryInterface;
use Hermes\Service\UpstoreService;
use Hermes\Service\UrlService;
use Hermes\Service\UrlServiceInterface;
use Redis;
use Silex\Application;
use Silex\ServiceProviderInterface;
use Symfony\Component\HttpFoundation\Request;

class ApplicationProvider implements ServiceProviderInterface
{
    /**
     * Registers services on the given app.
     *
     * This method should only be used to configure services and parameters.
     * It should not get services.
     *
     * @param Application $app
     */
    public function register(Application $app)
    {
        $app['dbal'] = $app->share(function () use ($app) {
            $params = [
                'driver'   => 'pdo_mysql',
                'host'     => $app['doctrine']['host'],
                'user'     => $app['doctrine']['user'],
                'password' => $app['doctrine']['password'],
                'dbname'   => $app['doctrine']['dbname']
            ];

            return DriverManager::getConnection($params);
        });

        $app['cache'] = $app->share(function() use ($app) {
            if ($app['debug']) {
                return new ArrayCache();
            }

            $redis = new Redis();
            $redis->connect($app['redis']['host'], $app['redis']['port']);

            $cache = new RedisCache();
            $cache->setRedis($redis);

            return $cache;
        });

        $app['objectManager'] = $app->share(function () use ($app) {

            $config = Setup::createAnnotationMetadataConfiguration(
                [__DIR__ . '/Entity'],
                $app['debug'],
                $app['doctrine']['proxy_path'],
                $app['cache']
            );

            return EntityManager::create($app['dbal'], $config);
        });

        $app['options.upstore'] = $app->share(function () use ($app) {
            return new UpstoreOptions(
                $app[UpstoreOptions::class]['apiKey'],
                $app[UpstoreOptions::class]['apiUri']
            );
        });

        $app['service.url'] = $app->share(function () use ($app) {
            return new UrlService($app['objectManager'], $app['service.upstore']);
        });

        $app['service.upstore'] = $app->share(function () use ($app) {
            return new UpstoreService($app['options.upstore']);
        });

        $app['controller.index'] = $app->share(function () use ($app) {

            /* @var $repository UrlRepositoryInterface */
            $repository = $app['objectManager']->getRepository(UrlEntity::class);

            /* @var $service UrlServiceInterface */
            $service = $app['service.url'];

            return new IndexController($repository, $service);
        });

        $app['controller.api'] = $app->share(function () use ($app) {

            $hydrator = new UrlHydrator();

            /* @var $repository UrlRepositoryInterface */
            $repository = $app['objectManager']->getRepository(UrlEntity::class);

            /* @var $service UrlServiceInterface */
            $service = $app['service.url'];

            return new ApiController($hydrator, $repository, $service);
        });
    }

    /**
     * Bootstraps the application.
     *
     * This method is called after all services are registered
     * and should be used for "dynamic" configuration (whenever
     * a service must be requested).
     *
     * @param Application $app
     */
    public function boot(Application $app)
    {
        $injectHost = function (Request $request) {

            // Get the host as the query param, if it does not exist fallback to current host
            $host = $request->query->get('host', $request->getHost());

            // Store the host as an attribute so it can be accessed in the controller
            $request->attributes->set('host', $host);
            $request->attributes->set('request', $request);
        };

        $authenticate = function (Request $request) use ($app) {
            $auth = $request->headers->get('Authorization');

            if (strcmp($auth, $app['api_token']) !== 0) {
                return new FailedAuthenticationResponse();
            }
        };

        $app->get('/{key}', 'controller.index:redirect')
            ->assert('key', '[a-zA-Z0-9-]+')
            ->before($injectHost)
            ->bind('redirect');

        $app->get('/', 'controller.index:pageNotFound')
            ->bind('404');

        $app->get('/api/url', 'controller.api:create')
            ->method('POST')
            ->before($injectHost)
            ->before($authenticate);

        $resourceActions = [
            'GET'    => 'get',
            'DELETE' => 'delete',
            'PATCH'  => 'update'
        ];

        foreach ($resourceActions as $method => $action) {

            $app->get('/api/url/{key}', 'controller.api:' . $action)
                ->before($injectHost)
                ->before($authenticate)
                ->method($method);
        }
    }
}
