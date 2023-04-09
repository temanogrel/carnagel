<?php
/**
 *
 *
 *
 */

namespace Hermes\Options;

/**
 * Class UpstoreOptions
 *
 * Immutable options class for the upstore api
 */
class UpstoreOptions
{
    /**
     * @var string
     */
    protected $apiKey;

    /**
     * @var string
     */
    protected $apiUri;

    /**
     * @param string $apiKey
     * @param string $apiUri
     */
    public function __construct($apiKey, $apiUri)
    {
        $this->apiUri = $apiUri;
        $this->apiKey = $apiKey;
    }

    /**
     * @return string
     */
    public function getApiKey()
    {
        return $this->apiKey;
    }

    /**
     * @return string
     */
    public function getApiUri()
    {
        return $this->apiUri;
    }
}
