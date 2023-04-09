<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Middleware;

use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Ultron\Domain\Sites;
use Zend\Crypt\Utils;
use Zend\Expressive\Router\RouteResult;

final class ApiAuthenticationMiddleware
{
    /**
     * @var Sites
     */
    private $sites;

    /**
     * ApiAuthenticationMiddleware constructor.
     *
     * @param Sites $sites
     */
    public function __construct(Sites $sites)
    {
        $this->sites = $sites;
    }

    /**
     * @param ServerRequestInterface $request
     * @param ResponseInterface      $response
     * @param callable               $next
     */
    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        /* @var $routeResult RouteResult */
        $routeResult = $request->getAttribute(RouteResult::class);

        if (
            $routeResult === null ||
            $routeResult->getMatchedRouteName() === false ||
            strpos($routeResult->getMatchedRouteName(), 'api') !== 0
        ) {
            return $next($request, $response);
        }

        $auth = $request->getHeaderLine('Authorization');
        if (Utils::compareStrings($this->sites->getApiAccessToken(), $auth) === false) {
            return $response->withStatus(401);
        }

        return $next($request, $response);
    }
}
