<?php
/**
 *
 *
 *
 */

namespace Hermes\Service;

use Hermes\Entity\UrlEntity;

interface UrlServiceInterface
{
    /**
     * Create a url entity from the string
     *
     * @param string $url
     * @param string $hostname
     *
     * @return UrlEntity
     */
    public function create($url, $hostname);

    /**
     * Update a url entity
     *
     * @param UrlEntity $url
     *
     * @return void
     */
    public function update(UrlEntity $url);

    /**
     * Delete a url entity
     *
     * @param UrlEntity $url
     *
     * @return void
     */
    public function delete(UrlEntity $url);

    /**
     * Increment the number of transmissions
     *
     * @param UrlEntity $url
     *
     * @return void
     */
    public function incrementTransmissions(UrlEntity $url);
}
