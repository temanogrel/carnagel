<?php
/**
 *
 *
 *  AB
 */

declare(strict_types=1);

namespace Aphrodite\Recording\Controller;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Entity\ValueObject\Images;
use Aphrodite\Recording\Repository\Exception\NonUniqueRecordingResultException;
use Aphrodite\Recording\Repository\Exception\RecordingNotFoundException;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Recording\Service\RecordingService;
use Aphrodite\Recording\Service\RecordingServiceInterface;
use Aphrodite\Stdlib\Hydrator\Strategy\DateTimeStrategy;
use DateTime;
use Zend\Http\Request;
use Zend\Http\Response;
use Zend\Hydrator\ClassMethods;
use Zend\Hydrator\Strategy\ClosureStrategy;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Http\Exception\Client\ConflictException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class RecordingResourceController
 *
 * @method Request getRequest
 * @method Response getResponse
 *
 * @method bool isGranted($permission, $context = null)
 */
class RecordingResourceController extends AbstractRestfulController
{
    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var RecordingServiceInterface
     */
    private $recordingService;

    /**
     * @param RecordingRepositoryInterface $recordingRepository
     * @param RecordingServiceInterface    $recordingService
     */
    public function __construct(
        RecordingRepositoryInterface $recordingRepository,
        RecordingServiceInterface $recordingService
    ) {

        $this->recordingRepository = $recordingRepository;
        $this->recordingService    = $recordingService;
    }

    /**
     * Retrieve a recording
     *
     * @throws BadRequestException
     * @throws UnauthorizedException
     * @throws NotFoundException
     * @throws ConflictException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        $recording = $this->getRecording();
        if (!$this->isGranted(RecordingService::PERMISSION_READ, $recording)) {
            throw new UnauthorizedException;
        }

        return new ResourceViewModel(['recording' => $recording], ['template' => 'recording/resource']);
    }

    /**
     * Update a recording
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     * @throws ConflictException
     * @throws BadRequestException
     *
     * @return ResourceViewModel
     */
    public function put()
    {
        $recording = $this->getRecording();
        if (!$this->isGranted(RecordingService::PERMISSION_UPDATE, $recording)) {
            throw new UnauthorizedException;
        }

        try {
            $data = Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY);
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
        }

        // Don't handle this part.
        if (isset($data['publishedOn'])) {
            unset($data['publishedOn']);
        }

        $imageUrlStrategy = new ClosureStrategy(null, function (array $values) {
            return new Images($values['thumb'] ?? null, $values['large'] ?? null, $values['gallery'] ?? null);
        });

        // todo: no validation or hydrators yet
        $hydrator = new ClassMethods();
        $hydrator->addStrategy('createdAt', new DateTimeStrategy(DateTime::RFC3339));
        $hydrator->addStrategy('updatedAt', new DateTimeStrategy(DateTime::RFC3339));
        $hydrator->addStrategy('lastCheckedAt', new DateTimeStrategy(DateTime::RFC3339));
        $hydrator->addStrategy('imageUrls', $imageUrlStrategy);
        $hydrator->hydrate($data, $recording);

        $this->recordingService->update($recording);

        return new ResourceViewModel(['recording' => $recording], ['template' => 'recording/resource']);
    }

    /**
     * Delete a recording
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     * @throws ConflictException
     * @throws BadRequestException
     *
     * @return Response
     */
    public function delete()
    {
        $recording = $this->getRecording();
        if (!$this->isGranted(RecordingService::PERMISSION_DELETE, $recording)) {
            throw new UnauthorizedException;
        }

        $this->recordingService->delete($recording);

        $response = $this->getResponse();
        $response->setStatusCode(204);

        return $response;
    }

    /**
     * @throws BadRequestException
     * @throws ConflictException
     * @throws NotFoundException
     *
     * @return RecordingEntity
     */
    private function getRecording(): RecordingEntity
    {
        $splitArgs = function ($arg) {
            $parts = explode(':', base64_decode($arg));

            if (count($parts) !== 2) {
                throw new BadRequestException('Failed to process the identifier');
            }

            return $parts;
        };

        $identifier = $this->params('recordingId');
        $identifyBy = $this->getRequest()->getQuery('identifier');

        try {

            switch ($identifyBy) {
                case 'old-id':
                    return $this->recordingRepository->getByOldId((int) $identifier);

                case 'gallery-url':
                    return $this->recordingRepository->getByGalleryUrl(base64_decode($identifier));

                case 'video-url':
                    return $this->recordingRepository->getByVideoUrl(base64_decode($identifier));

                default:
                    return $this->recordingRepository->getById((int) $identifier);
            }

        } catch (RecordingNotFoundException $e) {
            throw new NotFoundException('The recording was not found');
        } catch (NonUniqueRecordingResultException $e) {
            throw new ConflictException('Multiple recordings where matched', ['identifiers' => $e->getIdentifiers()]);
        }
    }
}
