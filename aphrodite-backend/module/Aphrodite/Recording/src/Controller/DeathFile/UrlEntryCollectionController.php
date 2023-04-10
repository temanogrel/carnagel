<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Controller\DeathFile;

use Aphrodite\Recording\DeathFilePermissions;
use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\InputFilter\DeathFile\UrlAddInputFilter;
use Aphrodite\Recording\Repository\DeathFile\UrlRepositoryInterface;
use Aphrodite\Recording\Repository\DeathFileRepositoryInterface;
use Aphrodite\Recording\Service\DeathFile\UrlServiceInterface;
use Doctrine\Common\Collections\Criteria;
use DoctrineModule\Stdlib\Hydrator\DoctrineObject;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class UrlEntryCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 * @method boolean isGranted($permission, $context = null)
 */
class UrlEntryCollectionController extends AbstractRestfulController
{
    /**
     * @var UrlServiceInterface
     */
    private $service;

    /**
     * @var UrlRepositoryInterface
     */
    private $repository;

    /**
     * @var DeathFileRepositoryInterface
     */
    private $deathFileRepository;

    /**
     * UrlEntryCollectionController constructor.
     *
     * @param UrlServiceInterface          $service
     * @param UrlRepositoryInterface       $repository
     * @param DeathFileRepositoryInterface $deathFileRepository
     */
    public function __construct(
        UrlServiceInterface $service,
        UrlRepositoryInterface $repository,
        DeathFileRepositoryInterface $deathFileRepository
    ) {
        $this->service             = $service;
        $this->repository          = $repository;
        $this->deathFileRepository = $deathFileRepository;
    }

    /**
     * Add a new url entry to a death file
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     *
     * @return ResourceViewModel
     */
    public function post()
    {
        $deathFile = $this->deathFileRepository->getById($this->params('id'));
        if (!$deathFile) {
            throw new NotFoundException();
        }

        if (!$this->isGranted(DeathFilePermissions::ADD_ENTRY, $deathFile)) {
            throw new UnauthorizedException();
        }

        $data = $this->validateIncomingData(UrlAddInputFilter::class);

        /* @var $url UrlEntry */
        $url = $this->hydrateObject(DoctrineObject::class, new UrlEntry(), $data);

        $this->service->create($url, $deathFile);

        return new ResourceViewModel(['url' => $url], ['template' => 'death-file/url/resource']);
    }

    /**
     * Retrieve a paginated collection of urls
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        $deathFile = $this->deathFileRepository->getById($this->params('id'));
        if (!$deathFile) {
            throw new NotFoundException();
        }

        if (!$this->isGranted(DeathFilePermissions::LIST_ENTRIES, $deathFile)) {
            throw new UnauthorizedException();
        }

        $query = $this->getRequest()->getQuery();

        $limit  = $query->get('limit');

        $criteria = Criteria::create();
        $criteria->where(Criteria::expr()->eq('deathFile', $deathFile));

        $paginator = $this->repository->paginatedSearch($criteria);
        $paginator->setItemCountPerPage($limit);

        return new ResourceViewModel(['urls' => $paginator], ['template' => 'death-file/url/collection']);
    }
}
