<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Entity\ValueObject;

use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Embeddable;

/**
 * Class Images
 *
 * @Embeddable()
 */
class Images
{
    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $thumb;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $large;

    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $gallery;

    /**
     * Images constructor.
     *
     * @param string $thumb
     * @param string $large
     * @param string $gallery
     */
    public function __construct($thumb = null, $large = null, $gallery = null)
    {
        $this->thumb   = $thumb;
        $this->large   = $large;
        $this->gallery = $gallery;
    }

    /**
     * @param string $thumb
     */
    public function setThumb($thumb)
    {
        $this->thumb = $thumb;
    }

    /**
     * @param string $large
     */
    public function setLarge($large)
    {
        $this->large = $large;
    }

    /**
     * @param string $gallery
     */
    public function setGallery($gallery)
    {
        $this->gallery = $gallery;
    }

    /**
     * @return string
     */
    public function getThumb()
    {
        return $this->thumb;
    }

    /**
     * @return string
     */
    public function getLarge()
    {
        return $this->large;
    }

    /**
     * @return string
     */
    public function getGallery()
    {
        return $this->gallery;
    }
}
