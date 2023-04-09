<?php

use Hermes\ApplicationProvider;
use Silex\Provider\ServiceControllerServiceProvider;

chdir(__DIR__);

require 'vendor/autoload.php';

$app = new Silex\Application(require __DIR__ . '/config/config.php');
$app->register(new ServiceControllerServiceProvider());
$app->register(new ApplicationProvider());

return $app;
