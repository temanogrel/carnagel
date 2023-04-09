<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service;

use Doctrine\Common\Collections\Criteria;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\SiteConfiguration;
use Zend\Paginator\Paginator;

interface PageCacheServiceInterface
{
    /**
     * Creates the page cache for a given site configuration
     *
     * @param SiteConfiguration $site
     *
     * @return void
     */
    public function create(SiteConfiguration $site);

    /**
     * Updates the page cache of all sites after a new recording is added
     *
     * @param RecordingEntity $recording
     */
    public function addRecording(RecordingEntity $recording);

    /**
     * Gets the page cache entry for a given site and page
     *
     * @param Criteria $criteria
     * @param SiteConfiguration $site
     * @param int $page
     *
     * @return Paginator
     */
    public function get(Criteria $criteria, SiteConfiguration $site, int $page): Paginator;
}
