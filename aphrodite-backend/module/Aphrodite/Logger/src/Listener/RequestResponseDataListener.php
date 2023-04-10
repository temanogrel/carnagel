<?php
/**
 *
 * 
 */

declare(strict_types=1);

namespace Aphrodite\Logger\Listener;

use Aphrodite\Logger\Service\LogHandlerServiceInterface;
use Zend\EventManager\EventManagerInterface;
use Zend\EventManager\ListenerAggregateInterface;
use Zend\EventManager\ListenerAggregateTrait;
use Zend\Http\Request as HttpRequest;
use Zend\Http\Response as HttpResponse;
use Zend\Mvc\MvcEvent;

final class RequestResponseDataListener implements ListenerAggregateInterface
{
    use ListenerAggregateTrait;

    /**
     * @var LogHandlerServiceInterface
     */
    private $service;

    /**
     * LogRequestResponseDataListener constructor.
     * @param LogHandlerServiceInterface $service
     */
    public function __construct(LogHandlerServiceInterface $service)
    {
        $this->service = $service;
    }

    /**
     * {@inheritdoc}
     */
    public function attach(EventManagerInterface $events)
    {
        $this->listeners[] = $events->attach(MvcEvent::EVENT_FINISH, [$this, 'handleRequestResponseData'], 1000);
    }

    /**
     * Call log handler service with request/response and matched route
     *
     * @param MvcEvent $event
     */
    public function handleRequestResponseData(MvcEvent $event)
    {
        /* @var HttpRequest $request */
        $request = $event->getRequest();
        /* @var HttpResponse $response */
        $response = $event->getResponse();

        $routeMatch = $event->getRouteMatch();

        // We currently only support http request/responses
        // todo: add support for console requests too
        if ($request instanceof HttpRequest && $response instanceof HttpResponse) {
            $duration = defined('START_TIME') ? (microtime(true) - START_TIME) * 1000 : null;

            $this->service->handleRequestResponse($request, $response, $routeMatch, $duration);
        }
    }
}
