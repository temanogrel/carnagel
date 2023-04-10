<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Controller;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Entity\ValueObject\Images;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Recording\Service\RecordingService;
use Aphrodite\Recording\Service\RecordingServiceInterface;
use DateTime;
use Doctrine\Common\Collections\Criteria;
use Zend\Http\Response;
use Zend\Hydrator\ClassMethods;
use Zend\Hydrator\Strategy\ClosureStrategy;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class RecordingCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 *
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
     * @param RecordingRepositoryInterface $recordingRepository
     * @param RecordingServiceInterface    $recordingService
     */
    public function __construct(
        RecordingRepositoryInterface $recordingRepository,
        RecordingServiceInterface $recordingService
    ) {
        $this->recordingService    = $recordingService;
        $this->recordingRepository = $recordingRepository;
    }

    /**
     * Create a unassociated recording
     */
    public function post()
    {
        if (!$this->isGranted(RecordingService::PERMISSION_CREATE)) {
            throw new UnauthorizedException;
        }

        try {
            $data = Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY);
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
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

        $recording = new RecordingEntity();

        if(isset($data['id'])) {
            $recording->setOldId($data['id']);
        }

        $hydrator->hydrate($data, $recording);

        $this->recordingService->create($recording);

        $this
            ->getResponse()
            ->setStatusCode(Response::STATUS_CODE_201);

        return new ResourceViewModel(['recording' => $recording], ['template' => 'recording/resource']);
    }

    /**
     * Get all recording
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        if (!$this->isGranted(RecordingService::PERMISSION_READ)) {
            throw new UnauthorizedException;
        }

        $parameters = $this->getRequest()->getQuery();

        $limit  = $parameters->get('limit');
        $offset = $parameters->get('offset');

        $criteria = Criteria::create();
        $criteria->setMaxResults($limit);
        $criteria->setFirstResult($offset);

        $expr = $criteria->expr();

        if ($parameters->get('sort')) {
            $fields = explode(',', $parameters->get('sort'));

            $orderBy = [];
            foreach ($fields as $field) {
                $parts = explode(':', $field);

                // Ignore it
                if (count($parts) != 2) {
                    continue;
                }

                $orderBy[$parts[0]] = $parts[1];
            }

            $criteria->orderBy($orderBy);
        }

        if ($parameters->get('performer')) {
            $criteria->andWhere($expr->eq('performer', $parameters->get('performer')));
        }

        if ($parameters->get('stageName')) {
            $criteria->andWhere($expr->contains('stageName', $parameters->get('stageName')));
        }

        if ($parameters->get('service')) {
            $criteria->andWhere($expr->eq('service', $parameters->get('service')));
        }

        if ($parameters->get('state')) {
            $criteria->andWhere($expr->eq('state', $parameters->get('state')));
        }

        if ($parameters->get('since')) {
            $criteria->andWhere(
                $expr->gte('updatedAt', date('Y-m-d H:i:s', (int)$parameters->get('since')))
            );
        }

        if ($parameters->get('after')) {
            $criteria->andWhere(
                $criteria->expr()->lte('updatedAt', date('Y-m-d H:i:s', (int) $parameters->get('after')))
            );
        }

        if ($parameters->get('checkedAt')) {
            $criteria->andWhere(
                $expr->orX(
                    $expr->lte('lastCheckedAt', date('Y-m-d H:i:s', (int) $parameters->get('checkedAt'))),
                    $expr->isNull('lastCheckedAt')
                )
            );
        }

        if ($parameters->get('orphaned') !== null) {
            $criteria->andWhere(
                $expr->eq('orphaned', (bool) $parameters->get('orphaned'))
            );
        }

        $recordings = $this->recordingRepository->matching($criteria);
        $recordings->setItemCountPerPage($limit);

        return new ResourceViewModel(['recordings' => $recordings], ['template' => 'recording/collection']);
    }
}
