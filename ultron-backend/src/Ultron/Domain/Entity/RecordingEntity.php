<?php
/**
 *
 *
 *
 */

namespace Ultron\Domain\Entity;

use DateTime;
use Doctrine\ORM\Mapping\Cache;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Embedded;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Index;
use Doctrine\ORM\Mapping\JoinColumn;
use Doctrine\ORM\Mapping\ManyToOne;
use Doctrine\ORM\Mapping\Table;
use InvalidArgumentException;
use Ultron\Domain\Service\RecordingService;

/**
 * @Table(name="recordings", indexes={
 *  @Index(name="uid", columns={"uid"}),
 *  @Index(name="slug", columns={"slug"}),
 *  @Index(name="stageName", columns={"stageName"}),
 *  @Index(name="createdAt", columns={"createdAt"}),
 *  @Index(name="galleryUrl", columns={"imageUrls_galleryUrl"}),
 *  @Index(name="updatedAt", columns={"updatedAt"})
 * })
 *
 * @Entity(repositoryClass="Ultron\Infrastructure\Repository\RecordingRepository")
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
     * Internal id
     *
     * @var string
     *
     * @Column(type="integer", unique=true)
     */
    protected $uid;

    /**
     * @var string
     *
     * @Column(unique=true)
     */
    protected $slug;

    /**
     * The stage name used when recording the video.
     *
     * @var string
     *
     * @Column()
     */
    protected $stageName;

    /**
     * @var integer
     *
     * @Column(type="integer")
     */
    protected $duration;

    /**
     * @var integer
     *
     * @Column(type="bigint")
     */
    protected $size;

    /**
     * @var string
     *
     * @Column()
     */
    protected $bitRate;

    /**
     * @var string
     *
     * @Column()
     */
    protected $audio;

    /**
     * @var string
     *
     * @Column()
     */
    protected $video;

    /**
     * @var string
     *
     * @Column()
     */
    protected $videoUrl;

    /**
     * @var ValueObject\Images
     *
     * @Embedded(class="Ultron\Domain\Entity\ValueObject\Images")
     */
    protected $imageUrls;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $views = 0;

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
     * @var PerformerEntity
     *
     * @ManyToOne(targetEntity="Ultron\Domain\Entity\PerformerEntity", inversedBy="recordings", fetch="EAGER")
     * @JoinColumn(onDelete="CASCADE")
     */
    protected $performer;

    /**
     * Convert seconds to a human readable duration
     *
     * @param string $format
     *
     * @return string
     */
    public function getFormattedDuration(string $format = 'H:i:s'): string
    {
        return gmdate($format, $this->duration);
    }

    /**
     * Get the post title
     *
     * @return string
     */
    public function getRecordingTitle():string
    {
        return RecordingService::generatePostTitle($this);
    }

    /**
     * @return int
     */
    public function getId():int
    {
        return $this->id;
    }

    /**
     * @param int $id
     */
    public function setId($id)
    {
        $this->id = (int) $id;
    }

    /**
     * @return int
     */
    public function getUid():int
    {
        return $this->uid;
    }

    /**
     * @param int $uid
     */
    public function setUid(int $uid)
    {
        $this->uid = $uid;
    }

    /**
     * @return mixed
     */
    public function getSlug()
    {
        return $this->slug;
    }

    /**
     * @param mixed $slug
     */
    public function setSlug(string $slug)
    {
        $this->slug = $slug;
    }

    /**
     * @return string
     */
    public function getStageName():string
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
     * @return int
     */
    public function getDuration():int
    {
        return $this->duration;
    }

    /**
     * @param int $duration
     */
    public function setDuration(int $duration)
    {
        $this->duration = $duration;
    }

    /**
     * @return int
     */
    public function getSize():int
    {
        return $this->size;
    }

    /**
     * @param int $size
     */
    public function setSize(int $size)
    {
        $this->size = $size;
    }

    /**
     * @return string
     */
    public function getVideoUrl():string
    {
        return $this->videoUrl;
    }

    /**
     * @param string $videoUrl
     */
    public function setVideoUrl(string $videoUrl)
    {
        $this->videoUrl = $videoUrl;
    }

    /**
     * @return ValueObject\Images
     */
    public function getImageUrls():ValueObject\Images
    {
        return $this->imageUrls;
    }

    /**
     * @param ValueObject\Images $imageUrls
     */
    public function setImageUrls(ValueObject\Images $imageUrls)
    {
        $this->imageUrls = $imageUrls;
    }

    /**
     * @return DateTime
     */
    public function getCreatedAt():DateTime
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
    public function getUpdatedAt():DateTime
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
     * @return PerformerEntity
     */
    public function getPerformer():PerformerEntity
    {
        return $this->performer;
    }

    /**
     * @param PerformerEntity $performer
     */
    public function setPerformer(PerformerEntity $performer)
    {
        $this->performer = $performer;
    }

    /**
     * @return string
     */
    public function getBitRate():string
    {
        return $this->bitRate;
    }

    /**
     * @param string $bitRate
     */
    public function setBitRate(string $bitRate)
    {
        $this->bitRate = $bitRate;
    }

    /**
     * @return string
     */
    public function getAudio():string
    {
        return $this->audio;
    }

    /**
     * @param string $audio
     */
    public function setAudio(string $audio)
    {
        $this->audio = $audio;
    }

    /**
     * @return string
     */
    public function getVideo():string
    {
        return $this->video;
    }

    /**
     * @param string $video
     */
    public function setVideo(string $video)
    {
        $this->video = $video;
    }
}
