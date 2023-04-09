<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\View\Extension;

use Interop\Container\ContainerInterface;
use Ultron\Domain\Sites;

class UltronTwigExtensionFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return UltronTwigExtension
     */
    public function __invoke(ContainerInterface $container): UltronTwigExtension
    {
        /* @var $sites Sites */
        $sites = $container->get(Sites::class);

        return new UltronTwigExtension($sites);
    }
}
