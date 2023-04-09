<?php
/**
 *
 *
 *
 */

namespace Hermes\Entity;

use DateTime;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Index;
use Doctrine\ORM\Mapping\Table;

/**
 * Class UrlEntity
 *
 * @Table(
 *  name="urls",
 *  indexes={
 *      @Index(name="short_url", columns={"short_key", "hostname"}),
 *      @Index(name="original_url", columns={"originalUrl"}),
 *      @Index(name="has_upstore_hash", columns={"upstoreDownloadHash", "isUpstore"})
 *  },
 *  options={
 *  "collate": "utf8_bin"
 * })
 *
 * @Entity(repositoryClass="Hermes\Repository\UrlRepository")
 */
class UrlEntity
{
    /**
     * @var int
     *
     * @Id()
     * @Column(type="integer")
     * @GeneratedValue(strategy="AUTO")
     */
    protected $id;

    /**
     * The hostname the url was created under
     *
     * @var string
     *
     * @Column()
     */
    protected $hostname;

    /**
     * The short key used for identification
     *
     * @var string
     *
     * @Column(name="short_key", nullable=true)
     */
    protected $key;

    /**
     * The original url
     *
     * @var string
     *
     * @Column()
     */
    protected $originalUrl;

    /**
     * Using upstores own url short script
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $upstoreDownloadHash;

    /**
     * The number of times the short url has been used
     *
     * @var int
     *
     * @Column(type="integer")
     */
    protected $transmissions = 0;

    /**
     * @var bool
     *
     * @Column(type="boolean")
     */
    protected $isUpstore = true;

    /**
     * Creation date
     *
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $createdAt;

    /**
     * Last time of modification
     *
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $updatedAt;

    /**
     * @return int
     */
    public function getId()
    {
        return $this->id;
    }

    /**
     * @param int $id
     */
    public function setId($id)
    {
        $this->id = (int)$id;
    }

    /**
     * @return string
     */
    public function getHostname()
    {
        return $this->hostname;
    }

    /**
     * @param string $hostname
     */
    public function setHostname($hostname)
    {
        $this->hostname = (string)$hostname;
    }

    /**
     * @return string
     */
    public function getKey()
    {
        return $this->key;
    }

    /**
     * @param string $key
     */
    public function setKey($key)
    {
        $this->key = (string)$key;
    }

    /**
     * @return string
     */
    public function getOriginalUrl()
    {
        return $this->originalUrl;
    }

    /**
     * @param string $originalUrl
     */
    public function setOriginalUrl($originalUrl)
    {
        $this->originalUrl = (string)$originalUrl;
    }

    /**
     * @return int
     */
    public function getTransmissions()
    {
        return $this->transmissions;
    }

    /**
     * @param int $transmissions
     */
    public function setTransmissions($transmissions)
    {
        $this->transmissions = (int)$transmissions;
    }

    /**
     * @return DateTime
     */
    public function getCreatedAt()
    {
        return $this->createdAt;
    }

    /**
     * @param DateTime $createdAt
     */
    public function setCreatedAt(DateTime $createdAt)
    {
        $this->createdAt = $createdAt;
    }

    /**
     * @return DateTime
     */
    public function getUpdatedAt()
    {
        return $this->updatedAt;
    }

    /**
     * @param DateTime $updatedAt
     */
    public function setUpdatedAt(DateTime $updatedAt)
    {
        $this->updatedAt = $updatedAt;
    }

    /**
     * @return bool
     */
    public function hasUpstoreDownloadHash()
    {
        return $this->upstoreDownloadHash !== null && $this->upstoreDownloadHash !== '';
    }

    /**
     * @return string
     */
    public function getUpstoreDownloadHash()
    {
        return $this->upstoreDownloadHash;
    }

    /**
     * @param string $hash
     */
    public function setUpstoreDownloadHash(string $hash = null)
    {
        $this->upstoreDownloadHash = $hash;
    }

    /**
     * @return boolean
     */
    public function isUpstore()
    {
        return $this->isUpstore;
    }

    /**
     * @param boolean $isUpstore
     */
    public function setIsUpstore($isUpstore)
    {
        $this->isUpstore = (bool)$isUpstore;
    }
}
