<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Entity;

use DateTime;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\DiscriminatorColumn;
use Doctrine\ORM\Mapping\DiscriminatorMap;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Index;
use Doctrine\ORM\Mapping\InheritanceType;
use Doctrine\ORM\Mapping\OneToMany;
use Doctrine\ORM\Mapping\Table;
use Doctrine\ORM\Mapping\UniqueConstraint;
use DomainException;

/**
 * Class AbstractPerformerEntity
 *
 * @Table(
 *  name="performers",
 *  indexes={
 *      @Index(name="serviceId", columns={"service_id", "service", "id"}),
 *      @Index(name="section", columns={"section"}),
 *      @Index(name="service", columns={"service"}),
 *      @Index(name="online", columns={"online", "service", "id"}),
 *      @Index(name="blacklisted", columns={"blacklisted", "service"}),
 *      @Index(name="recording", columns={"is_recording", "service", "id"}),
 *      @Index(name="pending_recording", columns={"is_pending_recording", "service", "id"}),
 *      @Index(name="stageName", columns={"stage_name"}),
 *      @index(name="createdAt", columns={"created_at"}),
 *      @index(name="updatedAt", columns={"updated_at"})
 *  },
 *
 *  uniqueConstraints={
 *      @UniqueConstraint(name="service_and_id_unq", columns={"service_id", "service"})
 *  }
 * )
 *
 * @Entity(repositoryClass="Aphrodite\Performer\Repository\PerformerRepository")
 * @InheritanceType("SINGLE_TABLE")
 * @DiscriminatorColumn(name="service", type="string")
 * @DiscriminatorMap({
 *  "mfc" = "MyFreeCamsPerformer",
 *  "cbc" = "ChaturbatePerformer",
 *  "cam4" = "Cam4Performer"
 * })
 */
abstract class AbstractPerformerEntity
{
    /**
     * @var integer
     *
     * @Id()
     * @Column(type="integer")
     * @GeneratedValue()
     */
    protected $id;

    /**
     * Internal id
     *
     * @var string
     *
     * @Column()
     */
    protected $serviceId;

    /**
     * @var string
     *
     * @Column()
     */
    protected $stageName;

    /**
     * @var boolean
     *
     * @Column(type="boolean")
     */
    protected $blacklisted = false;

    /**
     * @var bool
     *
     * @Column(type="boolean")
     */
    protected $online = false;

    /**
     * @var bool
     *
     * @Column(type="boolean")
     */
    protected $isRecording = false;

    /**
     * @var bool
     *
     * @Column(type="boolean")
     */
    protected $isPendingRecording = false;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $recordingCount = 0;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $currentViewers = 0;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $peakViewerCount = 0;

    /**
     * @var string[]
     *
     * @Column(type="simple_array", nullable=true)
     */
    protected $aliases = [];

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
     * @OneToMany(targetEntity="Aphrodite\Recording\Entity\RecordingEntity", mappedBy="performer", fetch="LAZY")
     */
    protected $recordings;

    public function __construct()
    {
        $this->recordings = new ArrayCollection();
    }

    /**
     * Should return a string matching the serviceToEntityClassName
     *
     * @return string
     */
    abstract public function getService();

    /**
     * Convert the service name to the entity class name
     *
     * @param string $service
     *
     * @return string
     */
    public static function serviceToEntityClassName($service)
    {
        switch (strtolower($service))
        {
            case 'mfc':
                return MyFreeCamsPerformer::class;

            case 'cbc':
                return ChaturbatePerformer::class;

            case 'cam':
            case 'cam4':
                return Cam4Performer::class;

            default:
                throw new DomainException('Unknown service');
        }
    }

    /**
     * @param string $alias
     *
     * @return bool
     */
    public function hasAlias($alias)
    {
        return in_array($alias, $this->aliases);
    }

    /**
     * @param string $alias
     */
    public function addAlias($alias)
    {
        if (!$this->hasAlias($alias)) {
            $this->aliases[] = $alias;
        }
    }

    /**
     * @return string[]
     */
    public function getAliases()
    {
        return $this->aliases;
    }

    /**
     * @param string[] $aliases
     */
    public function setAliases(array $aliases)
    {
        $this->aliases = $aliases;
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
    public function setId($id)
    {
        $this->id = $id;
    }

    /**
     * @return string
     */
    public function getServiceId()
    {
        return $this->serviceId;
    }

    /**
     * @param string $serviceId
     */
    public function setServiceId($serviceId)
    {
        $this->serviceId = $serviceId;
    }

    /**
     * @return boolean
     */
    public function isBlacklisted()
    {
        return $this->blacklisted;
    }

    /**
     * @param boolean $blacklisted
     */
    public function setBlacklisted($blacklisted)
    {
        $this->blacklisted = (bool) $blacklisted;
    }

    /**
     * @return int
     */
    public function getCurrentViewers()
    {
        return $this->currentViewers;
    }

    /**
     * @param int $currentViewers
     */
    public function setCurrentViewers($currentViewers)
    {
        $this->currentViewers = $currentViewers;
    }

    /**
     * @return int
     */
    public function getPeakViewerCount()
    {
        return $this->peakViewerCount;
    }

    /**
     * @param int $peakViewerCount
     */
    public function setPeakViewerCount($peakViewerCount)
    {
        $this->peakViewerCount = $peakViewerCount;
    }

    /**
     * @return boolean
     */
    public function isOnline()
    {
        return $this->online;
    }

    /**
     * @param boolean $online
     */
    public function setOnline($online)
    {
        $this->online = (bool) $online;

        if (!$online) {
            $this->setIsRecording(false);
            $this->setIsPendingRecording(false);
        }
    }

    /**
     * @return boolean
     */
    public function isRecording()
    {
        return $this->isRecording;
    }

    /**
     * @param boolean $isRecording
     */
    public function setIsRecording($isRecording)
    {
        $this->isRecording = (bool) $isRecording;

        // Both cannot be true.
        if ($isRecording) {
            $this->setIsPendingRecording(false);
        }
    }

    /**
     * @return boolean
     */
    public function isPendingRecording()
    {
        return $this->isPendingRecording;
    }

    /**
     * @param boolean $isPendingRecording
     */
    public function setIsPendingRecording($isPendingRecording)
    {
        $this->isPendingRecording = (bool) $isPendingRecording;
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
    public function setStageName($stageName)
    {
        $this->stageName = $stageName;
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
    public function getRecordingCount()
    {
        return $this->recordingCount;
    }

    /**
     * @param int $recordingCount
     */
    public function setRecordingCount($recordingCount)
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
    public function getRecordings()
    {
        return $this->recordings;
    }

    /**
     * @param Collection $recordings
     */
    public function setRecordings($recordings)
    {
        $this->recordings = $recordings;
    }
}
