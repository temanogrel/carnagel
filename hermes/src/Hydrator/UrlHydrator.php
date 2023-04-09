<?php
/**
 *
 *
 *
 */

namespace Hermes\Hydrator;

use DateTime;
use Hermes\Entity\UrlEntity;

class UrlHydrator
{
    public function extract(UrlEntity $url)
    {
        return [
            'id'                  => $url->getId(),
            'key'                 => $url->getKey(),
            'hostname'            => $url->getHostname(),
            'shortUrl'            => sprintf('http://%s/%s', $url->getHostname(), $url->getKey()),
            'originalUrl'         => $url->getOriginalUrl(),
            'isUpstore'           => $url->isUpstore(),
            'upstoreDownloadHash' => $url->getUpstoreDownloadHash(),

            // Dates
            'createdAt' => $url->getCreatedAt()->format(DateTime::ISO8601),
            'updatedAt' => $url->getUpdatedAt()->format(DateTime::ISO8601)
        ];
    }
}
