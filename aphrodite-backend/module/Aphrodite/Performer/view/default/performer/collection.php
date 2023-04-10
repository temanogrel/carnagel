<?php
/**
 *
 *
 *
 */

use Zend\Paginator\Paginator;
use Zend\Stdlib\ArrayUtils;

$performers = array_map(function ($performer) {
    return $this->renderResource('performer/resource', ['performer' => $performer]);
}, ArrayUtils::iteratorToArray($this->performers));

if ($this->performers instanceof Paginator) {
    return [
        'meta' => $this->renderPaginator($this->performers),
        'data' => $performers
    ];
}

return [
    'meta' => [
        'total'  => count($performers),
        'limit'  => count($performers),
        'offset' => 0
    ],

    'data' => $performers
];
