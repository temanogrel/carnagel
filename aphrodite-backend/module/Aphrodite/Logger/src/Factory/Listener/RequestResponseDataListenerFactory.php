<?php
/**
 *
 * 
 */

declare(strict_types=1);

namespace Aphrodite\Logger\Factory\Listener;

use Aphrodite\Logger\Listener\RequestResponseDataListener;
use Aphrodite\Logger\Service\LogHandlerService;
use Psr\Container\ContainerInterface;

final class RequestResponseDataListenerFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return RequestResponseDataListener
     *
     * @throws \Psr\Container\ContainerExceptionInterface
     */
    public function __invoke(ContainerInterface $container): RequestResponseDataListener
    {
        return new RequestResponseDataListener($container->get(LogHandlerService::class));
    }
}
