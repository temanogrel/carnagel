<?php
/**
 *
 *
 *
 */

namespace Hermes\Service;

use Doctrine\Common\Collections\Collection;
use Hermes\Entity\UrlEntity;

interface UpstoreServiceInterface
{
    /**
     * Attempt to get the downlink hash for the url
     *
     * @param UrlEntity $url
     *
     * @return string
     */
    public function getDownlinkHash(UrlEntity $url);

    /**
     * Sync all the downlink hash with the collection
     *
     * Uses the batch upstore api to update all the urls with a downlink hash
     *
     * @param UrlEntity[] $collection
     *
     * @return void
     */
    public function syncUrlCollection($collection);
}
