<?php
/**
 *
 *
 */

declare(strict_types = 1);

chdir(__DIR__);

require __DIR__ . '/vendor/autoload.php';

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Symfony\Component\Console\Application;
use Zend\Expressive\Router\Route;
use Zend\Expressive\Router\RouterInterface;
use Doctrine\ORM\Tools\Console\ConsoleRunner as ORMConsoleRunner;

/* @var ContainerInterface $container */
$container = require __DIR__ . '/config/container.php';
$application = new Application('Application console');

/* @var $entityManager EntityManager */
$entityManager = $container->get(EntityManager::class);
$entityManager
    ->getConnection()
    ->exec('SET wait_timeout=28800; SET interactive_timeout=28800');

$helperSet = new \Symfony\Component\Console\Helper\HelperSet(array(
    'db' => new \Doctrine\DBAL\Tools\Console\Helper\ConnectionHelper($entityManager->getConnection()),
    'em' => new \Doctrine\ORM\Tools\Console\Helper\EntityManagerHelper($entityManager)
));

$application->setHelperSet($helperSet);

$commands = $container->get('config')['console']['commands'];

ORMConsoleRunner::addCommands($application);

foreach ($commands as $command) {
    $application->add($container->get($command));
}

/* @var RouterInterface $router */
$router = $container->get(RouterInterface::class);
$routes = $config['routes'];

foreach ($routes as $route) {
    $router->addRoute(new Route($route['path'], $route['middleware'], $route['allowed_methods'], $route['name']));
}

$application->run();
