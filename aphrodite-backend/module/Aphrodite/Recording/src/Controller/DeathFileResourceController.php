<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Controller;

use Aphrodite\Recording\InputFilter\DeathFileUpdateInputFilter;
use Aphrodite\Recording\DeathFilePermissions;
use Aphrodite\Recording\Repository\DeathFileRepositoryInterface;
use Aphrodite\Recording\Service\DeathFileServiceInterface;
use Zend\Stdlib\Hydrator\ClassMethods;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;
use ZfrRest\View\Model\ResourceViewModel;

/**
 * Class DeathFileCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 *
 * @method boolean isGranted($permission, $context = null)
 */
class DeathFileResourceController extends AbstractRestfulController
{
    /**
     * @var DeathFileServiceInterface
     */
    private $service;

    /**
     * @var DeathFileRepositoryInterface
     */
    private $repository;

    /**
     * DeathFileCollectionController constructor.
     *
     * @param DeathFileServiceInterface    $service
     * @param DeathFileRepositoryInterface $repository
     */
    public function __construct(DeathFileServiceInterface $service, DeathFileRepositoryInterface $repository)
    {
        $this->service    = $service;
        $this->repository = $repository;
    }

    /**
     * Retrieve a death file
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     *
     * @return ResourceViewModel
     */
    public function get()
    {
        $file = $this->repository->getById($this->params('id'));
        if (!$file) {
            throw new NotFoundException('Death file not found');
        }

        if (!$this->isGranted(DeathFilePermissions::VIEW_DEATH_FILE, $file)) {
            throw new UnauthorizedException();
        }

        return new ResourceViewModel(['file' => $file], ['template' => 'death-file/resource']);
    }

    /**
     * Update a death file
     *
     * @throws NotFoundException
     * @throws UnauthorizedException
     *
     * @return ResourceViewModel
     */
    public function put()
    {
        $file = $this->repository->getById($this->params('id'));
        if (!$file) {
            throw new NotFoundException('Death file not found');
        }

        if (!$this->isGranted(DeathFilePermissions::UPDATE_DEATH_FILE, $file)) {
            throw new UnauthorizedException();
        }

        $values = $this->validateIncomingData(DeathFileUpdateInputFilter::class);
        $this->hydrateObject(ClassMethods::class, $file, $values);

        $this->service->update($file);

        return new ResourceViewModel(['file' => $file], ['template' => 'death-file/resource']);
    }

    public function delete()
    {
        $file = $this->repository->getById($this->params('id'));
        if (!$file) {
            throw new NotFoundException('Death file not found');
        }

        if (!$this->isGranted(DeathFilePermissions::DELETE_DEATH_FILE, $file)) {
            throw new UnauthorizedException();
        }

        $this->service->delete($file);

        $response = $this->getResponse();
        $response->setStatusCode(204);

        return $response;
    }
}
