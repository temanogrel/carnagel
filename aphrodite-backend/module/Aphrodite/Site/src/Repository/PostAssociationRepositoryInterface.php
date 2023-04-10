<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Aphrodite\Site\Repository;

use Aphrodite\Site\Entity\PostAssociation;
use Aphrodite\Site\Repository\Exception\PostAssociationNotFoundException;

interface PostAssociationRepositoryInterface
{
    /**
     * Retrieve a post association by it's id
     * 
     * @param int $id
     *
     * @throws PostAssociationNotFoundException
     *
     * @return PostAssociation
     */
    public function getById(int $id):PostAssociation;
}
