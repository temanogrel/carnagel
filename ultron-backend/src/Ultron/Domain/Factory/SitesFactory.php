<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Factory;

use Doctrine\Common\Collections\ArrayCollection;
use Interop\Container\ContainerInterface;
use Ultron\Domain\SiteConfiguration;
use Ultron\Domain\Sites;

class SitesFactory
{
    public function __invoke(ContainerInterface $container)
    {
        /* @var $config array */
        $config = $container->get('config');
        $config = $config['ultron'] ?? [];

        $configurations = new ArrayCollection();

        foreach ($config['sites'] ?? [] as $site) {
            $configurations->set($site['domain'], new SiteConfiguration($site));
        }

        $sites = new Sites();
        $sites->setApiAccessToken($config['apiToken']);
        $sites->setSiteConfigurations($configurations);

        return $sites;
    }
}
