<?php
/**
 *
 *
 *
 */

namespace Ultron\Infrastructure\Paginator\Adapter;

use Doctrine\ORM\Tools\Pagination\Paginator;
use Zend\Paginator\Adapter\AdapterInterface;

class DoctrineAdapter implements AdapterInterface
{
    /**
     * @var Paginator
     */
    private $paginator;

    /**
     * @var bool
     */
    private $cache;

    /**
     * @param Paginator $paginator
     * @param bool      $cache
     */
    public function __construct(Paginator $paginator, $cache = false)
    {
        $this->cache     = $cache;
        $this->paginator = $paginator;
    }

    /**
     * {@inheritdoc}
     */
    public function getItems($offset, $itemCountPerPage)
    {
        $query = $this->paginator->getQuery();
        $query->setMaxResults($itemCountPerPage);
        $query->setFirstResult($offset);

        return $query->getResult();
    }

    /**
     * {@inheritdoc}
     */
    public function count()
    {
        return $this->paginator->count();
    }
}
