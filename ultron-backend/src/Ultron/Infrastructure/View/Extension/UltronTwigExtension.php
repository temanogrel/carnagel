<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\View\Extension;

use Twig_Extension;
use Twig_Extension_GlobalsInterface;
use Twig_SimpleFunction;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\ArrayUtils;

final class UltronTwigExtension extends Twig_Extension implements Twig_Extension_GlobalsInterface
{
    /**
     * @var Sites
     */
    private $sites;

    /**
     * UltronTwigExtension constructor.
     *
     * @param Sites $sites
     */
    public function __construct(Sites $sites)
    {
        $this->sites = $sites;
    }

    /**
     * {@inheritdoc}
     */
    public function getName():string
    {
        return 'ultron';
    }

    public function getFunctions():array
    {
        return [
            new Twig_SimpleFunction('unique', [$this, 'filterUnique']),
        ];
    }

    public function getGlobals():array
    {
        return [
            'site'          => $this->sites->getCurrentSite(),
            'base_keywords' => [
                'camgirl',
                'webcam',
                'videos',
                'upstore',
                'download',
                'myfreecams',
                'chaturbate',
                'cam4',
                'filehost',
            ],
        ];
    }

    public function filterUnique(...$args):array
    {
        return ArrayUtils::iunique(array_merge(...array_map([ArrayUtils::class, 'iteratorToArray'], $args)));
    }
}
