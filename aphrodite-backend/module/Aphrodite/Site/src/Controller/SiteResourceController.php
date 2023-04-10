<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Controller;

use Aphrodite\Site\Hydrator\SiteHydrator;
use Aphrodite\Site\Repository\SiteRepositoryInterface;
use Aphrodite\Site\Service\SiteService;
use Aphrodite\Site\Service\SiteServiceInterface;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use Zend\Stdlib\Parameters;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class SiteResourceController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 *
 * @method bool isGranted($permission, $context = null)
 */
class SiteResourceController extends AbstractRestfulController
{
    /**
     * @var SiteServiceInterface
     */
    private $service;

    /**
     * @var SiteRepositoryInterface
     */
    private $repository;

    /**
     * @param SiteRepositoryInterface $repository
     * @param SiteServiceInterface    $service
     */
    public function __construct(SiteRepositoryInterface $repository, SiteServiceInterface $service)
    {
        $this->service    = $service;
        $this->repository = $repository;
    }

    public function get()
    {
        $site = $this->repository->getById($this->params('siteId'));
        if (! $site) {
            throw new NotFoundException;
        }

        if (!$this->isGranted(SiteService::PERMISSION_READ, $site)) {
            throw new UnauthorizedException;
        }

        return new ResourceViewModel(['site' => $site], ['template' => 'site/resource']);
    }

    public function delete()
    {
        $site = $this->repository->getById($this->params('siteId'));
        if (! $site) {
            throw new NotFoundException;
        }

        if (!$this->isGranted(SiteService::PERMISSION_DELETE, $site)) {
            throw new UnauthorizedException;
        }

        $this->service->remove($site);

        $response = $this->getResponse();
        $response->setStatusCode(204);

        return $response;
    }

    public function put()
    {
        $site = $this->repository->getById($this->params('siteId'));
        if (! $site) {
            throw new NotFoundException;
        }

        if (!$this->isGranted(SiteService::PERMISSION_UPDATE, $site)) {
            throw new UnauthorizedException;
        }

        try {
            $data = new Parameters(Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY));
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
        }

        $hydrator = new SiteHydrator();
        $hydrator->hydrate($data->toArray(), $site);

        $this->service->update($site);

        return new ResourceViewModel(['site' => $site], ['template' => 'site/resource']);
    }
}
