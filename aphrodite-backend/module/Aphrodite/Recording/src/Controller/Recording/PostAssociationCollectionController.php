<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Controller\Recording;

use Aphrodite\Recording\Repository\Exception\RecordingNotFoundException;
use Aphrodite\Recording\Repository\RecordingRepositoryInterface;
use Aphrodite\Site\PostAssociationPermissions;
use Aphrodite\Site\Repository\SiteRepositoryInterface;
use Aphrodite\Site\Service\PostAssociationService;
use Aphrodite\Site\Service\PostAssociationServiceInterface;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use Zend\Stdlib\Response;
use ZfcRbac\Exception\UnauthorizedException;
use ZfrRest\Http\Exception\Client\BadRequestException;
use ZfrRest\Http\Exception\Client\NotFoundException;
use ZfrRest\Mvc\Controller\AbstractRestfulController;

/**
 * Class PostAssociationCollectionController
 *
 * @method \Zend\Http\Request getRequest
 * @method \Zend\Http\Response getResponse
 *
 * @method bool isGranted($permission, $context = null)
 */
class PostAssociationCollectionController extends AbstractRestfulController
{
    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var PostAssociationServiceInterface
     */
    private $associationService;

    /**
     * @var SiteRepositoryInterface
     */
    private $siteRepository;

    /**
     * @param RecordingRepositoryInterface    $recordingRepository
     * @param SiteRepositoryInterface         $siteRepository
     * @param PostAssociationServiceInterface $associationService
     */
    public function __construct(
        RecordingRepositoryInterface $recordingRepository,
        SiteRepositoryInterface $siteRepository,
        PostAssociationServiceInterface $associationService
    ) {
        $this->siteRepository      = $siteRepository;
        $this->associationService  = $associationService;
        $this->recordingRepository = $recordingRepository;
    }

    /**
     * Create a new association
     *
     * @throws BadRequestException
     * @throws NotFoundException
     * @throws UnauthorizedException
     */
    public function post()
    {
        try {
            $data = Json::decode($this->getRequest()->getContent(), Json::TYPE_ARRAY);
        } catch (RuntimeException $e) {
            throw new BadRequestException('Invalid json body provided');
        }

        try {
            $recording = $this->recordingRepository->getById($this->params('recordingId'));
        } catch (RecordingNotFoundException $e) {
            throw new NotFoundException('The recording was not found');
        }

        if (!$this->isGranted(PostAssociationPermissions::CREATE, $recording)) {
            throw new UnauthorizedException;
        }

        $site = $this->siteRepository->getById($data['site']);
        if (!$site) {
            throw new NotFoundException(sprintf('The site with id %s was not found', $data['site']));
        }

        $this->associationService->create($recording, $site, $data['post']);

        $response = $this->getResponse();
        $response->setStatusCode(204);

        return $response;
    }
}
