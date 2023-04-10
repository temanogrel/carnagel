<?php
/**
 *
 *
 *  AB
 */

declare(strict_types=1);

namespace Aphrodite\Performer\Controller\Performer;

use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Entity\ValueObject\Images;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Recording\Service\RecordingService;
use Aphrodite\Recording\Service\RecordingServiceInterface;
use Aphrodite\Stdlib\Hydrator\Strategy\DateTimeStrategy;
use DateTime;
use Zend\Http\Request;
use Zend\Http\Response;
use Zend\Hydrator\ClassMethods;
use Zend\Hydrator\Exception\BadMethodCallException;
use Zend\Hydrator\Strategy\ClosureStrategy;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class RecordingCollectionController
 *
 * @method Request getRequest
 * @method Response getResponse
 * @method bool isGranted($permission, $context = null)
 */
class RecordingCollectionController extends AbstractRestfulController
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
     * @param RecordingServiceInterface    $recordingService
     * @param PerformerRepositoryInterface $performerRepository
     * @param RecordingRepositoryInterface $recordingRepository
     */
    public function __construct(
        RecordingServiceInterface $recordingService,
        PerformerRepositoryInterface $performerRepository,
        RecordingRepositoryInterface $recordingRepository
    ) {
        $this->recordingService    = $recordingService;
        $this->recordingRepository = $recordingRepository;
        $this->performerRepository = $performerRepository;
    }

    /**
     * Create a new performer
     *
     * @throws NotFoundException
     * @throws BadRequestException
     * @throws UnauthorizedException
     * @throws BadMethodCallException
     *
     * @return ResourceViewModel
     */
    public function post()
    {
        if (!$this->isGranted(RecordingService::PERMISSION_CREATE)) {
            throw new UnauthorizedException;
        }
        try {

            $performer = $this->performerRepository->getById($this->params('performerId'));
            if (!$performer) {
                throw new NotFoundException;
            }

            try {
                $data = Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY);
            } catch (RuntimeException $e) {
                throw new BadRequestException('Invalid json body provided');
            }

            $recording = new RecordingEntity();

            $imageUrlStrategy = new ClosureStrategy(null, function (array $values) {
                return new Images($values['thumb'] ?? null, $values['large'] ?? null, $values['gallery'] ?? null);
            });


            if (isset($data['id'])) {
                $recording->setOldId($data['id']);
            }

            // Migrator is providing us with null instead of an empty array
            $data['images']  = isset($data['images']) && is_array($data['images']) ? $data['images'] : [];
            $data['sprites'] = isset($data['sprites']) && is_array($data['sprites']) ? $data['sprites'] : [];

            // todo: no validation or hydrators yet
            $hydrator = new ClassMethods();
            $hydrator->addStrategy('createdAt', new DateTimeStrategy(DateTime::RFC3339));
            $hydrator->addStrategy('updatedAt', new DateTimeStrategy(DateTime::RFC3339));
            $hydrator->addStrategy('lastCheckedAt', new DateTimeStrategy(DateTime::RFC3339));
            $hydrator->addStrategy('imageUrls', $imageUrlStrategy);
            $hydrator->hydrate($data, $recording);

            $this->recordingService->create($recording, $performer);

            $this
                ->getResponse()
                ->setStatusCode(Response::STATUS_CODE_201);

            return new ResourceViewModel(['recording' => $recording], ['template' => 'recording/resource']);
        } catch (\Throwable $e) {
            echo $e->getMessage() . PHP_EOL;
            echo $e->getTraceAsString();
            exit(1);
        }
    }

    /**
     * Get all the recordings for a performer
     *
     * @throws NotFoundException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        if (!$this->isGranted(RecordingService::PERMISSION_READ)) {
            throw new UnauthorizedException;
        }

        $performer = $this->performerRepository->getById($this->params('performerId'));
        if (!$performer) {
            throw new NotFoundException;
        }

        $recordings = $this->recordingRepository->getForPerformer($performer);

        return new ResourceViewModel(['recordings' => $recordings], ['template' => 'recording/collection']);
    }
}
