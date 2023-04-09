<?php
/**
 *
 *
 *
 */

namespace Hermes\Http\Response;

use Symfony\Component\HttpFoundation\JsonResponse;

class ResourceNotFoundResponse extends JsonResponse
{
    public function __construct($data = null, $status = 200, $headers = [])
    {
        parent::__construct(['message' => 'Resource not found', 'code' => 404], 404);
    }
}
