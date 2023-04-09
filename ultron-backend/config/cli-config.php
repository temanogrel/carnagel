<?php
chdir(dirname(__DIR__));

require 'vendor/autoload.php';

/* @var $container \Interop\Container\ContainerInterface */
$container = require __DIR__ . '/container.php';

return new \Symfony\Component\Console\Helper\HelperSet([
    'em' => new \Doctrine\ORM\Tools\Console\Helper\EntityManagerHelper(
        $container->get(\Doctrine\ORM\EntityManager::class)
    ),
]);
