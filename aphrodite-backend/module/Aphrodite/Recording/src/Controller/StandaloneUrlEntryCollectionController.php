<?php

namespace Aphrodite\Recording\Controller;

use Aphrodite\Recording\DeathFilePermissions;
use Aphrodite\Recording\Repository\DeathFile\UrlRepository;
use Doctrine\Common\Collections\Criteria;
use DoctrineModule\Paginator\Adapter\Selectable;
use Zend\Paginator\Paginator;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class StandaloneUrlEntryCollectionController
 *
 * @method boolean isGranted($permission, $context = null)
 */
class StandaloneUrlEntryCollectionController extends AbstractRestfulController
{
    /**
     * @var UrlRepository
     */
    private $repository;

    /**
     * StandaloneUrlEntryCollectionController constructor.
     *
     * @param UrlRepository $repository
     */
    public function __construct(UrlRepository $repository)
    {
        $this->repository = $repository;
    }

    /**
     * List all of url entries
     *
     * @throws UnauthorizedException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        if (!$this->isGranted(DeathFilePermissions::LIST_ENTRIES)) {
            throw new UnauthorizedException;
        }

        $criteria = new Criteria();

        if ($this->params()->fromQuery('state') !== null) {
            $criteria->andWhere($criteria->expr()->eq('state', $this->params()->fromQuery('state')));
        }

        $urls = new Paginator(new Selectable($this->repository, $criteria));
        $urls->setCurrentPageNumber($this->params()->fromQuery('page'));
        $urls->setItemCountPerPage($this->params()->fromQuery('limit'));

        return new ResourceViewModel(['urls' => $urls], ['template' => 'death-file/url/collection']);
    }
}
