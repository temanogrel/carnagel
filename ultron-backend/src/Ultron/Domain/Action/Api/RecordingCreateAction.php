<?php
/**
 *
 *
 *
 */

declare(strict_types=1);

namespace Ultron\Domain\Action\Api;

use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Ultron\Domain\Exception\PerformerNotFoundException;
use Ultron\Domain\Exception\RecordingMissingAudioException;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\Service\PerformerServiceInterface;
use Ultron\Domain\Service\RecordingServiceInterface;
use Ultron\Infrastructure\Repository\PerformerRepositoryInterface;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Zend\Stdlib\Parameters;

class RecordingCreateAction
{
    /**
     * @var RecordingServiceInterface
     */
    private $recordingService;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @var PerformerServiceInterface
     */
    private $performerService;

    /**
     * RecordingCreateAction constructor.
     *
     * @param RecordingServiceInterface    $recordingService
     * @param PerformerServiceInterface    $performerService
     * @param RecordingRepositoryInterface $recordingRepository
     * @param PerformerRepositoryInterface $performerRepository
     */
    public function __construct(
        RecordingServiceInterface $recordingService,
        PerformerServiceInterface $performerService,
        RecordingRepositoryInterface $recordingRepository,
        PerformerRepositoryInterface $performerRepository
    )
    {
        $this->recordingService = $recordingService;
        $this->performerService = $performerService;
        $this->recordingRepository = $recordingRepository;
        $this->performerRepository = $performerRepository;
    }

    /**
     * @param ServerRequestInterface $request
     * @param ResponseInterface      $response
     * @param callable|null          $next
     *
     * @return ResponseInterface
     */
    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next = null)
    {
        /* @var $data array */
        $data = $request->getParsedBody();

        $performerData = new Parameters($data['performer'] ?? []);
        $recordingData = new Parameters($data['recording'] ?? []);

        try {
            $performer = $this->performerRepository->getByUid($performerData->get('id'));
            $performer->setStageName($performerData->get('stageName'));

            $this->performerService->update($performer);

        } catch (PerformerNotFoundException $e) {
            $performer = $this->performerService->create($performerData);
        }

        try {
            try {
                if ($this->recordingRepository->getByUid($recordingData->get('id'))) {
                    return $response->withStatus(409, 'Recording already exists');
                }
            } catch (RecordingNotFoundException $e) {
                $this->recordingService->create($recordingData, $performer);
            }
        } catch (RecordingMissingAudioException $e) {
            // do nothing
        }

        return $response->withStatus(204);
    }
}
