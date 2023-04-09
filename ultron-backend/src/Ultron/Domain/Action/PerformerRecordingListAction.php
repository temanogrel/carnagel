<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;

use Doctrine\Common\Collections\Criteria;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use RecursiveArrayIterator;
use RecursiveIteratorIterator;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\Service\RecordingService;
use Ultron\Domain\SiteConfiguration;
use Ultron\Infrastructure\Repository\PerformerRepositoryInterface;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Zend\Diactoros\Response\HtmlResponse;
use Zend\Expressive\Template\TemplateRendererInterface;

class PerformerRecordingListAction
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
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * PerformerRecordingListAction constructor.
     *
     * @param SiteConfiguration $site
     * @param TemplateRendererInterface $template
     * @param PerformerRepositoryInterface $performerRepository
     * @param RecordingRepositoryInterface $recordingRepository
     */
    public function __construct(
        SiteConfiguration $site,
        TemplateRendererInterface $template,
        PerformerRepositoryInterface $performerRepository,
        RecordingRepositoryInterface $recordingRepository
    ) {
        $this->site                = $site;
        $this->template            = $template;
        $this->performerRepository = $performerRepository;
        $this->recordingRepository = $recordingRepository;
    }

    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        try {
            $performer = $this->performerRepository->getBySlug($request->getAttribute('slug'));
        } catch (RecordingNotFoundException $e) {
            return $response->withStatus(404);
        }

        $page = $request->getQueryParams()['page'] ?? 1;

        $criteria = Criteria::create();
        $criteria->orderBy(['createdAt' => 'desc']);
        $criteria->where(Criteria::expr()->eq('performer', $performer));

        $paginator = $this->recordingRepository->getPaginatedResult($criteria, $this->site);
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
            $performer->getStageName(),
            $performer->getService(true),
        ];

        return new HtmlResponse($this->template->render('app::list', [
            'recordings'  => $paginator,
            'hide_more'   => true,

            // SEO
            'keywords'    => $keywords,
            'description' => vsprintf('Page %d of %d for all the recordings available for the %s on %s', $descArgs),
        ]));
    }
}
