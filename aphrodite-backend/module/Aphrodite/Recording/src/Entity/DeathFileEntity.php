<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Entity;

use DateTime;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;

/**
 * Entity that represents
 *
 * @Entity(repositoryClass="Aphrodite\Recording\Repository\DeathFileRepository")
 */
class DeathFileEntity
{
    /**
     * @var integer
     *
     * @Id()
     * @GeneratedValue()
     * @Column(type="integer")
     */
    protected $id;

    /**
     * Path in the virtual filesystem
     *
     * @var string
     *
     * @Column(type="string", unique=true)
     */
    protected $location;

    /**
     * @var integer
     *
     * @Column(type="integer", nullable=true)
     */
    protected $entries;

    /**
     * @var integer
     *
     * @Column(type="integer")
     */
    protected $ignored = 0;

    /**
     * @var integer
     *
     * @Column(type="integer")
     */
    protected $pending = 0;

    /**
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $createdAt;

    /**
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
        $this->id = $id;
    }

    /**
     * @return string
     */
    public function getLocation()
    {
        return $this->location;
    }

    /**
     * @param string $location
     */
    public function setLocation($location)
    {
        $this->location = $location;
    }

    /**
     * @return int
     */
    public function getEntries()
    {
        return $this->entries;
    }

    /**
     * @param int $entries
     */
    public function setEntries($entries)
    {
        $this->entries = $entries;
    }

    /**
     * @return int
     */
    public function getIgnored()
    {
        return $this->ignored;
    }

    /**
     * @param int $ignored
     */
    public function setIgnored($ignored)
    {
        $this->ignored = (int) $ignored;
    }

    /**
     * @return int
     */
    public function getPending()
    {
        return $this->pending;
    }

    /**
     * @param int $pending
     */
    public function setPending($pending)
    {
        $this->pending = (int) $pending;
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
    public function setCreatedAt($createdAt)
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
    public function setUpdatedAt($updatedAt)
    {
        $this->updatedAt = $updatedAt;
    }
}
