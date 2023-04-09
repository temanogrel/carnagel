<?php
/**
 *
 *
 *
 */

namespace Ultron\Domain;

use Zend\Stdlib\AbstractOptions;

class SiteConfiguration extends AbstractOptions
{
    /**
     * @var bool
     */
    protected $enabled;

    /**
     * @var string
     */
    protected $urlRoot;

    /**
     * @var string
     */
    protected $domain;

    /**
     * @var string
     */
    protected $service;

    /**
     * @var string
     */
    protected $section;

    /**
     * @var string
     */
    protected $theme;

    /**
     * @var string
     */
    protected $favicon;

    /**
     * @var int
     */
    protected $pageSize = 90;

    /**
     * @return boolean
     */
    public function isEnabled()
    {
        return $this->enabled;
    }

    /**
     * @param boolean $enabled
     */
    public function setEnabled($enabled)
    {
        $this->enabled = (bool) $enabled;
    }

    /**
     * @return string
     */
    public function getUrlRoot()
    {
        return $this->urlRoot;
    }

    /**
     * @param string $urlRoot
     */
    public function setUrlRoot($urlRoot)
    {
        $this->urlRoot = $urlRoot;
    }

    /**
     * @return string
     */
    public function getService()
    {
        return $this->service;
    }

    /**
     * @param string $service
     */
    public function setService($service)
    {
        $this->service = $service;
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
    public function setSection($section)
    {
        $this->section = $section;
    }

    /**
     * @return string
     */
    public function getTheme()
    {
        return $this->theme;
    }

    /**
     * @param string $theme
     */
    public function setTheme($theme)
    {
        $this->theme = $theme;
    }

    /**
     * @return string
     */
    public function getFavicon()
    {
        return $this->favicon;
    }

    /**
     * @param string $favicon
     */
    public function setFavicon($favicon)
    {
        $this->favicon = $favicon;
    }

    /**
     * @return string
     */
    public function getDomain()
    {
        return $this->domain;
    }

    /**
     * @param string $domain
     */
    public function setDomain($domain)
    {
        $this->domain = $domain;
    }

    /**
     * @return int
     */
    public function getPageSize(): int
    {
        return $this->pageSize;
    }

    /**
     * @param int $pageSize
     */
    public function setPageSize(int $pageSize)
    {
        $this->pageSize = $pageSize;
    }

}
