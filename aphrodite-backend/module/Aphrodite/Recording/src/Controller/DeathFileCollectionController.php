<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Controller;

use Aphrodite\Recording\InputFilter\DeathFileUploadInputFilter;
use Aphrodite\Recording\DeathFilePermissions;
use Aphrodite\Recording\Repository\DeathFileRepositoryInterface;
use Aphrodite\Recording\Service\DeathFileServiceInterface;
use Doctrine\Common\Collections\Criteria;
use Zend\Stdlib\Parameters;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\MethodNotAllowedException;
use ZfrRest\Http\Exception\Client\UnprocessableEntityException;
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
class DeathFileCollectionController extends AbstractRestfulController
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
     * @var DeathFileUploadInputFilter
     */
    private $uploadInputFilter;

    /**
     * DeathFileCollectionController constructor.
     *
     * @param DeathFileServiceInterface    $service
     * @param DeathFileRepositoryInterface $repository
     * @param DeathFileUploadInputFilter   $uploadInputFilter
     */
    public function __construct(
        DeathFileServiceInterface $service,
        DeathFileRepositoryInterface $repository,
        DeathFileUploadInputFilter $uploadInputFilter
    ) {
        $this->service           = $service;
        $this->repository        = $repository;
        $this->uploadInputFilter = $uploadInputFilter;
    }

    /**
     * Upload a new death list file
     */
    public function post()
    {
        if (!$this->isGranted(DeathFilePermissions::UPLOAD_DEATH_FILE)) {
            throw new UnauthorizedException();
        }

        $request = $this->getRequest();
        if (!$request->isPost()) {
            throw new MethodNotAllowedException('', null, ['POST']);
        }

        $data = array_merge(
            $request->getPost()->toArray(),
            $request->getFiles()->toArray()
        );

        $this->uploadInputFilter->setData($data);
        if (!$this->uploadInputFilter->isValid()) {
            throw new UnprocessableEntityException('Validation error', $this->uploadInputFilter->getMessages());
        }

        $file = new Parameters($this->uploadInputFilter->getValue('file'));

        $entity = $this->service->createFromFile($file);

        return new ResourceViewModel(['file' => $entity], ['template' => 'death-file/resource']);
    }

    public function get()
    {
        if (!$this->isGranted(DeathFilePermissions::LIST_DEATH_FILES)) {
            throw new UnauthorizedException();
        }

        $criteria = new Criteria();
        $criteria->orderBy(['id' => 'desc']);

        $files = $this->repository->matching($criteria);

        return new ResourceViewModel(['files' => $files], ['template' => 'death-file/collection']);
    }
}
