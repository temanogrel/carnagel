<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Service\DeathFile;

use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\Entity\DeathFileEntity;

interface UrlServiceInterface
{
    /**
     * Create a new url entry
     *
     * @param UrlEntry        $url
     * @param DeathFileEntity $deathFile
     *
     * @return void
     */
    public function create(UrlEntry $url, DeathFileEntity $deathFile);

    /**
     * Update a url entry
     *
     * @param UrlEntry $url
     *
     * @return void
     */
    public function update(UrlEntry $url);

    /**
     * Remove a url entry
     *
     * @param UrlEntry $url
     *
     * @return void
     */
    public function remove(UrlEntry $url);
}
