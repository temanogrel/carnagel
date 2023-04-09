<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Service;

use Redis;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\Exception\CacheEntryMissingException;
use Ultron\Domain\SiteConfiguration;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Zend\Stdlib\ArrayUtils;

class CacheService
{
    /**
     * @var Redis
     */
    private $redis;

    /**
     * @var Sites
     */
    private $sites;

    /**
     * CacheService constructor.
     *
     * @param Redis $redis
     * @param Sites $sites
     */
    public function __construct(Redis $redis, Sites $sites)
    {
        $this->redis = $redis;
        $this->sites = $sites;
    }

    private function generateCacheKey(SiteConfiguration $site):string
    {
        return sprintf('ultron[%s]total_records', $site->getDomain());
    }

    /**
     * Build the cache for all sites
     *
     * Requires the repository as an argument because else it would cause circular dependencies
     *
     * @param RecordingRepositoryInterface $recordingRepository
     *
     * @return void
     */
    public function buildCache(RecordingRepositoryInterface $recordingRepository)
    {
        foreach ($this->sites->getSiteConfigurations() as $site) {
            if (!$site->isEnabled()) {
                continue;
            }
            
            $this->redis->set($this->generateCacheKey($site), $recordingRepository->getTotalCount($site));
        }
    }

    /**
     * @param SiteConfiguration $site
     *
     * @return int
     */
    public function getCountForSite(SiteConfiguration $site): int
    {
        $key = $this->generateCacheKey($site);

        if (!$this->redis->exists($key)) {
            throw new CacheEntryMissingException();
        }

        return (int) $this->redis->get($this->generateCacheKey($site));
    }

    public function addRecording(RecordingEntity $recording): void
    {
        foreach ($this->sites->getSiteConfigurations() as $site) {
            if ($site->isEnabled() && $recording->getPerformer()->belongsTo($site)) {
                $this->incrementCountForSite($site);
            }
        }
    }

    public function removeRecording(RecordingEntity $recording): void
    {
        foreach ($this->sites->getSiteConfigurations() as $site) {
            if ($site->isEnabled() && $recording->getPerformer()->belongsTo($site)) {
                $this->decrementCountForSite($site);
            }
        }
    }

    private function incrementCountForSite(SiteConfiguration $site): void
    {
        $this->redis->incr($this->generateCacheKey($site));
    }

    private function decrementCountForSite(SiteConfiguration $site): void
    {
        $this->redis->decr($this->generateCacheKey($site));
    }

    public function setCount(SiteConfiguration $site, int $count): void
    {
        $this->redis->set($this->generateCacheKey($site), $count);
    }
}
