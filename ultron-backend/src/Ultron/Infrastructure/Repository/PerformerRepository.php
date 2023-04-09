<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Repository;

use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityRepository;
use Doctrine\ORM\Tools\Pagination\Paginator as DoctrinePaginator;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Exception\PerformerNotFoundException;
use Ultron\Infrastructure\Paginator\Adapter\DoctrineAdapter;
use Zend\Paginator\Adapter\Callback;
use Zend\Paginator\Paginator;

final class PerformerRepository extends EntityRepository implements PerformerRepositoryInterface
{
    /**
     * {@inheritdoc}
     */
    public function getById(int $id):PerformerEntity
    {
        $performer = $this->findOneBy(['id' => $id]);
        if (!$performer) {
            throw new PerformerNotFoundException();
        }

        return $performer;
    }

    /**
     * {@inheritdoc}
     */
    public function getBySlug(string $slug):PerformerEntity
    {
        $performer = $this->findOneBy(['slug' => $slug]);
        if (!$performer) {
            throw new PerformerNotFoundException();
        }

        return $performer;
    }

    /**
     * {@inheritdoc}
     */
    public function getByUid(int $uid):PerformerEntity
    {
        $performer = $this->findOneBy(['uid' => $uid]);
        if (!$performer) {
            throw new PerformerNotFoundException();
        }

        return $performer;
    }

    /**
     * {@inheritdoc}
     */
    public function getPaginatedResult(Criteria $criteria = null)
    {
        $builder = $this->createQueryBuilder('performer');
        $builder->addCriteria($criteria);

        return new Paginator(new DoctrineAdapter(new DoctrinePaginator($builder)));
    }

    /**
     * {@inheritdoc}
     */
    public function getByCriteria(Criteria $criteria): array
    {
        $builder = $this->createQueryBuilder('performer');
        $builder->addCriteria($criteria);

        return $builder->getQuery()->getResult();
    }
}
