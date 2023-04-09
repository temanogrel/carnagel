<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;


use Doctrine\ORM\NonUniqueResultException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\SiteConfiguration;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Repository\RecordingRepository;
use Zend\Diactoros\Response\RedirectResponse;
use Zend\Expressive\Helper\UrlHelper;

final class CamgirlGalleryAction
{
    /**
     * @var RecordingRepository
     */
    private $recordingRepository;
    /**
     * @var UrlHelper
     */
    private $urlHelper;
    /**
     * @var Sites
     */
    private $sites;

    /**
     * CamgirlGalleryAction constructor.
     *
     * @param Sites               $sites
     * @param RecordingRepository $recordingRepository
     * @param UrlHelper           $urlHelper
     */
    public function __construct(Sites $sites, RecordingRepository $recordingRepository, UrlHelper $urlHelper)
    {
        $this->sites               = $sites;
        $this->urlHelper           = $urlHelper;
        $this->recordingRepository = $recordingRepository;
    }

    public function __invoke(ServerRequestInterface $request, ResponseInterface $response, callable $next = null)
    {
        $url       = $request->getQueryParams()['url'];
        $performer = $request->getQueryParams()['p'];

        try {
            return $this->redirectToRandomSite($this->recordingRepository->getByGalleryUrl($url));
        } catch (NonUniqueResultException $e) {
            return new RedirectResponse($this->urlHelper->generate('recording.search', ['query' => $performer]));
        } catch (RecordingNotFoundException $e) {
            return new RedirectResponse($this->urlHelper->generate('recording.search', ['query' => $performer]));
        }
    }

    private function redirectToRandomSite(RecordingEntity $recording)
    {
        $sites = $this->sites->getSiteConfigurations()->filter(function (SiteConfiguration $site) use ($recording) {
            if (!$site->isEnabled()) {
                return false;
            }

            if (!$recording->getPerformer()->belongsTo($site)) {
                return false;
            }

            return true;
        });

        /* @var $site SiteConfiguration */
        $site = $sites->getValues()[random_int(0, $sites->count() - 1)];

        $this->urlHelper->setBasePath($site->getDomain());

        $urlArgs = [
            'slug'   => $recording->getSlug(),
            'prefix' => $site->getUrlRoot(),
        ];

        return new RedirectResponse(
            'http:/' . $this->urlHelper->generate('recording.details', $urlArgs)
        );
    }
}
