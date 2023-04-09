<?php
/**
 *
 *
 *
 */

if (extension_loaded('newrelic')) {
    newrelic_set_appname('hermes');
}

use Silex\Application;

/* @var $app Application */
$app = require '../bootstrap.php';
$app->run();
