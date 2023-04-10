<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Repository\DeathFile;

use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityRepository;
use Zend\Paginator\Adapter\Callback;
use Zend\Paginator\Paginator;

class UrlRepository extends EntityRepository implements UrlRepositoryInterface
{
    /**
     * @inheritDoc
     */
    public function getByUrl($url)
    {
        return $this->findOneBy(['url' => $url]);
    }

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
    public function paginatedSearch(Criteria $criteria)
    {
        $builder = $this->createQueryBuilder('url');
        $builder->addCriteria($criteria);

        $countBuilder = clone $builder;
        $countBuilder
            ->select('COUNT(url.url)')
            ->resetDQLPart('orderBy')
            ->setMaxResults(null)
            ->setFirstResult(null);

        $countQuery = $countBuilder->getQuery();
        $countQuery->useResultCache(true, 600);

        $count = function () use ($countQuery) {
            return $countQuery->getSingleScalarResult();
        };

        $data = function () use ($builder) {
            return $builder->getQuery()->getResult();
        };

        return new Paginator(new Callback($data, $count));
    }

}
