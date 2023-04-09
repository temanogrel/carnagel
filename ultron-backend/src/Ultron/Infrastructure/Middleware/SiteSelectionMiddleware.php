<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Middleware;

use Doctrine\Common\Collections\Criteria;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Ultron\Domain\Sites;

final class SiteSelectionMiddleware
{
    /**
     * @var Sites
     */
    private $sites;

    /**
     * SiteSelectionMiddleware constructor.
     *
     * @param Sites $sites
     */
    public function __construct(Sites $sites)
    {
        $this->sites = $sites;
    }

    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next)
    {
        $host = $request->getUri()->getHost();

        $criteria = Criteria::create();
        $criteria->andWhere(Criteria::expr()->eq('domain', $host));

        $siteConfigurations = $this->sites->getSiteConfigurations();

        $result = $siteConfigurations->matching($criteria);

        $site = $result->count() === 1 ?
            $result->first() :
            $siteConfigurations->get($this->sites->getDefaultHostname());

        $this->sites->setCurrentSite($site);

        return $next($request, $response);
    }
}
