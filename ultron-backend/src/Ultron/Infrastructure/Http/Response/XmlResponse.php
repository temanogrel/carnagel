<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Http\Response;

use Zend\Diactoros\Response\TextResponse;

class XmlResponse extends TextResponse
{
    /**
     * @inheritDoc
     */
    public function __construct($text, $status = 200, array $headers = [])
    {
        $headers['content-type'] = 'application/xml';

        parent::__construct($text, $status, $headers);
    }
}
