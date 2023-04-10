<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Repository;

use Doctrine\ORM\EntityRepository;
use Zend\Stdlib\ParametersInterface;

class SiteRepository extends EntityRepository implements SiteRepositoryInterface
{
    public function getById($id)
    {
        return $this->findOneBy(['id' => $id]);
    }

    public function matchingQueryParameters(ParametersInterface $parameters)
    {
        $builder = $this->createQueryBuilder('site');

        if ($parameters->get('enabled') !== null) {
            $builder->andWhere($builder->expr()->eq('site.enabled', (int) $parameters->get('enabled')));
        }

        return $builder->getQuery()->getResult();
    }
}
