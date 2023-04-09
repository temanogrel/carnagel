<?php
/**
 *
 *
 *
 */

namespace Hermes\Controller;

use Hermes\Entity\UrlEntity;
use Hermes\Http\Response\FailedDataValidationResponse;
use Hermes\Http\Response\ResourceNotFoundResponse;
use Hermes\Hydrator\UrlHydrator;
use Hermes\Repository\UrlRepositoryInterface;
use Hermes\Service\UrlServiceInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

class ApiController
{
    /**
     * @var UrlServiceInterface
     */
    private $service;

    /**
     * @var UrlHydrator
     */
    private $hydrator;

    /**
     * @var UrlRepositoryInterface
     */
    private $repository;

    /**
     * @param UrlHydrator            $hydrator
     * @param UrlRepositoryInterface $repository
     * @param UrlServiceInterface    $service
     */
    public function __construct(UrlHydrator $hydrator, UrlRepositoryInterface $repository, UrlServiceInterface $service)
    {
        $this->service    = $service;
        $this->hydrator   = $hydrator;
        $this->repository = $repository;
    }

    /**
     * Create a new short url
     *
     * @param string  $host
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function create(Request $request, $host)
    {
        $rawUrl      = $request->request->get('url');
        $filteredUrl = filter_var($rawUrl, FILTER_VALIDATE_URL);

        if (!$filteredUrl) {
            return new FailedDataValidationResponse([
                'url' => [
                    'invalidUrl' => sprintf('The url %s failed url validation', $rawUrl)
                ]
            ]);
        }

        $entity = $this->service->create($filteredUrl, $host);

        return new JsonResponse($this->hydrator->extract($entity));
    }

    /**
     * Change the url of a existing key
     *
     * @param Request $request
     * @param string  $key
     * @param string  $host
     *
     * @return FailedDataValidationResponse|Response
     */
    public function update(Request $request, $key, $host)
    {
        $rawUrl      = $request->request->get('url');
        $filteredUrl = filter_var($rawUrl, FILTER_VALIDATE_URL);

        if (!$filteredUrl) {
            return new FailedDataValidationResponse([
                'url' => [
                    'invalidUrl' => sprintf('The url %s failed url validation', $rawUrl)
                ]
            ]);
        }

        $entity = $this->repository->getByKeyAndHostname($key, $host);
        if (!$entity) {
            return new ResourceNotFoundResponse();
        }

        $entity->setOriginalUrl($filteredUrl);
        $entity->setUpstoreDownloadHash(null);

        $this->service->update($entity);

        return new Response('', Response::HTTP_NO_CONTENT);
    }

    /**
     * Retrieve a url object
     *
     * @param Request $request
     * @param string  $key
     * @param string  $host
     *
     * @return ResourceNotFoundResponse
     */
    public function get(Request $request, $key, $host)
    {
        switch ($request->query->get('identifier', 'key'))
        {
            case 'id':
                $url = $this->repository->getById($key);
                break;

            case 'originalUrl':
                $url = $this->getByOriginalUrl(base64_decode($key));
                break;

            default:
                $url = $this->repository->getByKeyAndHostname($key, $host);
                break;
        }

        if (!$url) {
            return new ResourceNotFoundResponse();
        }

        return new JsonResponse($this->hydrator->extract($url));
    }

    /**
     * Delete a object
     *
     * @param string $key
     * @param string $host
     *
     * @return ResourceNotFoundResponse|Response
     */
    public function delete($key, $host)
    {
        $url = $this->repository->getByKeyAndHostname($key, $host);
        if (!$url) {
            return new ResourceNotFoundResponse();
        }

        $this->service->delete($url);

        return new Response('', Response::HTTP_NO_CONTENT);
    }

    /**
     * Attempt to retrieve a url by it's original url
     *
     * @param $url
     *
     * @return UrlEntity|null
     */
    private function getByOriginalUrl($url)
    {
        $entity = $this->repository->getByOriginalUrl($url);
        if ($entity) {
            return $entity;
        }

        $matches = [];

        // We failed to get an exacting match, so now we try and extract the key
        if (!preg_match('/^http:\/\/(upsto\.re|upstore\.net)\/(?P<key>[a-zA-Z0-9]+)(\/(.*))?$/', $url, $matches)) {
            return null;
        }

        return $this->repository->getByUpstoreCode($matches['key']);
    }
}
