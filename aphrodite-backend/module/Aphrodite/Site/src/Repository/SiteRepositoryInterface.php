<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\Site\Repository;

use Aphrodite\Site\Entity\Site;
use Zend\Stdlib\ParametersInterface;

interface SiteRepositoryInterface
{
    /**
     * Retrieve all wordpress sites matching
     *
     * @param ParametersInterface $parameters
     *
     * @return Site[]
     */
    public function matchingQueryParameters(ParametersInterface $parameters);

    /**
     * @param int $id
     *
     * @return Site|null
     */
    public function getById($id);
}
