<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Controller;

use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Performer\Service\IntersectionServiceInterface;
use Aphrodite\Performer\Service\PerformerService;
use Doctrine\Common\Collections\Criteria;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use Zend\View\Model\JsonModel;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Http\Exception\Client\MethodNotAllowedException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class PerformerCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 * @method bool isGranted($permission, $context = null)
 */
class PerformerCollectionController extends AbstractRestfulController
{
    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @var IntersectionServiceInterface
     */
    private $intersectionService;

    /**
     * @param PerformerRepositoryInterface  $performerRepository
     * @param IntersectionServiceInterface  $intersectionService
     */
    public function __construct(
        PerformerRepositoryInterface $performerRepository,
        IntersectionServiceInterface $intersectionService
    ) {
        $this->performerRepository  = $performerRepository;
        $this->intersectionService  = $intersectionService;
    }

    /**
     * Intersect the current list of online performers
     *
     * @throws MethodNotAllowedException    Only supports post requests
     * @throws BadRequestException          Malformed json
     *
     * @return ResourceViewModel
     */
    public function intersectAction(): ResourceViewModel
    {
        if (!$this->isGranted(PerformerService::PERMISSION_INTERSECT)) {
            throw new UnauthorizedException;
        }

        if (!$this->getRequest()->isPost()) {
            throw new MethodNotAllowedException('Only post methods are allowed', null, ['POST']);
        }

        // Extract information from the request
        $service = $this->getRequest()->getQuery('service');
        $content = $this->getRequest()->getContent();

        try {
            $data = Json::decode($content, Json::TYPE_ARRAY);
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
        }

        $this
            ->getResponse()
            ->setStatusCode(200);

        // Dispatch
        $performers = $this->intersectionService->process($service, $data);

        return new ResourceViewModel(['performers' => $performers], ['template' => 'performer/collection']);
    }

    /**
     * Get all performers
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        if (!$this->isGranted(PerformerService::PERMISSION_READ)) {
            throw new UnauthorizedException;
        }

        $parameters = $this->getRequest()->getQuery();

        // Criteria internally casts to int if not null
        $limit   = $parameters->get('limit');
        $offset  = $parameters->get('offset');
        $service = $parameters->get('service');
        $indexBy = $parameters->get('index-by');

        $criteria = Criteria::create();
        $criteria->setMaxResults($limit);
        $criteria->setFirstResult($offset);

        if ($parameters->get('sort')) {
            $fields = explode(',', $parameters->get('sort'));

            $orderBy = [];
            foreach ($fields as $field) {
                $parts = explode(':', $field);

                // Ignore it
                if (count($parts) != 2 || in_array($parts[0], ['service'])) continue;

                $orderBy[$parts[0]] = $parts[1];
            }

            $criteria->orderBy($orderBy);
        }

        if ($parameters->get('online') !== null) {
            $criteria->andWhere($criteria->expr()->eq('online', (int)$parameters->get('online')));
        }

        if ($parameters->get('state')) {
            $criteria->andWhere($criteria->expr()->eq('state', $parameters->get('state')));
        }

        if ($parameters->get('stageName')) {
            $criteria->andWhere($criteria->expr()->contains('stageName', $parameters->get('stageName')));
        }

        if ($parameters->get('recording')) {
            $criteria->andWhere($criteria->expr()->eq('isRecording', $parameters->get('recording')));
        }

        if ($parameters->get('since')) {
            $criteria->andWhere(
                $criteria->expr()->gte('updatedAt', date('Y-m-d H:i:s', (int) $parameters->get('since')))
            );
        }

        if ($parameters->get('after')) {
            $criteria->andWhere(
                $criteria->expr()->lte('updatedAt', date('Y-m-d H:i:s', (int) $parameters->get('after')))
            );
        }

        // If no service is provided it will hit all of them
        $performers = $this->performerRepository->search($criteria, $service, $indexBy);
        $performers->setItemCountPerPage($limit);

        return new ResourceViewModel(['performers' => $performers], ['template' => 'performer/collection']);
    }
}
