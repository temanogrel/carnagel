<?php
/**
 *
 *
 *
 */

namespace Hermes\Controller;

use Symfony\Component\HttpFoundation\Response;

class ErrorController
{
    /**
     * Render the 404 page
     *
     * @return Response
     */
    public function pageNotFound()
    {
        return Response::create('Page not found', Response::HTTP_NOT_FOUND);
    }
}
