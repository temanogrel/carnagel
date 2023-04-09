<?php
/**
 *
 *
 *
 */

namespace Hermes\Repository;

use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityRepository;
use Doctrine\ORM\NonUniqueResultException;
use Doctrine\ORM\NoResultException;
use Doctrine\ORM\Tools\Pagination\Paginator;
use Hermes\Entity\UrlEntity;

class UrlRepository extends EntityRepository implements UrlRepositoryInterface
{
    /**
     * @inheritDoc
     */
    public function getById($id)
    {
        return $this->findOneBy(['id' => $id]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByKeyAndHostname($key, $hostname)
    {
        return $this->findOneBy(['key' => $key, 'hostname' => $hostname]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByOriginalUrl($url)
    {
        return $this->findOneBy(['originalUrl' => $url]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByOriginalUrlWithWildcard($url)
    {
        $builder = $this->createQueryBuilder('u');
        $builder->andWhere('u.originalUrl LIKE :url');
        $builder->setParameter('url', $url . '%');

        try {
            return $builder->getQuery()->getOneOrNullResult();
        } catch (NonUniqueResultException $e) {
            return null;
        }
    }

    public function getByUpstoreCode(string $code)
    {
        $builder = $this->createQueryBuilder('u');
        $builder->andWhere('u.originalUrl LIKE :first OR u.originalUrl LIKE :second');
        $builder->setParameter('first', sprintf('http://upstor.re/%s%%', $code));
        $builder->setParameter('second', sprintf('http://upstore.net/%s%%', $code));

        try {
            return $builder->getQuery()->getOneOrNullResult();
        } catch (NonUniqueResultException $e) {
            return null;
        }
    }

    /**
     * {@inheritdoc}
     */
    public function getCount(Criteria $criteria = null)
    {
        $criteria = $criteria ?: Criteria::create();

        $builder = $this->createQueryBuilder('u');
        $builder
            ->select('COUNT(u)')
            ->addCriteria($criteria);

        try {
            return $builder->getQuery()->getSingleScalarResult();
        } catch (NoResultException $e) {
            return 0;
        } catch (NonUniqueResultException $e) {
            return 0;
        }
    }

    /**
     * {@inheritdoc}
     */
    public function getPaginator(Criteria $criteria)
    {
        $builder = $this->createQueryBuilder('u');
        $builder->addCriteria($criteria);

        return new Paginator($builder->getQuery());
    }

}
