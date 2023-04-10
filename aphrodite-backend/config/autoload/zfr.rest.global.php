<?php

use ZfcRbac\Exception\UnauthorizedException;
use ZfrOAuth2\Server\Exception\InvalidAccessTokenException;
use ZfrRest\Http\Exception\Client\UnauthorizedException as HttpUnauthorizedException;

return [
    'zfr_rest' => [
        'exception_map' => [
            UnauthorizedException::class       => HttpUnauthorizedException::class,
            InvalidAccessTokenException::class => HttpUnauthorizedException::class,
        ],
    ]
];
