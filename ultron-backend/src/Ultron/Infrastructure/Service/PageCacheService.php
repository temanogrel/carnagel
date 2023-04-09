<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service;

use Doctrine\Common\Collections\Criteria;
use Redis;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\CacheService;
use Ultron\Domain\SiteConfiguration;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Ultron\Infrastructure\Service\Exception\PageCacheEntryCreationException;
use Ultron\Infrastructure\Service\Exception\PageCacheEntryNotFoundException;
use Ultron\Infrastructure\Service\Exception\PageCountEntryNotFoundException;
use Zend\Paginator\Adapter\Callback;
use Zend\Paginator\Paginator;
use Zend\Stdlib\ArrayUtils;

final class PageCacheService implements PageCacheServiceInterface
{
    /**
     * @var Redis
     */
    private $redis;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var CacheService
     */
    private $cacheService;

    /**
     * @var Sites
     */
    private $sites;

    /**
     * PaginationCacheService constructor.
     * @param Redis $redis
     * @param RecordingRepositoryInterface $recordingRepository
     * @param CacheService $cacheService
     * @param Sites $sites
     */
    public function __construct(
        Redis $redis,
        RecordingRepositoryInterface $recordingRepository,
        CacheService $cacheService,
        Sites $sites
    ) {
        $this->recordingRepository = $recordingRepository;
        $this->redis               = $redis;
        $this->cacheService        = $cacheService;
        $this->sites               = $sites;
    }

    public function create(SiteConfiguration $site)
    {
        $rowCount = $this->cacheService->getCountForSite($site);

        $pageCountKey = $this->getPageCountRedisKey($site->getDomain());

        $processed = 0;
        $minId     = 1;

        for ($page = 0; ; $page++) {
            $offset = 5000;

            do {
                $processedTemp = $processed;

                $pageInfo = $this
                    ->recordingRepository
                    ->getPageInformation($site, $site->getPageSize(), $minId, $minId + $offset);

                $processedTemp += $pageInfo->getCount();

                $offset *= 2;
            } while ($pageInfo->getCount() < $site->getPageSize() && $processedTemp < $rowCount);

            $processed += $pageInfo->getCount();

            // Don't store the last page if it is not complete
            if ($pageInfo->getCount() === $site->getPageSize()) {
                $this->redis->set($pageCountKey, $page + 1);
            } else {
                break;
            }

            $key   = $this->getPageRedisKey($site->getDomain(), $page + 1);
            $value = $this->encodeRedisValue($minId, $pageInfo->getMaxId());

            if (!$this->redis->set($key, $value)) {
                throw new PageCacheEntryCreationException();
            }

            $minId = $pageInfo->getMaxId() + 1;
        }
    }

    public function addRecording(RecordingEntity $recording)
    {
        foreach ($this->sites->getSiteConfigurations() as $site) {
            if (!$site->isEnabled() || !$recording->getPerformer()->belongsTo($site)) {
                continue;
            }

            $pageCountKey = $this->getPageCountRedisKey($site->getDomain());

            if (!$this->redis->exists($pageCountKey)) {
                throw new PageCountEntryNotFoundException();
            }

            $pageCount = (int) $this->redis->get($pageCountKey);

            $value    = $this->getFromRedis($this->getPageRedisKey($site->getDomain(), $pageCount));
            $pageInfo = $this
                ->recordingRepository
                ->getPageInformation($site, $site->getPageSize(), $value['max'] + 1, null);

            // Do we have a complete page to put in page cache ?
            if ($pageInfo->getCount() === $site->getPageSize()) {
                $key   = $this->getPageRedisKey($site->getDomain(), $pageCount + 1);
                $value = $this->encodeRedisValue($value['max'] + 1, $pageInfo->getMaxId());

                if (!$this->redis->set($key, $value)) {
                    throw new PageCacheEntryCreationException();
                }

                $this->redis->set($pageCountKey, $pageCount + 1);
            }
        }
    }

    public function get(Criteria $criteria, SiteConfiguration $site, int $page): Paginator
    {
        $pageCountKey = $this->getPageCountRedisKey($site->getDomain());
        $page         = (int) $this->redis->get($pageCountKey) - ($page - 1);

        try {
            return $this->loadFromRedis($criteria, $site, $page);
        } catch (PageCacheEntryNotFoundException $e) {
            return $this->recordingRepository->getPaginatedResult($criteria, $site);
        }
    }

    private function getFromRedis(string $key)
    {
        if (!$this->redis->exists($key)) {
            throw new PageCacheEntryNotFoundException();
        }

        return json_decode((string) $this->redis->get($key), true);
    }

    private function getPageCountRedisKey(string $domain): string
    {
        return $domain . ':page-count';
    }

    private function getPageRedisKey(string $domain, int $page): string
    {
        return $domain . ':page-' . $page;
    }

    private function encodeRedisValue(int $min, int $max)
    {
        return json_encode([
            'min' => $min,
            'max' => $max
        ]);
    }

    /**
     * @param Criteria $criteria
     * @param SiteConfiguration $site
     * @param int $page
     *
     * @return Paginator
     */
    private function loadFromRedis(Criteria $criteria, SiteConfiguration $site, int $page): Paginator
    {
        $value = $this->getFromRedis($this->getPageRedisKey($site->getDomain(), $page));

        $itemsCallback = function () use ($criteria, $site, $value) {
            return $this->recordingRepository->getBetweenIds($criteria, $site, $value['min'], $value['max']);
        };

        $countCallback = function () use ($site) {
            return $this->cacheService->getCountForSite($site);
        };

        return new Paginator(new Callback($itemsCallback, $countCallback));
    }
}
