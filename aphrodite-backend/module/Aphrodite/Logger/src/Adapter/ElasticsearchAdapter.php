<?php

declare(strict_types=1);

namespace Aphrodite\Logger\Adapter;

use Elastica\Client;
use Elastica\Document;

final class ElasticsearchAdapter extends AbstractAdapter
{
    /**
     * @var Client
     */
    private $client;

    /**
     * @param Client $client
     */
    public function __construct(Client $client)
    {
        $this->client  = $client;
    }

    /**
     * {@inheritdoc}
     */
    public function write(array $data, string $type = null): bool
    {
        $index = $this->client->getIndex(sprintf('http-logs-%s', date('Y-m-d')));
        $type  = $index->getType($type);

        $document = new Document();
        $document->setData($data);

        $type->addDocument($document);

        return true;
    }
}
