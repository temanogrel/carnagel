<?php
/**
 *
 *
 *
 */

namespace Hermes\Repository;

use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Collections\Selectable;
use Doctrine\Common\Persistence\ObjectRepository;
use Hermes\Entity\UrlEntity;

interface UrlRepositoryInterface extends ObjectRepository, Selectable
{
    /**
     * Retrieve a url by it's internal id
     *
     * @param integer $id
     *
     * @return UrlEntity|null
     */
    public function getById($id);

    /**
     * Retrieve a url entity by key and hostname
     *
     * @param string $key
     * @param string $hostname
     *
     * @return UrlEntity|null
     */
    public function getByKeyAndHostname($key, $hostname);

    /**
     * Retrieve a recording by it's original url
     *
     * @param string $url
     *
     * @return UrlEntity|null
     */
    public function getByOriginalUrl($url);

    /**
     * Retrieve a recording by it's original url with a appended wildcard
     *
     * @param $url|null
     *
     * @return UrlEntity
     */
    public function getByOriginalUrlWithWildcard($url);

    /**
     * @param string $code
     *
     * @return UrlEntity|null
     */
    public function getByUpstoreCode(string $code);

    /**
     * Retrieve a count for all the url meeting the given criteria
     *
     * @param Criteria|null $criteria
     *
     * @return int
     */
    public function getCount(Criteria $criteria = null);
}
