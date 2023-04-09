<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action\Api;


use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\Service\RecordingService;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;

class RecordingDeleteAction
{
    /**
     * @var RecordingService
     */
    private $service;

    /**
     * @var RecordingRepositoryInterface
     */
    private $repository;

    /**
     * RecordingDeleteAction constructor.
     *
     * @param RecordingService             $service
     * @param RecordingRepositoryInterface $repository
     */
    public function __construct(RecordingService $service, RecordingRepositoryInterface $repository)
    {
        $this->service    = $service;
        $this->repository = $repository;
    }

    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        try {
            $recording = $this->repository->getById($request->getAttribute('id'));

            $this->service->remove($recording);

            return $response->withStatus(204);

        } catch (RecordingNotFoundException $e) {
            return $response->withStatus(404);
        }
    }
}
