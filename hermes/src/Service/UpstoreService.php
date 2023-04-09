<?php
/**
 *
 *
 *
 */

namespace Hermes\Service;

use GuzzleHttp\Client;
use GuzzleHttp\ClientInterface;
use Hermes\Entity\UrlEntity;
use Hermes\Options\UpstoreOptions;
use Hermes\Service\Exception\FailedToMatchUrlException;

class UpstoreService implements UpstoreServiceInterface
{
    /**
     * @var ClientInterface
     */
    private $client;

    /**
     * @var UpstoreOptions
     */
    private $options;

    /**
     * @param UpstoreOptions       $options
     * @param ClientInterface|null $client
     */
    public function __construct(UpstoreOptions $options, ClientInterface $client = null)
    {
        $this->client  = $client ?: new Client();
        $this->options = $options;
    }

    private function extractHash($url)
    {
        $matches = [];

        if (preg_match('/upsto\.re\/(?P<hash>[a-zA-Z0-9]+)/', $url, $matches)) {
            return $matches['hash'];
        }

        if (preg_match('/upstore\.net\/(?P<hash>[a-zA-Z0-9]+)$/', $url, $matches)) {
            return $matches['hash'];
        }

        throw new FailedToMatchUrlException();
    }

    /**
     * {@inheritdoc}
     */
    public function getDownlinkHash(UrlEntity $url)
    {
        $hash = $this->extractHash($url->getOriginalUrl());

        $response = $this->client->get($this->options->getApiUri(), [
            'query' => [
                'hash' => $hash,
                'key'  => $this->options->getApiKey()
            ]
        ]);

        if ($response->getStatusCode() !== 200) {
            throw new \RuntimeException($response->getBody()->getContents());
        }

        $json = json_decode($response->getBody()->getContents(), true);

        return $json['links'][$hash];
    }

    /**
     * {@inheritdoc}
     */
    public function syncUrlCollection($collection)
    {
        if (count($collection) > 1000) {
            throw new \OutOfBoundsException('Collection may not be larger than 1000 elements');
        }

        $urls   = [];
        $hashes = [];

        foreach ($collection as $url) {

            try {
                $hash = $this->extractHash($url->getOriginalUrl());

                // Need to do it this way, because possible to have multiple hash for the same
                $urls[]   = $url;
                $hashes[] = $hash;

            } catch (FailedToMatchUrlException $e) {
                $url->setIsUpstore(false);
            }
        }

        $response = $this->client->post($this->options->getApiUri(), [
            'query' => [
                'key' => $this->options->getApiKey()
            ],
            'form_params' => [
                'hash' => implode(',', $hashes)
            ]
        ]);

        $json = json_decode($response->getBody()->getContents(), true);

        /* @var $url UrlEntity */
        foreach ($urls as $url) {
            $url->setUpstoreDownloadHash($json['links'][$this->extractHash($url->getOriginalUrl())]);
        }
    }
}
