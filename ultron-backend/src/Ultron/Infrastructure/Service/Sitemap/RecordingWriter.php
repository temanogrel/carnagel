<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service\Sitemap;

use DOMDocument;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\RecordingService;
use Ultron\Domain\SiteConfiguration;
use Zend\Expressive\Helper\UrlHelper;

class RecordingWriter
{
    const MAX_SIZE = 50000;
    const FILE_PATTERN = 'public/sitemap/%s-recordings-%d.xml';

    /**
     * @var array
     */
    private $recordings = [];

    /**
     * @var int
     */
    private $recordingCount = 0;

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
     * @param UrlHelper $urlHelper
     */
    public function __construct(SiteConfiguration $site, UrlHelper $urlHelper)
    {
        $this->urlHelper = $urlHelper;
        $this->site      = $site;
    }

    private function reset()
    {
        $this->recordings     = [];
        $this->recordingCount = 0;
    }

    public function flush()
    {
        $doc = new DOMDocument('1.0', 'utf-8');

        // root element
        $urlSet = $doc->createElement('urlset');
        $urlSet->setAttribute('xmlns', 'http://www.sitemaps.org/schemas/sitemap/0.9');
        $urlSet->setAttribute('xmlns:image', 'http://www.google.com/schemas/sitemap-image/1.1');

        $this
            ->urlHelper
            ->setBasePath($this->site->getDomain());

        /* @var $recording RecordingEntity */
        foreach ($this->recordings as $recording) {

            $urlArgs = [
                'slug'   => $recording->getSlug(),
                'prefix' => $this->site->getUrlRoot()
            ];

            $url       = 'http:/' . $this->urlHelper->generate('recording.details', $urlArgs);
            $updatedAt = $recording->getUpdatedAt()->format('Y-m-d');

            $image = $doc->createElement('image:image');
            $image->appendChild($doc->createElement('image:loc', $recording->getImageUrls()->getLarge()));
            $image->appendChild($doc->createElement('image:title', RecordingService::generatePostTitle($recording)));

            $el = $doc->createElement('url');
            $el->appendChild($doc->createElement('loc', $url));
            $el->appendChild($doc->createElement('lastmod', $updatedAt));
            $el->appendChild($doc->createElement('changefreq', 'monthly'));
            $el->appendChild($doc->createElement('priority', (string) 0.8));
            $el->appendChild($image);

            $urlSet->appendChild($el);
        }

        $file = sprintf(static::FILE_PATTERN, $this->site->getDomain(), $this->iterations++);

        $doc->appendChild($urlSet);
        $doc->save($file);

        $this->reset();
    }

    public function add(RecordingEntity $recording)
    {
        if (!$recording->getPerformer()->belongsTo($this->site)) {
            return;
        }

        $this->recordings[] = $recording;
        $this->recordingCount++;

        if ($this->recordingCount >= static::MAX_SIZE) {
            $this->flush();
        }
    }
}
