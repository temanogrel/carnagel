<?php
/**
 *
 *
 *
 */

namespace Ultron\Domain\Entity;

use DateTime;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Index;
use Doctrine\ORM\Mapping\OneToMany;
use Doctrine\ORM\Mapping\Table;
use Doctrine\ORM\Mapping\UniqueConstraint;
use DomainException;
use Ultron\Domain\SiteConfiguration;

/**
 * @Table(
 *  name="performers",
 *  indexes={
 *      @Index(name="uid", columns={"uid"}),
 *      @Index(name="slug", columns={"slug"}),
 *      @Index(name="section", columns={"section"}),
 *      @Index(name="service", columns={"service"}),
 *      @Index(name="stageName", columns={"stageName"}),
 *      @index(name="createdAt", columns={"createdAt"}),
 *      @index(name="updatedAt", columns={"updatedAt"})
 *  },
 *
 *  uniqueConstraints={
 *      @UniqueConstraint(name="uid_per_service", columns={"uid", "service"})
 *  }
 * )
 *
 * @Entity(repositoryClass="Ultron\Infrastructure\Repository\PerformerRepository")
 */
class PerformerEntity
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
     * Internal id
     *
     * @var integer
     *
     * @Column(type="integer")
     */
    protected $uid;

    /**
     * @var string
     *
     * @Column(unique=true)
     */
    protected $slug;

    /**
     * @var string
     *
     * @Column()
     */
    protected $stageName;

    /**
     * @var array
     *
     * @Column(type="json_array")
     */
    protected $aliases = [];

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $section;

    /**
     * @var string
     *
     * @Column()
     */
    protected $service;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $recordingCount = 0;

    /**
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $updatedAt;

    /**
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $createdAt;

    /**
     * @var Collection
     *
     * @OneToMany(targetEntity="Ultron\Domain\Entity\RecordingEntity", mappedBy="performer")
     */
    protected $recordings;

    public function __construct()
    {
        $this->recordings = new ArrayCollection();
    }

    /**
     * Check if the performer is supposed to be displayed on the given site
     *
     * @param SiteConfiguration $configuration
     *
     * @return bool
     */
    public function belongsTo(SiteConfiguration $configuration):bool
    {
        if ($configuration->getService() !== null && $configuration->getService() !== $this->getService()) {
            return false;
        }

        if ($configuration->getSection() !== null && $configuration->getSection() !== $this->section) {
            return false;
        }

        return true;
    }

    /**
     * Retrieve the service name
     *
     * @param bool $fullName
     *
     * @throws DomainException
     *
     * @return string
     */
    public function getService(bool $fullName = false):string
    {
        if (!$fullName) {
            return $this->service;
        }

        switch ($this->service) {
            case 'mfc':
                return 'MyFreeCams';

            case 'cbc':
                return 'Chaturbate';

            case 'cam':
                return 'cam9';

            default:
                throw new DomainException('Unknown service');
        }
    }

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
    public function setId(int $id)
    {
        $this->id = $id;
    }

    /**
     * @return int
     */
    public function getUid():int
    {
        return $this->uid;
    }

    /**
     * @param integer $uid
     */
    public function setUid(int $uid)
    {
        $this->uid = $uid;
    }

    /**
     * @return string
     */
    public function getSlug()
    {
        return $this->slug;
    }

    /**
     * @param string $slug
     */
    public function setSlug(string $slug)
    {
        $this->slug = $slug;
    }

    /**
     * @return string
     */
    public function getStageName()
    {
        return $this->stageName;
    }

    /**
     * @param string $stageName
     */
    public function setStageName(string $stageName)
    {
        $this->stageName = $stageName;
    }

    /**
     * @return array
     */
    public function getAliases():array
    {
        return $this->aliases;
    }

    /**
     * @param array $aliases
     */
    public function setAliases(array $aliases)
    {
        $this->aliases = $aliases;
    }

    /**
     * @return string
     */
    public function getSection()
    {
        return $this->section;
    }

    /**
     * @param string $section
     */
    public function setSection(string $section)
    {
        $this->section = $section;
    }


    /**
     * @param string $service
     */
    public function setService(string $service)
    {
        $this->service = $service;
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
     * @return int
     */
    public function getRecordingCount():int
    {
        return $this->recordingCount;
    }

    /**
     * @param int $recordingCount
     */
    public function setRecordingCount(int $recordingCount)
    {
        $this->recordingCount = $recordingCount;
    }

    /**
     * Increment the number of recordings
     *
     * @return void
     */
    public function incrementRecordingCount()
    {
        $this->recordingCount++;
    }

    /**
     * @return Collection
     */
    public function getRecordings():Collection
    {
        return $this->recordings;
    }

    /**
     * @param Collection $recordings
     */
    public function setRecordings(Collection $recordings)
    {
        $this->recordings = $recordings;
    }
}
