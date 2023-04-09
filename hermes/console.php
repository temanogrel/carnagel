<?php
/**
 *
 *
 *
 */

use Doctrine\ORM\EntityManager;
use Doctrine\DBAL\Tools\Console\ConsoleRunner as DBALConsoleRunner;
use Doctrine\ORM\Tools\Console\ConsoleRunner as ORMConsoleRunner;
use Hermes\Command\Import\CurImport;
use Hermes\Command\Import\PipImport;
use Hermes\Command\Sync\UpstoreSync;
use Knp\Provider\ConsoleServiceProvider;

/* @var $app \Silex\Application */
$app = require __DIR__ . '/bootstrap.php';
$app->register(new ConsoleServiceProvider(), [
    'console.name'              => 'HermesCli',
    'console.version'           => '1.0.2',
    'console.project_directory' => __DIR__,
]);

/* @var $entityManager EntityManager */
$entityManager = $app['objectManager'];

$helperSet = new \Symfony\Component\Console\Helper\HelperSet([
    'db' => new \Doctrine\DBAL\Tools\Console\Helper\ConnectionHelper($entityManager->getConnection()),
    'em' => new \Doctrine\ORM\Tools\Console\Helper\EntityManagerHelper($entityManager),
]);

/* @var $console \Knp\Console\Application */
$console = $app['console'];

// Register the entity manager
$console->setHelperSet($helperSet);

// Register the doctrine commands
ORMConsoleRunner::addCommands($console);
DBALConsoleRunner::addCommands($console);

// Add out custom commands
$console->add(new UpstoreSync());

// Run the application
$console->run();
