<?php
/**
 *
 *
 *
 */

use Zend\Stdlib\ArrayUtils;

$files = ArrayUtils::iteratorToArray($this->files);

return [
    'data' => array_map(function ($file) {
        return $this->renderResource('death-file/resource', ['file' => $file]);
    }, $files)
];
