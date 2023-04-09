<?php
/**
 *
 *
 *
 */

declare(strict_types=1);

namespace Ultron\Domain\Action;

use Doctrine\Common\Collections\Criteria;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use RecursiveArrayIterator;
use RecursiveIteratorIterator;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Service\RecordingService;
use Ultron\Domain\SiteConfiguration;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Zend\Diactoros\Response\HtmlResponse;
use Zend\Diactoros\Response\RedirectResponse;
use Zend\Expressive\Helper\UrlHelper;
use Zend\Expressive\Template\TemplateRendererInterface;

class RecordingSearchAction
{
    /**
     * @var UrlHelper
     */
    private $urlHelper;

    /**
     * @var TemplateRendererInterface
     */
    private $template;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;
    /**
     * @var SiteConfiguration
     */
    private $site;

    /**
     * RecordingSearchAction constructor.
     *
     * @param UrlHelper $urlHelper
     * @param SiteConfiguration $site
     * @param TemplateRendererInterface $template
     * @param RecordingRepositoryInterface $recordingRepository
     */
    public function __construct(
        UrlHelper $urlHelper,
        SiteConfiguration $site,
        TemplateRendererInterface $template,
        RecordingRepositoryInterface $recordingRepository
    ) {
        $this->template            = $template;
        $this->urlHelper           = $urlHelper;
        $this->recordingRepository = $recordingRepository;
        $this->site                = $site;
    }

    /**
     * @param ServerRequestInterface $request
     * @param ResponseInterface $response
     * @param callable $next
     *
     * @return HtmlResponse
     */
    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        if ($request->getMethod() !== 'GET') {
            $query = $request->getParsedBody()['query'] ?? '';

            if (empty($query)) {
                return new RedirectResponse($this->urlHelper->generate('home'));
            }

            return new RedirectResponse(
                $this->urlHelper->generate('recording.search', ['query' => $query])
            );
        }

        $page  = (int) ($request->getQueryParams()['page'] ?? 1);
        $query = $request->getAttribute('query');

        $criteria = Criteria::create();
        $criteria->orderBy(['createdAt' => 'desc']);

        $paginator = $this->recordingRepository->searchByPerformer($query, $criteria, $this->site);
        $paginator->setPageRange(5);
        $paginator->setCurrentPageNumber($page);

        $recordingKeywords = array_map(function (RecordingEntity $recording) {
            return RecordingService::getRecordingKeywords($recording);
        }, iterator_to_array($paginator));

        // Create a recursive iterator so we can merge the recording keywords
        $iterator = new RecursiveIteratorIterator(new RecursiveArrayIterator($recordingKeywords));
        $keywords = iterator_to_array($iterator, false);

        $descArgs = [
            $query,
            $page,
            $paginator->count(),
        ];

        return new HtmlResponse($this->template->render('app::search', [
            'recordings'  => $paginator,
            'query'       => $query,

            // SEO
            'keywords'    => $keywords,
            'description' => vsprintf('Search results for \'%s\', page %d out of %d', $descArgs),
        ]));
    }
}
