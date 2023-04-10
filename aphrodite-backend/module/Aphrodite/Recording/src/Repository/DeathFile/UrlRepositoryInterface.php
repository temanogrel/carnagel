<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Repository\DeathFile;

use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Collections\Selectable;
use Zend\Paginator\Paginator;

interface UrlRepositoryInterface extends Selectable
{
    /**
     * Search for all urls matching the given criteria
     *
     * Works just like {@see Selectable::matching()} but this returns a paginator
     *
     * @param Criteria $criteria
     *
     * @return Paginator
     */
    public function paginatedSearch(Criteria $criteria);

    /**
     * Retrieve a entry by it's url
     *
     * @param string $url
     *
     * @return UrlEntry|null
     */
    public function getByUrl($url);

    /**
     * Retrieve a entry by it's id
     *
     * @param integer $id
     *
     * @return UrlEntry|null
     */
    public function getById($id);
}
