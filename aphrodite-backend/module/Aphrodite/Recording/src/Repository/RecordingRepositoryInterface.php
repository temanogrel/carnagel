<?php
/**
 *
 *
 *  AB
 */

declare(strict_types = 1);

namespace Aphrodite\Recording\Repository;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Repository\Exception\NonUniqueRecordingResultException;
use Aphrodite\Recording\Repository\Exception\RecordingNotFoundException;
use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Collections\Selectable;
use Generator;
use Zend\Paginator\Paginator;

interface RecordingRepositoryInterface extends Selectable
{
    const MODE_HOURLY = 'hourly';
    const MODE_MINUTELY = 'minutely';
    const MODE_DAILY = 'daily';

    /**
     * Get a recording by it's id
     *
     * @param int $id
     *
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getById(int $id):RecordingEntity;

    /**
     * Get a recording by it's old ID
     *
     * @param int $id
     *
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getByOldId(int $id):RecordingEntity;

    /**
     * Retrieve a video by it's url
     *
     * @param string $url
     *
     * @throws NonUniqueRecordingResultException
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getByVideoUrl(string $url):RecordingEntity;

    /**
     * Retrieve a recording by it's gallery url
     *
     * @param string $url
     *
     * @throws NonUniqueRecordingResultException
     * @throws RecordingNotFoundException
     *
     * @return RecordingEntity
     */
    public function getByGalleryUrl(string $url):RecordingEntity;

    /**
     * Get the recordings for a given performer
     *
     * @param AbstractPerformerEntity $performer
     *
     * @return Paginator
     */
    public function getForPerformer(AbstractPerformerEntity $performer):Paginator;

    /**
     * @param string $mode
     * @param int    $limit
     *
     * @return array
     */
    public function getRecordingTimeLine(string $mode = self::MODE_HOURLY, $limit = 10):array;

    /**
     * Retrieves the number of recordings in the various available states
     *
     * @return Generator
     */
   public function getRecordingsPerState(): Generator;

    /**
     * Retrieve the number of recordings matching a given state
     *
     * @param Criteria $criteria
     *
     * @return int
     */
    public function getRecordingCount(Criteria $criteria):int;
}
