<?php
/**
 *
 *
 *
 */

namespace Ultron\Infrastructure\Logging;

use const DATE_RFC3339;
use Elastica\Client;
use Elastica\Document;
use function get_class;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Throwable;

class ErrorListener
{
    /**
     * @var Client
     */
    private $client;

    /**
     * ErrorListener constructor.
     *
     * @param Client $client
     */
    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function __invoke($error, ServerRequestInterface $request, ResponseInterface $response)
    {
        $index = $this->client->getIndex(sprintf('ultron-%s', date('Y-m-d')));
        $type = $index->getType('errors');

        if ($error instanceof Throwable) {
            $data = [
                '@timestamp' => date(DATE_RFC3339),
                'message'    => $error->getMessage(),
                'class'      => get_class($error),
                'stacktrace' => $error->getTraceAsString(),
            ];

            $this->recurseException($data, $error->getPrevious());
        } else {
            $data = [
                '@timestamp' => date(DATE_RFC3339),
                'message'    => 'None exception error occurred',
                'class'      => get_class($error),
                'payload'    => $error,
            ];
        }

        $document = new Document();
        $document->setData($data);

        $type->addDocument($document);
    }

    private function recurseException(&$data, Throwable $exception = null)
    {
        if ($exception === null) {
            return;
        }

        $data['previous'] = [
            'class'      => get_class($exception),
            'message'    => $exception->getMessage(),
            'stacktrace' => $exception->getTraceAsString(),
        ];

        $this->recurseException($data['previous'], $exception->getPrevious());
    }
}
