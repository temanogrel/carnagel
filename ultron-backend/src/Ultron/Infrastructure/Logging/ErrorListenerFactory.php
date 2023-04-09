<?php
/**
 *
 *
 *
 */

namespace Ultron\Infrastructure\Logging;

use Elastica\Client;
use Psr\Container\ContainerInterface;

class ErrorListenerFactory
{
    public function __invoke(ContainerInterface $container): ErrorListener
    {
        $client = new Client(['servers' => $container->get('config')['ultron']['elasticsearch']]);

        return new ErrorListener($client);
    }
}
