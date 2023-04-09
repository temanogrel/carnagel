<?php
/**
 *
 *
 *
 */

namespace Hermes\Http\Response;

use Symfony\Component\HttpFoundation\JsonResponse;

class FailedAuthenticationResponse extends JsonResponse
{
    public function __construct()
    {
        parent::__construct(['message' => 'Failed authentication'], JsonResponse::HTTP_UNAUTHORIZED);
    }
}
