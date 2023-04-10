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
use Doctrine\ORM\EntityRepository;

final class PostAssociationRepository extends EntityRepository implements PostAssociationRepositoryInterface
{
    /**
     * @inheritDoc
     */
    public function getById(int $id):PostAssociation
    {
        $association = $this->findOneBy(['id' => $id]);
        if (!$association) {
            throw new PostAssociationNotFoundException();
        }

        return $association;
    }
}
