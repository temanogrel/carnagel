<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Aphrodite\Site\Controller;

use Aphrodite\Site\PostAssociationPermissions;
use Aphrodite\Site\Repository\Exception\PostAssociationNotFoundException;
use Aphrodite\Site\Repository\PostAssociationRepositoryInterface;
use Aphrodite\Site\Service\PostAssociationService;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;

/**
 * Class PostAssociationResourceController
 *
 * @method boolean isGranted(string $permission, $context = null)
 * @method \Zend\Http\Response getResponse
 */
final class PostAssociationResourceController extends AbstractRestfulController
{
    /**
     * @var PostAssociationService
     */
    private $service;

    /**
     * @var PostAssociationRepositoryInterface
     */
    private $repository;

    public function __construct(PostAssociationRepositoryInterface $repository, PostAssociationService $service)
    {
        $this->service    = $service;
        $this->repository = $repository;
    }

    public function delete()
    {
        try {
            $association = $this->repository->getById((int)$this->params('id'));

            if (!$this->isGranted(PostAssociationPermissions::DELETE, $association)) {
                throw new UnauthorizedException();
            }

            $this->service->delete($association);

            $response = $this->getResponse();
            $response->setStatusCode(204);

            return $response;

        } catch (PostAssociationNotFoundException $e) {
            throw new NotFoundException;
        }
    }
}
