<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Entity;

use Aphrodite\Recording\Entity\RecordingEntity;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\JoinColumn;
use Doctrine\ORM\Mapping\ManyToOne;

/**
 * Class PostAssociation
 *
 * @Entity(repositoryClass="Aphrodite\Site\Repository\PostAssociationRepository")
 */
class PostAssociation
{
    /**
     * @var int
     *
     * @Id()
     * @Column(type="integer")
     * @GeneratedValue()
     */
    protected $id;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $postId;

    /**
     * @var RecordingEntity
     *
     * @ManyToOne(targetEntity="Aphrodite\Recording\Entity\RecordingEntity", inversedBy="publishedOn")
     * @JoinColumn(onDelete="CASCADE")
     */
    protected $recording;

    /**
     * @var Site
     *
     * @ManyToOne(targetEntity="Site", inversedBy="posts")
     * @JoinColumn(onDelete="CASCADE")
     */
    protected $site;

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
    public function getPostId()
    {
        return $this->postId;
    }

    /**
     * @param int $postId
     */
    public function setPostId($postId)
    {
        $this->postId = $postId;
    }

    /**
     * @return Site
     */
    public function getSite()
    {
        return $this->site;
    }

    /**
     * @param Site $site
     */
    public function setSite(Site $site)
    {
        $this->site = $site;
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
    public function setRecording(RecordingEntity $recording)
    {
        $this->recording = $recording;
    }
}
