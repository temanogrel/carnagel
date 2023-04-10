<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Controller;

use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use Aphrodite\Performer\Service\PerformerService;
use Doctrine\Common\Collections\ArrayCollection;
use DoctrineModule\Paginator\Adapter\Collection;
use Zend\Paginator\Paginator;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class BlacklistCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 * @method bool isGranted($permission, $context = null)
 */
class BlacklistCollectionController extends AbstractRestfulController
{
    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @param PerformerRepositoryInterface $performerRepository
     */
    public function __construct(PerformerRepositoryInterface $performerRepository)
    {
        $this->performerRepository = $performerRepository;
    }

    public function get()
    {
        if (!$this->isGranted(PerformerService::PERMISSION_READ)) {
            throw new UnauthorizedException;
        }

        $performers = $this->performerRepository->getBlacklisted();
        $paginator = new Paginator(new Collection(new ArrayCollection($performers)));

        return new ResourceViewModel(['performers' => $paginator], ['template' => 'performer/collection']);
    }
}
