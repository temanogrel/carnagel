<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Repository;

use Doctrine\ORM\EntityRepository;

class DeathFileRepository extends EntityRepository implements DeathFileRepositoryInterface
{
    /**
     * {@inheritdoc}
     */
    public function getById($id)
    {
        return $this->findOneBy(['id' => $id]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByUrl($url)
    {
        return $this->findOneBy(['url' => $url]);
    }
}
