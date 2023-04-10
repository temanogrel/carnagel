<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Entity\DeathFile;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Entity\DeathFileEntity;
use DateTime;
use Doctrine\ORM\Mapping as ORM;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Index;
use Doctrine\ORM\Mapping\JoinColumn;
use Doctrine\ORM\Mapping\ManyToOne;
use Doctrine\ORM\Mapping\OneToOne;
use Doctrine\ORM\Mapping\Table;

/**
 * Class UrlEntry
 *
 * @Entity(repositoryClass="Aphrodite\Recording\Repository\DeathFile\UrlRepository")
 * @Table(
 *     name="death_file_entries",
 *     options={"collate": "utf8_bin"},
 *     indexes={
 *          @Index(columns={"filename"})
 *     }
 * )
 */
class UrlEntry
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
     * Upstore url
     *
     * @var string
     *
     * @Column(unique=true)
     */
    protected $url;

    /**
     * File name of the upstore recording
     *
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $filename;

    /**
     * @var integer
     *
     * @Column(type="integer", nullable=true)
     */
    protected $hermesId;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $state;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $ignoreReason;

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
     * @var RecordingEntity
     *
     * @OneToOne(targetEntity="Aphrodite\Recording\Entity\RecordingEntity")
     * @JoinColumn(onDelete="CASCADE", nullable=true)
     */
    protected $recording;

    /**
     * @var DeathFileEntity
     *
     * @ManyToOne(targetEntity="Aphrodite\Recording\Entity\DeathFileEntity")
     * @JoinColumn(onDelete="CASCADE")
     */
    protected $deathFile;

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
     * @return int
     */
    public function getHermesId()
    {
        return $this->hermesId;
    }

    /**
     * @param int $hermesId
     */
    public function setHermesId($hermesId)
    {
        $this->hermesId = $hermesId;
    }

    /**
     * @return string
     */
    public function getUrl()
    {
        return $this->url;
    }

    /**
     * @param string $url
     */
    public function setUrl($url)
    {
        $this->url = $url;
    }

    /**
     * @return string
     */
    public function getFilename()
    {
        return $this->filename;
    }

    /**
     * @param string $filename
     */
    public function setFilename(string $filename)
    {
        $this->filename = $filename;
    }

    /**
     * @return string
     */
    public function getState()
    {
        return $this->state;
    }

    /**
     * @param string $state
     */
    public function setState($state)
    {
        $this->state = $state;
    }

    /**
     * @return string
     */
    public function getIgnoreReason()
    {
        return $this->ignoreReason;
    }

    /**
     * @param string $ignoreReason
     */
    public function setIgnoreReason($ignoreReason)
    {
        $this->ignoreReason = $ignoreReason;
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
     * @return RecordingEntity
     */
    public function getRecording()
    {
        return $this->recording;
    }

    /**
     * @param RecordingEntity $recording
     */
    public function setRecording(RecordingEntity $recording = null)
    {
        $this->recording = $recording;
    }

    /**
     * @return DeathFileEntity
     */
    public function getDeathFile()
    {
        return $this->deathFile;
    }

    /**
     * @param DeathFileEntity $deathFile
     */
    public function setDeathFile(DeathFileEntity $deathFile)
    {
        $this->deathFile = $deathFile;
    }
}
