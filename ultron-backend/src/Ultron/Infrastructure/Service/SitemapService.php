<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service;

use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityManager;
use Generator;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\SiteConfiguration;
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Repository\PerformerRepositoryInterface;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Ultron\Infrastructure\Service\Sitemap\PerformerWriter;
use Ultron\Infrastructure\Service\Sitemap\RecordingWriter;
use Zend\Expressive\Helper\UrlHelper;

final class SitemapService implements SitemapServiceInterface
{
    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var Sites
     */
    private $sites;

    /**
     * @var UrlHelper
     */
    private $urlHelper;

    /**
     * @var EntityManager
     */
    private $entityManager;

    /**
     * SitemapService constructor.
     * @param PerformerRepositoryInterface $performerRepository
     * @param RecordingRepositoryInterface $recordingRepository
     * @param Sites $sites
     * @param UrlHelper $urlHelper
     * @param EntityManager $entityManager
     */
    public function __construct(
        PerformerRepositoryInterface $performerRepository,
        RecordingRepositoryInterface $recordingRepository,
        Sites $sites,
        UrlHelper $urlHelper,
        EntityManager $entityManager
    ) {
        $this->performerRepository = $performerRepository;
        $this->recordingRepository = $recordingRepository;
        $this->sites               = $sites;
        $this->urlHelper           = $urlHelper;
        $this->entityManager       = $entityManager;
    }

    /**
     * {@inheritdoc}
     */
    public function getSitemapUrls(SiteConfiguration $site): Generator
    {
        $hostname = $site->getDomain();

        foreach (glob(sprintf('public/sitemap/%s-*', $hostname)) as $url) {
            yield 'http://' . $hostname . substr($url, 6);
        }
    }

    /**
     * {@inheritdoc}
     */
    public function createPerformerSitemaps()
    {
        $criteria = Criteria::create();
        $criteria->orderBy(['id' => Criteria::ASC]);
        $criteria->andWhere(Criteria::expr()->gte('recordingCount', 1));
        $criteria->setMaxResults(10000);

        $writers = [];

        foreach ($this->sites->getSiteConfigurations() as $site) {
            if (!$site->isEnabled()) {
                continue;
            }
            
            $lastSeenId = 1;

            $writer    = new PerformerWriter($site, $this->urlHelper);
            $writers[] = $writer;

            do {
                $c = clone $criteria;
                $c->andWhere(Criteria::expr()->gte('id', $lastSeenId));

                /* @var RecordingEntity[] $items */
                $items = $this->performerRepository->getByCriteria($c);
                $this->addToWriter($items, $writer);

                if (array_key_exists(9999, $items)) {
                    $lastSeenId = $items[9999]->getId() + 1;
                }

                $this->entityManager->clear();
            } while (count($items) === 10000);
        }

        /* @var PerformerWriter $writer */
        foreach ($writers as $writer) {
            $writer->flush();
        }
    }

    /**
     * {@inheritdoc}
     */
    public function createRecordingSitemaps()
    {
        $writers = [];

        $rowCount = $this->recordingRepository->getTotalCount();

        /* @var SiteConfiguration $site */
        foreach ($this->sites->getSiteConfigurations() as $site) {
            if (!$site->isEnabled()) {
                continue;
            }

            $lastSeenId = 1;

            $writer    = new RecordingWriter($site, $this->urlHelper);
            $writers[] = $writer;

            do {
                $offset = 50000;

                do {
                    $pageInfo = $this
                        ->recordingRepository
                        ->getPageInformation($site, 10000, $lastSeenId, $lastSeenId + $offset);

                    $offset *= 2;
                } while ($pageInfo->getCount() < $site->getPageSize() && $pageInfo->getMaxId() < $rowCount);

                /* @var RecordingEntity[] $items */
                $items = $this
                    ->recordingRepository
                    ->getBetweenIds(Criteria::create(), $site, $lastSeenId, $pageInfo->getMaxId());
                $this->addToWriter($items, $writer);

                $lastSeenId = $pageInfo->getMaxId() + 1;

                $this->entityManager->clear();
            } while (count($items) === 10000);
        }

        /* @var RecordingWriter $writer */
        foreach ($writers as $writer) {
            $writer->flush();
        }
    }

    /**
     * @param array $items
     * @param $writer
     */
    private function addToWriter(array $items, $writer)
    {
        foreach ($items as $item) {
            $writer->add($item);
        }
    }
}
