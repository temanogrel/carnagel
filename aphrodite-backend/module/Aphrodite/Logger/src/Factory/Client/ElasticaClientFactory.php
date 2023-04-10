<?php
/**
 *
 *
 * 
 */
declare(strict_types=1);

namespace Aphrodite\Logger\Factory\Client;

use Elastica\Client;
use Aphrodite\Logger\Options\ElasticsearchOptions;
use Interop\Container\ContainerInterface;

final class ElasticaClientFactory
{
    /**
     * @param ContainerInterface $container
     * @return Client
     * @throws \Psr\Container\ContainerExceptionInterface
     */
    public function __invoke(ContainerInterface $container): Client
    {
        /* @var ElasticsearchOptions $options */
        $options = $container->get(ElasticsearchOptions::class);

        return new Client([
            'host'     => $options->getHost(),
            'port'     => $options->getPort(),
            'username' => $options->getUsername(),
            'password' => $options->getPassword(),
            'timeout'  => 1,
        ]);
    }
}
