<?php
/**
 *
 *
 *
 */
declare(strict_types=1);

namespace Aphrodite\Logger\Factory\Adapter;

use Elastica\Client;
use Aphrodite\Logger\Adapter\ElasticsearchAdapter;
use Aphrodite\Logger\Options\ElasticsearchOptions;
use Interop\Container\ContainerInterface;

final class ElasticsearchAdapterFactory
{
    /**
     * @param ContainerInterface $container
     * @return ElasticsearchAdapter
     */
    public function __invoke(ContainerInterface $container): ElasticsearchAdapter
    {
        /* @var Client $elasticaClient */
        $elasticaClient = $container->get(Client::class);

        return new ElasticsearchAdapter($elasticaClient);
    }
}
