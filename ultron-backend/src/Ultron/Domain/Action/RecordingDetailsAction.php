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
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\Service\RecordingService;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Zend\Diactoros\Response\HtmlResponse;
use Zend\Expressive\Template\TemplateRendererInterface;

class RecordingDetailsAction
{
    /**
     * @var TemplateRendererInterface
     */
    private $template;

    /**
     * @var RecordingRepositoryInterface
     */
    private $repository;

    /**
     * RecordingDetailsAction constructor.
     *
     * @param TemplateRendererInterface    $template
     * @param RecordingRepositoryInterface $recordingRepository
     */
    public function __construct(TemplateRendererInterface $template, RecordingRepositoryInterface $recordingRepository)
    {
        $this->template   = $template;
        $this->repository = $recordingRepository;
    }

    /**
     * @param ServerRequestInterface $request
     * @param ResponseInterface      $response
     * @param callable               $next
     *
     * @return HtmlResponse
     */
    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        try {
            $recording = $this->repository->getBySlug($request->getAttribute('slug'));
        } catch (RecordingNotFoundException $e) {
            return $response->withStatus(404);
        }

        // Basic stats
        $this->repository->incrementViewCount($recording);
        $keywords = RecordingService::getRecordingKeywords($recording);

        $descArgs = [
            $recording->getPerformer()->getStageName(),
            $recording->getCreatedAt()->format('d/m/y H:i'),
        ];

        return new HtmlResponse($this->template->render('app::details', [
            'recording' => $recording,

            // SEO
            'keywords'    => $keywords,
            'description' => vsprintf('A webcam recording of %s from %s', $descArgs),
        ]));
    }
}
