<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Controller;

use Aphrodite\Site\Entity\Site;
use Aphrodite\Site\Hydrator\SiteHydrator;
use Aphrodite\Site\Repository\SiteRepositoryInterface;
use Aphrodite\Site\Service\SiteService;
use Aphrodite\Site\Service\SiteServiceInterface;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use Zend\Stdlib\Parameters;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class SiteCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 *
 * @method bool isGranted($permission, $context = null)
 */
class SiteCollectionController extends AbstractRestfulController
{
    /**
     * @var SiteRepositoryInterface
     */
    private $repository;

    /**
     * @var SiteServiceInterface
     */
    private $siteService;

    /**
     * @param SiteRepositoryInterface $repository
     * @param SiteServiceInterface             $siteService
     */
    public function __construct(SiteRepositoryInterface $repository, SiteServiceInterface $siteService)
    {
        $this->repository  = $repository;
        $this->siteService = $siteService;
    }

    /**
     * Find all sites matching the given arguments
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        if (!$this->isGranted(SiteService::PERMISSION_READ)) {
            throw new UnauthorizedException;
        }

        $sites = $this->repository->matchingQueryParameters($this->getRequest()->getQuery());

        return new ResourceViewModel(['sites' => $sites], ['template' => 'site/collection']);
    }

    public function post()
    {
        if (!$this->isGranted(SiteService::PERMISSION_CREATE)) {
            throw new UnauthorizedException;
        }

        try {
            $data = new Parameters(Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY));
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
        }

        $site = new Site();

        $hydrator = new SiteHydrator();
        $hydrator->hydrate($data->toArray(), $site);

        $this->siteService->create($site);

        return new ResourceViewModel(['site' => $site], ['template' => 'site/resource']);
    }
}
