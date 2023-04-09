<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;


use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Ultron\Domain\SiteConfiguration;
use Ultron\Infrastructure\Http\Response\XmlResponse;
use Ultron\Infrastructure\Service\SitemapService;
use Zend\Expressive\Template\TemplateRendererInterface;

final class SiteMapIndexAction
{
    /**
     * @var SitemapService
     */
    private $siteMapService;

    /**
     * @var TemplateRendererInterface
     */
    private $template;

    /**
     * @var SiteConfiguration
     */
    private $site;

    /**
     * SitemapController constructor.
     *
     * @param SiteConfiguration         $site
     * @param SitemapService            $siteMapService
     * @param TemplateRendererInterface $template
     */
    public function __construct(
        SiteConfiguration $site,
        SitemapService $siteMapService,
        TemplateRendererInterface $template
    ) {
        $this->site           = $site;
        $this->template       = $template;
        $this->siteMapService = $siteMapService;
    }

    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        $args = [
            'urls' => $this->siteMapService->getSitemapUrls($this->site)
        ];

        return new XmlResponse($this->template->render('app::sitemap', $args));
    }
}
