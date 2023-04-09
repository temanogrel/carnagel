<?php
/**
 *
 *
 *
 */

namespace Hermes\Http\Response;

use Symfony\Component\HttpFoundation\JsonResponse;

class FailedDataValidationResponse extends JsonResponse
{
    public function __construct(array $errors)
    {
        $data = [
            'message' => 'Invalid data',
            'errors'  => $errors
        ];

        parent::__construct($data, JsonResponse::HTTP_BAD_REQUEST);
    }
}
