<?php
/**
 *
 * 
 */

declare(strict_types=1);

namespace Aphrodite\Logger\Factory\Listener;

use Aphrodite\Logger\Listener\ErrorListener;
use Aphrodite\Logger\Service\LogHandlerService;
use Psr\Container\ContainerInterface;

final class ErrorListenerFactory
{
    /**
     * @param ContainerInterface $container
     * @return ErrorListener
     * @throws \Psr\Container\ContainerExceptionInterface
     */
    public function __invoke(ContainerInterface $container): ErrorListener
    {
        return new ErrorListener($container->get(LogHandlerService::class));
    }
}
