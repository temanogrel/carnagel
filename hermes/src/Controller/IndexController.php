<?php
/**
 *
 *
 *
 */

namespace Hermes\Controller;

use Hermes\Repository\UrlRepositoryInterface;
use Hermes\Service\UrlServiceInterface;
use Symfony\Component\HttpFoundation\Response;

class IndexController
{
    /**
     * @var UrlRepositoryInterface
     */
    private $repository;
    /**
     * @var UrlServiceInterface
     */
    private $service;

    /**
     * @param UrlRepositoryInterface $repository
     * @param UrlServiceInterface    $service
     */
    public function __construct(UrlRepositoryInterface $repository, UrlServiceInterface $service)
    {
        $this->service    = $service;
        $this->repository = $repository;
    }

    public function pageNotFound()
    {
        return new Response('Page not found.', 404);
    }
    
    /**
     * Redirect the user if a matching url is found using the hostname and key
     *
     * Will return the following responses
     *  308 - When a match is found
     *  404 - No match is found
     *
     * @param string $key
     * @param string $host
     *
     * @return Response
     */
    public function redirect($key, $host)
    {
        $url = $this->repository->getByKeyAndHostname($key, $host);

        if (!$url) {
            return new Response('Page not found', Response::HTTP_NOT_FOUND);
        }

        $this->service->incrementTransmissions($url);

        $response = new Response('', Response::HTTP_SEE_OTHER);

        if ($url->hasUpstoreDownloadHash()) {
            $response->headers->add(['location' => 'http://downl.ink/' . $url->getUpstoreDownloadHash()]);
        } else {
            $response->headers->add(['location' => $url->getOriginalUrl()]);
        }

        return $response;
    }
}
