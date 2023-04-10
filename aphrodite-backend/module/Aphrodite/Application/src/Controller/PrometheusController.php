<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Application\Controller;

use A;
use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Prometheus\CollectorRegistry;
use Prometheus\RenderTextFormat;
use Prometheus\Storage\InMemory;
use Zend\Http\Response;
use Zend\Mvc\Controller\AbstractActionController;

/**
 * Class PrometheusController
 *
 * @method Response getResponse
 */
class PrometheusController extends AbstractActionController
{
    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;
    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * PrometheusController constructor.
     *
     * @param RecordingRepositoryInterface $recordingRepository
     * @param PerformerRepositoryInterface $performerRepository
     */
    public function __construct(
        RecordingRepositoryInterface $recordingRepository,
        PerformerRepositoryInterface $performerRepository
    ) {
        $this->recordingRepository = $recordingRepository;
        $this->performerRepository = $performerRepository;
    }

    public function metricsAction(): Response
    {
        $prom = new CollectorRegistry(new InMemory());

        $recordingStateGauge = $prom->getOrRegisterGauge(
            'aphrodite',
            'recording_state',
            'Number of recordings in various states',
            ['state', 'service']
        );

        foreach ($this->recordingRepository->getRecordingsPerState() as $count => list($state, $service)) {
            $recordingStateGauge->set($count, [$state, $service]);
        }

        $performerStateGauge = $prom->getOrRegisterGauge(
            'aphrodite',
            'performer_recording_state',
            'Number of performers in various states',
            ['state', 'service']
        );

        foreach ($this->performerRepository->getPerformerStats() as $count => list($state, $service)) {
            $performerStateGauge->set($count, [$state, $service]);
        }

        $response = $this->getResponse();
        $response->getHeaders()->addHeaderLine('Content-Type', RenderTextFormat::MIME_TYPE);
        $response->setContent((new RenderTextFormat())->render($prom->getMetricFamilySamples()));

        return $response;
    }
}
