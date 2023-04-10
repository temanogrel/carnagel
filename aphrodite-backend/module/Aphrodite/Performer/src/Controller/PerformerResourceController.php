<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Controller;

use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Performer\Service\PerformerService;
use Aphrodite\Performer\Service\PerformerServiceInterface;
use Aphrodite\Stdlib\Hydrator\Strategy\DateTimeStrategy;
use DateTime;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use Zend\Stdlib\Hydrator\ClassMethods;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class PerformerResourceController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 * @method bool isGranted($permission, $context = null)
 */
class PerformerResourceController extends AbstractRestfulController
{
    /**
     * @var PerformerServiceInterface
     */
    private $performerService;

    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @param PerformerRepositoryInterface $performerRepository
     * @param PerformerServiceInterface    $performerService
     */
    public function __construct(
        PerformerRepositoryInterface $performerRepository,
        PerformerServiceInterface $performerService
    ) {
        $this->performerService    = $performerService;
        $this->performerRepository = $performerRepository;
    }

    /**
     * Find a performer
     *
     * @throws NotFoundException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        if (!$this->isGranted(PerformerService::PERMISSION_READ)) {
            throw new UnauthorizedException;
        }

        $splitArgs = function ($arg) {
            $parts = explode(':', $arg);

            if (count($parts) != 2) {
                throw new BadRequestException('Failed to process the identifier');
            }

            return $parts;
        };

        switch ($this->getRequest()->getQuery('identifier')) {
            case 'service-id':
                $performer = $this->performerRepository->getByServiceId(...$splitArgs($this->params('performerId')));
                break;

            default:
                $performer = $this->performerRepository->getById($this->params('performerId'));
                break;
        }

        if (!$performer) {
            throw new NotFoundException;
        }

        return new ResourceViewModel(['performer' => $performer], ['template' => 'performer/resource']);
    }

    /**
     * Find a performer
     *
     * @throws BadRequestException
     * @throws NotFoundException
     *
     * @return ResourceViewModel
     */
    public function put()
    {
        if (!$this->isGranted(PerformerService::PERMISSION_UPDATE)) {
            throw new UnauthorizedException;
        }

        $performer = $this->performerRepository->getById($this->params('performerId'));
        if (!$performer) {
            throw new NotFoundException;
        }

        try {
            $data = Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY);
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
        }

        // todo: no validation or hydrators yet
        $hydrator = new ClassMethods();
        $hydrator->addStrategy('createdAt', new DateTimeStrategy(DateTime::RFC3339));
        $hydrator->addStrategy('updatedAt', new DateTimeStrategy(DateTime::RFC3339));

        $hydrator->hydrate($data, $performer);

        $this->performerService->update($performer);

        return new ResourceViewModel(['performer' => $performer], ['template' => 'performer/resource']);
    }
}
