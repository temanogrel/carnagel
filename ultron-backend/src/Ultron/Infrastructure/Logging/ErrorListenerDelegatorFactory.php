<?php
/**
 *
 *
 *
 */

namespace Ultron\Infrastructure\Logging;

use Interop\Container\ContainerInterface;
use Zend\Stratigility\Middleware\ErrorHandler;

class ErrorListenerDelegatorFactory
{
    /**
     * @param ContainerInterface $container
     * @param string             $name
     * @param callable           $callback
     *
     * @return \Zend\Stratigility\Middleware\ErrorHandler
     * @throws \Psr\Container\ContainerExceptionInterface
     * @throws \Psr\Container\NotFoundExceptionInterface
     */
    public function __invoke(ContainerInterface $container, $name, callable $callback)
    {
        $listener = $container->get(ErrorListener::class);

        /* @var $errorHandler ErrorHandler */
        $errorHandler = $callback();
        $errorHandler->attachListener($listener);


        return $errorHandler;
    }
}
