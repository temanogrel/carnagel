<?php
/**
 *
 *
 * 
 */

namespace Aphrodite\Logger\Service;

use Throwable;
use Zend\Http\Request;
use Zend\Http\Response;
use Zend\Mvc\Router\RouteMatch;

interface LogHandlerServiceInterface
{
    /**
     * Handle and log exceptions
     *
     * @param Throwable $exception
     * @param array $context
     */
    public function handleException(Throwable $exception, array $context = []);

    /**
     * Handle and log request/response details
     *
     * @param Request $request
     * @param Response $response
     * @param RouteMatch|null $routeMatch
     * @param float|null $duration
     * @param array $context
     * @return
     */
    public function handleRequestResponse(
        Request $request,
        Response $response,
        RouteMatch $routeMatch = null,
        float $duration = null,
        array $context = []
    );

    /**
     * Returns a list of active adapters
     *
     * @return array
     */
    public function getAdapters(): array;
}
