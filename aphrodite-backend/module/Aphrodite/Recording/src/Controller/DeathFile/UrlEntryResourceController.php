<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Controller\DeathFile;

use Aphrodite\Recording\DeathFilePermissions;
use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\InputFilter\DeathFile\UrlUpdateInputFilter;
use Aphrodite\Recording\Repository\DeathFile\UrlRepositoryInterface;
use Aphrodite\Recording\Service\DeathFile\UrlServiceInterface;
use Zend\Http\Response;
use Zend\Stdlib\Hydrator\ClassMethods;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Http\Exception\Client\UnprocessableEntityException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class UrlEntryResourceController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 *
 * @method boolean isGranted($permission, $context = null)
 */
class UrlEntryResourceController extends AbstractRestfulController
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
     * UrlEntryResourceController constructor.
     *
     * @param UrlServiceInterface    $service
     * @param UrlRepositoryInterface $repository
     */
    public function __construct(UrlServiceInterface $service, UrlRepositoryInterface $repository)
    {
        $this->service    = $service;
        $this->repository = $repository;
    }

    /**
     * Retrieve a url
     *
     * @throws NotFoundException
     *
     * @return UrlEntry|null
     */
    private function getUrl()
    {
        $id         = $this->params()->fromRoute('id');
        $identifier = $this->params()->fromQuery('identifier', 'id');

        switch ($identifier) {
            case 'url':
                $file = $this->repository->getByUrl(base64_decode($id));
                break;

            default:
                $file = $this->repository->getById($id);
                break;
        }

        if (!$file) {
            throw new NotFoundException('Url entry not found');
        }

        return $file;
    }

    /**
     * View a death file url entry
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        $url = $this->getUrl();
        if (!$this->isGranted(DeathFilePermissions::VIEW_ENTRY, $url)) {
            throw new UnauthorizedException();
        }

        return new ResourceViewModel(['url' => $url], ['template' => 'death-file/url/resource']);
    }

    /**
     * Update a url entry
     *
     * {@see UrlUpdateInputFilter} for the available fields to update
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     * @throws UnprocessableEntityException
     *
     * @return ResourceViewModel
     */
    public function put()
    {
        $url = $this->getUrl();
        if (!$this->isGranted(DeathFilePermissions::UPDATE_ENTRY, $url)) {
            throw new UnauthorizedException();
        }

        $data = $this->validateIncomingData(UrlUpdateInputFilter::class);
        $this->hydrateObject(ClassMethods::class, $url, $data);

        $this->service->update($url);

        return new ResourceViewModel(['url' => $url], ['template' => 'death-file/url/resource']);
    }

    /**
     * Remove a recording
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     *
     * @return Response
     */
    public function remove()
    {
        $url = $this->getUrl();
        if (!$this->isGranted(DeathFilePermissions::REMOVE_ENTRY)) {
            throw new UnauthorizedException();
        }

        $this->service->remove($url);

        $response = $this->getResponse();
        $response->setStatusCode(204);

        return $response;
    }
}
