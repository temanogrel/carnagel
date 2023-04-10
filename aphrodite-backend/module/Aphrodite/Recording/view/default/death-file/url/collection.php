<?php
/**
 *
 *
 *
 */

use Zend\Stdlib\ArrayUtils;

$urls = ArrayUtils::iteratorToArray($this->urls);

return [
    'meta' => $this->renderPaginator($this->urls),
    'data' => array_map(function ($url) {
        return $this->renderResource('death-file/url/resource', ['url' => $url]);
    }, $urls)
];
