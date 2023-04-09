<?php
/**
 *
 *
 *
 */

namespace Ultron\Domain\Entity\ValueObject;

use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Embeddable;
use Zend\Stdlib\ParametersInterface;

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
     * @Column()
     */
    protected $thumb;

    /**
     * @var string
     *
     * @Column()
     */
    protected $large;

    /**
     * @var string
     *
     * @Column(options={"collation":"utf8_bin"})
     */
    protected $galleryUrl;

    /**
     * Images constructor.
     *
     * @param string $thumb
     * @param string $large
     * @param string $galleryUrl
     */
    public function __construct($thumb, $large, $galleryUrl)
    {
        $this->thumb      = $thumb;
        $this->large      = $large;
        $this->galleryUrl = $galleryUrl;
    }

    public static function fromParameters(ParametersInterface $parameters):self
    {
        return new self(
            $parameters->get('thumb'),
            $parameters->get('large'),
            $parameters->get('gallery')
        );
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
    public function getGalleryUrl()
    {
        return $this->galleryUrl;
    }
}
