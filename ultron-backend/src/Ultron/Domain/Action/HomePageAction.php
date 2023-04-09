<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;

use Doctrine\Common\Collections\Criteria;
use Exception;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use RecursiveArrayIterator;
use RecursiveIteratorIterator;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\RecordingService;
use Ultron\Domain\SiteConfiguration;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Ultron\Infrastructure\Service\PageCacheServiceInterface;
use Zend\Diactoros\Response\HtmlResponse;
use Zend\Expressive\Router;
use Zend\Expressive\Template;
use Zend\Expressive\Template\TemplateRendererInterface;

class HomePageAction
{

    /**
     * @var SiteConfiguration
     */
    private $site;

    /**
     * @var TemplateRendererInterface
     */
    private $template;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var PageCacheServiceInterface
     */
    private $pageCacheService;

    /**
     * HomePageAction constructor.
     *
     * @param SiteConfiguration $site
     * @param TemplateRendererInterface $template
     * @param RecordingRepositoryInterface $recordingRepository
     * @param PageCacheServiceInterface $pageCacheService
     */
    public function __construct(
        SiteConfiguration $site,
        TemplateRendererInterface $template,
        RecordingRepositoryInterface $recordingRepository,
        PageCacheServiceInterface $pageCacheService
    ) {
        $this->site                = $site;
        $this->template            = $template;
        $this->recordingRepository = $recordingRepository;
        $this->pageCacheService    = $pageCacheService;
    }

    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next = null)
    {
        $query = $request->getQueryParams();

        $criteria = Criteria::create();
        $criteria->orderBy(['createdAt' => 'DESC']);

        $page = (int)($query['page'] ?? 1);

        $paginator = $this->pageCacheService->get($criteria, $this->site, $page);
        $paginator->setCurrentPageNumber($page);
        $paginator->setPageRange(5);

        $recordingKeywords = array_map(function (RecordingEntity $recording) {
            return RecordingService::getRecordingKeywords($recording);
        }, iterator_to_array($paginator));

        // Create a recursive iterator so we can merge the recording keywords
        $iterator = new RecursiveIteratorIterator(new RecursiveArrayIterator($recordingKeywords));
        $keywords = iterator_to_array($iterator, false);

        $descArgs = [
            $page,
            $paginator->count(),
        ];

        return new HtmlResponse($this->template->render('app::list', [
            'site'       => $this->site,
            'recordings' => $paginator,

            // SEO
            'keywords'   => $keywords,
            'description' => vsprintf('Page %d of %d for most recent recordings', $descArgs),
        ]));
    }
}
