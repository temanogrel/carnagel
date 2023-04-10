<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Stdlib;

use Doctrine\ORM\Mapping\Column;

trait SectionAwareTrait
{
    /**
     * @var string
     *
     * @Column(nullable=true)
     */
    protected $section;

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
    public function setSection($section)
    {
        $this->section = (string) $section;
    }
}
