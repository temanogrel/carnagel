<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Repository;

use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Collections\Selectable;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\SiteConfiguration;
use Ultron\Infrastructure\Service\ValueObject\PageInformation;
use Zend\Paginator\Paginator;

interface RecordingRepositoryInterface extends Selectable
{
    /**
     * Get a recording by it's id
     *
     * @param int $id
     *
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getById($id): RecordingEntity;

    /**
     * Get a recording by it's internal uid
     *
     * @param string $uid
     *
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getByUid($uid): RecordingEntity;

    /**
     * Get a recording by it's slug
     *
     * @param string $slug
     *
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getBySlug($slug): RecordingEntity;

    /**
     * Get a recording by it's gallery url
     *
     * @param string $url
     *
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getByGalleryUrl($url): RecordingEntity;

    /**
     * Search for performers by the given stage name and criteria
     *
     * @param string                 $stageName
     * @param Criteria|null          $criteria
     * @param SiteConfiguration|null $site
     *
     * @return RecordingEntity[]|Paginator
     */
    public function searchByPerformer(string $stageName, Criteria $criteria = null, SiteConfiguration $site = null): Paginator;

    /**
     * Get a paginated result of recording matching the criteria
     *
     * @param Criteria          $criteria
     * @param SiteConfiguration $site
     *
     * @return RecordingEntity[]|Paginator
     */
    public function getPaginatedResult(Criteria $criteria, SiteConfiguration $site): Paginator;

    /**
     * Increment the number of views
     *
     * @param RecordingEntity $recording
     *
     * @return void
     */
    public function incrementViewCount(RecordingEntity $recording);

    /**
     * Get the total count for the given site
     *
     * @param SiteConfiguration $site
     *
     * @return int
     */
    public function getTotalCount(SiteConfiguration $site = null): int;

    /**
     * Get the max id of a page for the page cache
     *
     * @param SiteConfiguration $site
     * @param int $limit
     * @param int $minId
     * @param int $maxId
     *
     * @return PageInformation
     */
    public function getPageInformation(
        SiteConfiguration $site,
        int $limit,
        int $minId,
        int $maxId = null
    ): PageInformation;

    /**
     * Gets all recordings between entries
     *
     * @param Criteria $criteria
     * @param SiteConfiguration $site
     * @param int $min
     * @param int $max
     *
     * @return array
     */
    public function getBetweenIds(Criteria $criteria, SiteConfiguration $site, int $min, int $max): array;

    /**
     * Gets all recordings from criteria and site
     *
     * @param Criteria $criteria
     * @param SiteConfiguration $site
     * @return array
     */
    public function getBySiteAndCriteria(Criteria $criteria, SiteConfiguration $site): array;

    /**
     * Calculates the recording count of each performer
     *
     * @return array
     */
    public function getRecordingCountOfPerformers(): array;
}
