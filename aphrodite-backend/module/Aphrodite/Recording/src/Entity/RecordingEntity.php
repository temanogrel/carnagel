<?php
/**
 *
 *
 *  AB
 */

declare(strict_types=1);

namespace Aphrodite\Recording\Entity;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Recording\Entity\ValueObject\Images;
use Aphrodite\Site\Entity\Site;
use DateTime;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\ORM\Mapping as ORM;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\DiscriminatorColumn;
use Doctrine\ORM\Mapping\DiscriminatorMap;
use Doctrine\ORM\Mapping\Embedded;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Index;
use Doctrine\ORM\Mapping\InheritanceType;
use Doctrine\ORM\Mapping\OneToMany;
use Doctrine\ORM\Mapping\Table;

/**
 * Class AbstractRecordingEntity
 *
 * @Table(name="recordings", indexes={
 *  @Index(name="state", columns={"state"}),
 *  @Index(name="old_id", columns={"old_id"}),
 *  @Index(name="state_within_period", columns={"state", "created_at"}),
 *  @Index(name="service", columns={"service"}),
 *  @Index(name="stageName", columns={"stage_name"}),
 *  @Index(name="videoUrl", columns={"video_url"}),
 *  @Index(name="createdAt", columns={"created_at"}),
 *  @Index(name="updatedAt", columns={"updated_at"}),
 *  @Index(name="lastCheckedAt", columns={"last_checked_at"}),
 *  @Index(name="videoMp4Uuid", columns={"video_mp4uuid"}),
 *  @Index(name="videoHlsUuid", columns={"video_hls_uuid"})
 * })
 *
 * @Entity(repositoryClass="Aphrodite\Recording\Repository\RecordingRepository")
 */
