<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service\Sitemap;

use DOMDocument;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\SiteConfiguration;
use Zend\Expressive\Helper\UrlHelper;

class PerformerWriter
{
    const MAX_SIZE = 50000;
    const FILE_PATTERN = 'public/sitemap/%s-performers-%d.xml';

    /**
     * @var array
     */
    private $performers = [];

    /**
     * @var int
     */
    private $performerCount = 0;

    /**
     * @var int
     */
    private $iterations = 0;

    /**
     * @var SiteConfiguration
     */
    private $site;

    /**
     * @var UrlHelper
     */
    private $urlHelper;

    /**
     * @param SiteConfiguration $site
     * @param UrlHelper   $urlHelper
     */
    public function __construct(SiteConfiguration $site, UrlHelper $urlHelper)
    {
        $this->urlHelper = $urlHelper;
        $this->site      = $site;
    }

    private function reset()
    {
        $this->performers     = [];
        $this->performerCount = 0;
    }

    public function flush()
    {
        $doc = new DOMDocument('1.0', 'utf-8');

        // root element
        $urlSet = $doc->createElement('urlset');
        $urlSet->setAttribute('xmlns', 'http://www.sitemaps.org/schemas/sitemap/0.9');

        // Set the proper hostname context for the url generator
        $this
            ->urlHelper
            ->setBasePath($this->site->getDomain());

        /* @var $performer PerformerEntity */
        foreach ($this->performers as $performer) {

            $urlArgs = [
                'slug'   => $performer->getSlug(),
                'prefix' => $this->site->getUrlRoot()
            ];

            $url       = 'http:/' . $this->urlHelper->generate('performer.list-recordings', $urlArgs);
            $updatedAt = $performer->getUpdatedAt()->format('Y-m-d');

            $el = $doc->createElement('url');
            $el->appendChild($doc->createElement('loc', $url));
            $el->appendChild($doc->createElement('lastmod', $updatedAt));
            $el->appendChild($doc->createElement('changefreq', 'monthly'));
            $el->appendChild($doc->createElement('priority', (string) 0.4));

            $urlSet->appendChild($el);
        }

        $file = sprintf(static::FILE_PATTERN, $this->site->getDomain(), $this->iterations++);

        $doc->appendChild($urlSet);
        $doc->save($file);

        $this->reset();
    }

    public function add(PerformerEntity $performer)
    {
        if (!$performer->belongsTo($this->site)) {
            return;
        }

        $this->performers[] = $performer;
        $this->performerCount++;

        if ($this->performerCount >= static::MAX_SIZE) {
            $this->flush();
        }
    }
}
