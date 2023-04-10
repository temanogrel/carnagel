<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Repository;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Doctrine\Common\Collections\Collection;
use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Collections\Selectable;
use Generator;
use Zend\Paginator\Paginator;

interface PerformerRepositoryInterface extends Selectable
{
    /**
     * Retrieve a performer by it's id
     *
     * @param int $id
     *
     * @return AbstractPerformerEntity|null
     */
    public function getById($id);

    /**
     * Retrieve a performer by it's service id and service
     *
     * @param string $service
     * @param string $id
     *
     * @return AbstractPerformerEntity
     */
    public function getByServiceId($service, $id);

    /**
     * Mark all performers as offline
     *
     * The return value is the number of performers affected by this.
     *
     * @param string|null $service
     *
     * @return int
     */
    public function markAllAsOffline($service = null);

    /**
     * Retrieve the number of performers matching the given state
     *
     * The following states are available as well as null, which means all performers
     * - recording
     * - pending
     *
     * @param string|null $state
     *
     * @return int
     */
    public function getPerformerCount($state = null);

    /**
     * Update the updatedAt on all performers that are online for the given service
     *
     * @param string $service
     *
     * @return int
     */
    public function updateAllOnline($service);

    /**
     * Identical to matching apart from the fact that it returns a paginator
     *
     * @param Criteria $criteria
     * @param null|string $service
     * @param null|string $indexBy
     *
     * @return Paginator
     */
    public function search(Criteria $criteria, $service = null, $indexBy = null);

    /**
     * Extend the basic criteria matching with service selection
     *
     * @param Criteria $criteria
     * @param null|string $service
     * @param null|string $indexBy
     *
     * @return Paginator
     */
    public function matching(Criteria $criteria, $service = null, $indexBy = null);

    /**
     * Get all blacklisted performers
     *
     * @return AbstractPerformerEntity[]|Collection
     */
    public function getBlacklisted();

    /**
     * @return string[]
     */
    public function getAvailableServices();

    /**
     * Get the number of performers with is_recording / is_pending_recording per service
     *
     * @return Generator
     */
    public function getPerformerStats(): Generator;
}
