<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service;

use Generator;
use Ultron\Domain\SiteConfiguration;

interface SitemapServiceInterface
{
    /**
     * @param SiteConfiguration $site
     * @return Generator
     */
    public function getSitemapUrls(SiteConfiguration $site): Generator;

    /**
     * Create all the performer sitemaps
     *
     * @return void
     */
    public function createPerformerSitemaps();

    /**
     * Create all the recording sitemaps
     *
     * @return void
     */
    public function createRecordingSitemaps();
}
