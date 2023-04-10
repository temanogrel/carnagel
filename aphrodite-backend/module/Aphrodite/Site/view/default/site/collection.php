<?php
/**
 *
 *
 *
 */

use Zend\Stdlib\ArrayUtils;

$sites = ArrayUtils::iteratorToArray($this->sites);

return [
    'data' => array_map(function ($site) {
        return $this->renderResource('site/resource', ['site' => $site]);
    }, $sites)
];