class RecordingEntity
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
     * @var integer
     *
     * @Column(type="integer", nullable=true)
     */
    protected $oldId;

    /**
     * @var string
     *
     * @Column()
     */
    protected $state;

    /**
     * The stage name used when recording the video.
     *
     * @var string
     *
     * @Column()
     */
    protected $stageName;

    /**
     * The section in which the recording was listed
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $section;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $service;

    /**
     * @var integer
     *
     * @Column(type="integer")
     */
    protected $duration = 0;

    /**
     * @var integer
     *
     * @Column(type="bigint")
     */
    protected $size264 = 0;

    /**
     * @var int
     *
     * @Column(type="bigint", nullable=true)
     */
    protected $size265 = 0;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $description;

    /**
     * This will contain the UUID for both the mp4 encoded file and the raw flv file since they are technically
     * both encoded in h264.
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $videoMp4Uuid;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $videoHlsUuid;

    /**
     * @var string
     *
     * @Column(nullable=true, type="text")
     */
    protected $videoManifest;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $wordpressCollageUuid;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $infinityCollageUuid;

    /**
     * Sprites are the large collages that infinity uses to create the video scrub
     *
     * @var string[]
     *
     * @Column(type="json_array", nullable=true)
     */
    protected $sprites = [];

    /**
     * Images are the images listed bellow a video in infinity
     *
     * @var array
     *
     * @Column(type="json_array", nullable=true)
     */
    protected $images = [];

    /**
     * @var int
     *
     * @Column(type="integer", nullable=true)
     */
    protected $bitRate;

    /**
     * The encoding used
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $encoding;

    /**
     * Get the audio metadata
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $audio;

    /**
     * Get the video metadata
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $video;

    /**
     * @var string
     *
     * @Column(nullable=true, options={"collation": "utf8_bin"})
     */
    protected $videoUrl;

    /**
     * @var boolean
     *
     * @Column(type="boolean")
     */
    protected $videoUrlValid = false;

    /**
     * @var ValueObject\Images
     *
     * @Embedded(class="Aphrodite\Recording\Entity\ValueObject\Images")
     */
    protected $imageUrls;

    /**
     * @var boolean
     *
     * @Column(type="boolean")
     */
    protected $imageUrlsValid = false;

    /**
     * @var Site[]
     *
     * @OneToMany(targetEntity="Aphrodite\Site\Entity\PostAssociation", mappedBy="recording", fetch="EAGER")
     */
    protected $publishedOn;

    /**
     * @var bool
     *
     * @Column(type="integer")
     */
    protected $orphaned = false;

    /**
     * @var string|null
     *
     * @Column(type="string", length=30, nullable=true)
     */
    protected $upstoreHash = null;

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
     * @var DateTime
     *
     * @Column(type="datetime", nullable=true)
     */
    protected $lastCheckedAt;

    /**
     * @var AbstractPerformerEntity
     *
     * @ORM\ManyToOne(targetEntity="Aphrodite\Performer\Entity\AbstractPerformerEntity", inversedBy="recordings")
     * @ORM\JoinColumn(nullable=true)
     */
    protected $performer = null;

    public function __construct()
    {
        $this->imageUrls   = new Images(null, null, null);
        $this->publishedOn = new ArrayCollection();
    }

    /**
     * @return AbstractPerformerEntity
     */
    public function getPerformer():? AbstractPerformerEntity
    {
        return $this->performer;
    }

    /**
     * @param AbstractPerformerEntity $performer
     */
    public function setPerformer(AbstractPerformerEntity $performer = null): void
    {
        $this->performer = $performer;
    }

    public function getImages(): array
    {
        return $this->images;
    }

    public function setImages(array $images): void
    {
        $this->images = $images;
    }

    /**
     * @return Site[]
     */
    public function getPublishedOn(): array
    {
        return $this->publishedOn->toArray();
    }

    /**
     * @param Site[] $publishedOn
     */
    public function setPublishedOn(array $publishedOn): void
    {
        $this->publishedOn->clear();

        array_map([$this->publishedOn, 'add'], $publishedOn);
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * @param int $id
     */
    public function setId(int $id): void
    {
        $this->id = $id;
    }

    /**
     * @return int
     */
    public function getOldId():? int
    {
        return $this->oldId;
    }

    /**
     * @param int $oldId
     */
    public function setOldId(int $oldId = null)
    {
        $this->oldId = $oldId;
    }

    /**
     * @return boolean
     */
    public function isOrphaned(): bool
    {
        return (bool) $this->orphaned;
    }

    /**
     * @param bool $orphaned
     */
    public function setOrphaned(bool $orphaned): void
    {
        $this->orphaned = $orphaned;
    }

    /**
     * @return string
     */
    public function getStageName(): string
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
     * @return string
     */
    public function getSection():? string
    {
        return $this->section;
    }

    /**
     * @param string $section
     */
    public function setSection(string $section = null): void
    {
        $this->section = $section;
    }

    /**
     * @return string
     */
    public function getService(): string
    {
        return $this->service;
    }

    /**
     * @param string $service
     */
    public function setService(string $service): void
    {
        $this->service = $service;
    }

    /**
     * @return boolean
     */
    public function isImageUrlsValid(): bool
    {
        return $this->imageUrlsValid;
    }

    /**
     * @param boolean $imageUrlsValid
     */
    public function setImageUrlsValid(bool $imageUrlsValid): void
    {
        $this->imageUrlsValid = $imageUrlsValid;
    }

    /**
     * @return boolean
     */
    public function isVideoUrlValid(): bool
    {
        return $this->videoUrlValid;
    }

    /**
     * @param boolean $videoUrlValid
     */
    public function setVideoUrlValid(bool $videoUrlValid): void
    {
        $this->videoUrlValid = $videoUrlValid;
    }

    /**
     * @return int
     */
    public function getDuration():? int
    {
        return $this->duration;
    }

    /**
     * @param int $duration
     */
    public function setDuration(int $duration): void
    {
        $this->duration = (int) $duration;
    }

    /**
     * @return int
     */
    public function getSize264():? int
    {
        // Doctrine stores size264 as a string to keep x86 compat
        if ($this->size264 !== null) {
            return (int) $this->size264;
        }

        return null;
    }

    /**
     * @param int $size264
     */
    public function setSize264(int $size264 = null): void
    {
        $this->size264 = $size264;
    }

    /**
     * @return int
     */
    public function getSize265():? int
    {
        // Doctrine stores size265 as a string to keep x86 compat
        if ($this->size265 !== null) {
            return (int) $this->size265;
        }

        return null;
    }

    /**
     * @param int $size265
     */
    public function setSize265(int $size265 = null)
    {
        $this->size265 = $size265;
    }

    /**
     * @return string
     */
    public function getVideoMp4Uuid(): ?string
    {
        return $this->videoMp4Uuid;
    }

    /**
     * @param string $videoMp4Uuid
     */
    public function setVideoMp4Uuid(string $videoMp4Uuid = null): void
    {
        if ($videoMp4Uuid === '00000000-0000-0000-0000-000000000000') {
            $videoMp4Uuid = null;
        }

        $this->videoMp4Uuid = $videoMp4Uuid;
    }

    /**
     * @return string
     */
    public function getVideoHlsUuid(): ?string
    {
        return $this->videoHlsUuid;
    }

    /**
     * @param string $videoHlsUuid
     */
    public function setVideoHlsUuid(string $videoHlsUuid = null): void
    {
        $this->videoHlsUuid = $videoHlsUuid;
    }

    /**
     * @return string
     */
    public function getVideoManifest(): ?string
    {
        return $this->videoManifest;
    }

    /**
     * @param string $videoManifest
     */
    public function setVideoManifest(string $videoManifest = null): void
    {
        $this->videoManifest = $videoManifest;
    }

    /**
     * @return string
     */
    public function getWordpressCollageUuid():? string
    {
        return $this->wordpressCollageUuid;
    }

    /**
     * @param string $wordpressCollageUuid
     */
    public function setWordpressCollageUuid(string $wordpressCollageUuid = null): void
    {
        $this->wordpressCollageUuid = $wordpressCollageUuid;
    }

    /**
     * @return string
     */
    public function getInfinityCollageUuid():? string
    {
        return $this->infinityCollageUuid;
    }

    /**
     * @param string $infinityCollageUuid
     */
    public function setInfinityCollageUuid(string $infinityCollageUuid = null): void
    {
        $this->infinityCollageUuid = $infinityCollageUuid;
    }

    /**
     * @return string[]
     */
    public function getSprites(): array
    {
        return $this->sprites;
    }

    /**
     * @param string[] $sprites
     */
    public function setSprites(array $sprites = []): void
    {
        $this->sprites = $sprites;
    }

    /**
     * @return string
     */
    public function getVideoUrl(): ?string
    {
        return $this->videoUrl;
    }

    /**
     * @param string $videoUrl
     */
    public function setVideoUrl(string $videoUrl = null): void
    {
        $this->videoUrl = $videoUrl;
    }

    /**
     * @return ValueObject\Images
     */
    public function getImageUrls(): ValueObject\Images
    {
        return $this->imageUrls;
    }

    /**
     * @param ValueObject\Images $imageUrls
     */
    public function setImageUrls(ValueObject\Images $imageUrls): void
    {
        $this->imageUrls = $imageUrls;
    }

    /**
     * @return DateTime
     */
    public function getCreatedAt(): ?DateTime
    {
        return $this->createdAt;
    }

    /**
     * @param DateTime $createdAt
     */
    public function setCreatedAt(DateTime $createdAt): void
    {
        $this->createdAt = $createdAt;
    }

    /**
     * @return DateTime
     */
    public function getUpdatedAt(): ?DateTime
    {
        return $this->updatedAt;
    }

    /**
     * @param DateTime $updatedAt
     */
    public function setUpdatedAt(DateTime $updatedAt): void
    {
        $this->updatedAt = $updatedAt;
    }

    /**
     * @return string
     */
    public function getBitRate():? int
    {
        if ($this->bitRate !== null) {
            return (int) $this->bitRate;
        }

        return null;
    }

    /**
     * @param string $bitRate
     */
    public function setBitRate(string $bitRate = null)
    {
        $this->bitRate = $bitRate;
    }

    /**
     * @return string
     */
    public function getEncoding():? string
    {
        return $this->encoding;
    }

    /**
     * @param string $encoding
     */
    public function setEncoding(string $encoding = null): void
    {
        $this->encoding = $encoding;
    }

    /**
     * @return string
     */
    public function getAudio():? string
    {
        return $this->audio;
    }

    /**
     * @param string $audio
     */
    public function setAudio(string $audio = null): void
    {
        $this->audio = $audio;
    }

    /**
     * @return string
     */
    public function getVideo():? string
    {
        return $this->video;
    }

    /**
     * @param string $video
     */
    public function setVideo(string $video = null): void
    {
        $this->video = $video;
    }

    /**
     * @return string
     */
    public function getState(): string
    {
        return $this->state;
    }

    /**
     * @param string $state
     */
    public function setState(string $state): void
    {
        $this->state = $state;
    }

    /**
     * @return DateTime
     */
    public function getLastCheckedAt():? DateTime
    {
        return $this->lastCheckedAt;
    }

    /**
     * @param DateTime $lastCheckedAt
     */
    public function setLastCheckedAt(DateTime $lastCheckedAt = null): void
    {
        $this->lastCheckedAt = $lastCheckedAt;
    }

    /**
     * @return null|string
     */
    public function getUpstoreHash():? string
    {
        return $this->upstoreHash;
    }

    /**
     * @param null|string $upstoreHash
     */
    public function setUpstoreHash(string $upstoreHash = null): void
    {
        $this->upstoreHash = $upstoreHash;
    }
}
