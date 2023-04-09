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
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Exception\PerformerNotFoundException;
use Zend\Paginator\Paginator;

interface PerformerRepositoryInterface extends Selectable
{
    /**
     * Retrieve a performer by it's id
     *
     * @param int $id
     *
     * @throws PerformerNotFoundException
     *
     * @return PerformerEntity
     */
    public function getById(int $id):PerformerEntity;

    /**
     * Retrieve a recording
     *
     * @param int $uid
     *
     * @return PerformerEntity
     */
    public function getByUid(int $uid):PerformerEntity;

    /**
     * Get a performer by it's slug
     *
     * @param string $slug
     *
     * @throws PerformerNotFoundException
     *
     * @return PerformerEntity
     */
    public function getBySlug(string $slug):PerformerEntity;

    /**
     * Get a paginated result of performers matching the given criteria
     *
     * @param Criteria|null $criteria
     *
     * @return PerformerEntity[]|Paginator
     */
    public function getPaginatedResult(Criteria $criteria = null);

    /**
     * Gets all performers between min and max id
     *
     * @param Criteria $criteria
     *
     * @return array
     */
    public function getByCriteria(Criteria $criteria): array;
}
