<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain;

use Doctrine\Common\Collections\ArrayCollection;
use Zend\Stdlib\AbstractOptions;

class Sites extends AbstractOptions
{
    /**
     * @var string
     */
    protected $apiAccessToken = 'hello';

    /**
     * @var string
     */
    protected $defaultHostname = 'ultron.dev';

    /**
     * @var SiteConfiguration
     */
    protected $currentSite;

    /**
     * @var ArrayCollection|SiteConfiguration[]
     */
    protected $siteConfigurations;

    /**
     * @return string
     */
    public function getApiAccessToken():string
    {
        return $this->apiAccessToken;
    }

    /**
     * @param string $apiAccessToken
     */
    public function setApiAccessToken(string $apiAccessToken)
    {
        $this->apiAccessToken = $apiAccessToken;
    }

    /**
     * @return string
     */
    public function getDefaultHostname():string
    {
        return $this->defaultHostname;
    }

    /**
     * @param string $defaultHostname
     */
    public function setDefaultHostname(string $defaultHostname)
    {
        $this->defaultHostname = $defaultHostname;
    }

    /**
     * @return SiteConfiguration
     */
    public function getCurrentSite():SiteConfiguration
    {
        return $this->currentSite;
    }

    /**
     * @param SiteConfiguration $currentSite
     */
    public function setCurrentSite(SiteConfiguration $currentSite)
    {
        $this->currentSite = $currentSite;
    }

    /**
     * @return ArrayCollection|SiteConfiguration[]
     */
    public function getSiteConfigurations():ArrayCollection
    {
        return $this->siteConfigurations;
    }

    /**
     * @param ArrayCollection|SiteConfiguration[] $siteConfigurations
     */
    public function setSiteConfigurations(ArrayCollection $siteConfigurations)
    {
        $this->siteConfigurations = $siteConfigurations;
    }
}
