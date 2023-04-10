<?php
/**
 *
 *
 *
 */

use Zend\Stdlib\ArrayUtils;

$recordings = ArrayUtils::iteratorToArray($this->recordings);

return [
    'meta' => $this->renderPaginator($this->recordings),
    'data' => array_map(function ($recording) {
        return $this->renderResource('recording/resource', ['recording' => $recording]);
    }, $recordings)
];
