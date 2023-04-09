<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Repository;

use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\EntityRepository;
use Doctrine\ORM\Mapping;
use Doctrine\ORM\NonUniqueResultException;
use Doctrine\ORM\Query\Expr\Join;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Domain\Service\CacheService;
use Ultron\Domain\Service\Exception\CacheEntryMissingException;
use Ultron\Domain\SiteConfiguration;
use Ultron\Infrastructure\Service\ValueObject\PageInformation;
use Zend\Paginator\Adapter\ArrayAdapter;
use Zend\Paginator\Adapter\Callback;
use Zend\Paginator\Paginator;

final class RecordingRepository extends EntityRepository implements RecordingRepositoryInterface, CacheServiceAwareInterface
{
    /**
     * @var CacheService
     */
    private $cacheService;

    /**
     * @inheritDoc
     */
    public function __construct(EntityManager $em, Mapping\ClassMetadata $class, CacheService $cacheService)
    {
        parent::__construct($em, $class);

        $this->cacheService = $cacheService;
    }

    /**
     * {@inheritdoc}
     */
    public function getById($id):RecordingEntity
    {
        $recording = $this->findOneBy(['id' => $id]);
        if (!$recording) {
            throw new RecordingNotFoundException;
        }

        return $recording;
    }

    /**
     * {@inheritdoc}
     */
    public function getByUid($uid):RecordingEntity
    {
        $recording = $this->findOneBy(['uid' => $uid]);
        if (!$recording) {
            throw new RecordingNotFoundException;
        }

        return $recording;
    }

    /**
     * {@inheritdoc}
     */
    public function getBySlug($slug):RecordingEntity
    {
        $recording = $this->findOneBy(['slug' => $slug]);
        if (!$recording) {
            throw new RecordingNotFoundException;
        }

        return $recording;
    }

    /**
     * {@inheritdoc}
     */
    public function getByGalleryUrl($url):RecordingEntity
    {
        $builder = $this->createQueryBuilder('recording');
        $builder
            ->addSelect('performer')
            ->innerJoin('recording.performer', 'performer')
            ->where('recording.imageUrls.galleryUrl = :url')
            ->setParameter('url', $url);

        try {

            $recording = $builder->getQuery()->getOneOrNullResult();
            if (!$recording) {
                throw new RecordingNotFoundException();
            }

            return $recording;

        } catch (NonUniqueResultException $e) {
            throw new RecordingNotFoundException('', 0, $e);
        }
    }

    /**
     * {@inheritdoc}
     */
    public function incrementViewCount(RecordingEntity $recording)
    {
        $builder = $this->getEntityManager()->createQueryBuilder();
        $builder
            ->update(RecordingEntity::class, 'r')
            ->set('r.views', 'r.views + 1')
            ->where('r.id = :id')
            ->setParameter('id', $recording->getId());

        $builder->getQuery()->execute();
    }

    /**
     * {@inheritdoc}
     */
    public function getTotalCount(SiteConfiguration $site = null): int
    {
        $builder = $this->createQueryBuilder('recording');
        $builder->select($builder->expr()->count('recording.id'));
        $builder->innerJoin('recording.performer', 'performer');

        if ($site && $site->getService()) {
            $builder
                ->andWhere('performer.service = :service')
                ->setParameter('service', $site->getService());
        }

        if ($site && $site->getSection()) {
            $builder
                ->andWhere('performer.section = :section')
                ->setParameter('section', $site->getSection());
        }

        return (int) $builder->getQuery()->getSingleScalarResult();
    }

    /**
     * {@inheritdoc}
     */
    public function getPageInformation(
        SiteConfiguration $site,
        int $limit,
        int $minId,
        int $maxId = null
    ): PageInformation {
        $builder = $this->createQueryBuilder('recording');

        $maxIdJoinQuery = $maxId ? ' AND recording.id <= :maxId' : '';

        $builder->select('recording.id');
        $builder->innerJoin(
            'recording.performer',
            'performer',
            Join::WITH,
            'recording.performer = performer.id AND recording.id >= :minId' . $maxIdJoinQuery
        );

        if ($site->getService()) {
            $builder
                ->andWhere('performer.service = :service')
                ->setParameter('service', $site->getService());
        }

        if ($site->getSection()) {
            $builder
                ->andWhere('performer.section = :section')
                ->setParameter('section', $site->getSection());
        }

        $builder->setParameter('minId', $minId);

        if ($maxId) {
            $builder->setParameter('maxId', $maxId);
        }

        $builder->setMaxResults($limit);
        $builder->orderBy('recording.id', 'ASC');

        $result = $builder->getQuery()->getArrayResult();

        return new PageInformation((int) end($result)['id'], count($result));
    }

    /**
     * {@inheritdoc}
     */
    public function searchByPerformer(string $stageName, Criteria $criteria = null, SiteConfiguration $site = null):Paginator
    {
        $escapedStageName = addcslashes($stageName, '%_.^$[](){}');

        $builder    = $this->createQueryBuilder('recording');
        $expression = $builder->expr();

        $builder->addCriteria($criteria);
        $builder->addSelect('performer');
        $builder->innerJoin('recording.performer', 'performer');

        $builder
            ->andWhere($expression->like('recording.stageName', ':stageName'))
            ->setParameter('stageName', $escapedStageName . '%');

        if ($site && $site->getService()) {
            $builder
                ->andWhere('performer.service = :service')
                ->setParameter('service', $site->getService());
        }

        if ($site && $site->getSection()) {
            $builder
                ->andWhere('performer.section = :section')
                ->setParameter('section', $site->getSection());
        }

        $dataCallback = function ($offset, $limit) use ($builder) {
            $builder->setFirstResult($offset);
            $builder->setMaxResults($limit);

            return $builder->getQuery()->getResult();
        };

        $countCallback = function () use ($builder) {
            $count = clone $builder;
            $count->select('COUNT(recording.id)');
            $count->resetDQLPart('orderBy');
            $count->setFirstResult(null);
            $count->setMaxResults(null);

            return $count->getQuery()->getSingleScalarResult();
        };

        return new Paginator(new Callback($dataCallback, $countCallback));
    }

    /**
     * {@inheritdoc}
     */
    public function getPaginatedResult(Criteria $criteria, SiteConfiguration $site):Paginator
    {
        $builder = $this->createQueryBuilder('recording');
        $builder
            ->addSelect('performer')
            ->addCriteria($criteria)
            ->innerJoin('recording.performer', 'performer');

        if ($site->getService()) {
            $builder
                ->andWhere('performer.service = :service')
                ->setParameter('service', $site->getService());
        }

        if ($site->getSection()) {
            $builder
                ->andWhere('performer.section = :section')
                ->setParameter('section', $site->getSection());
        }

        $dataCallback = function ($offset, $limit) use ($builder): array {
            $builder->setFirstResult($offset);
            $builder->setMaxResults($limit);

            return $builder->getQuery()->getResult();
        };

        $countCallback = function () use ($builder, $site): int {
            $countBuilder = clone $builder;
            $countBuilder->select('COUNT(recording.id)');
            $countBuilder->resetDQLPart('orderBy');
            $countBuilder->setMaxResults(null);
            $countBuilder->setFirstResult(null);

            return (int) $countBuilder->getQuery()->getSingleScalarResult();
        };

        return new Paginator(new Callback($dataCallback, $countCallback));
    }

    /**
     * {@inheritdoc}
     */
    public function getBetweenIds(Criteria $criteria, SiteConfiguration $site, int $min, int $max): array
    {
        $builder = $this->createQueryBuilder('recording');
        $builder
            ->addSelect('performer')
            ->addCriteria($criteria)
            ->innerJoin('recording.performer', 'performer');

        if ($site->getService()) {
            $builder
                ->andWhere('performer.service = :service')
                ->setParameter('service', $site->getService());
        }

        if ($site->getSection()) {
            $builder
                ->andWhere('performer.section = :section')
                ->setParameter('section', $site->getSection());
        }

        $builder
            ->andWhere($builder->expr()->gte('recording.id', ':min'))
            ->andWhere($builder->expr()->lte('recording.id', ':max'))
            ->setParameter('min', $min)
            ->setParameter('max', $max);

        return $builder->getQuery()->getResult();
    }

    /**
     * {@inheritdoc}
     */
    public function getBySiteAndCriteria(Criteria $criteria, SiteConfiguration $site): array
    {
        $builder = $this->createQueryBuilder('recording');
        $builder
            ->addSelect('performer')
            ->addCriteria($criteria)
            ->innerJoin('recording.performer', 'performer');

        if ($site->getService()) {
            $builder
                ->andWhere('performer.service = :service')
                ->setParameter('service', $site->getService());
        }

        if ($site->getSection()) {
            $builder
                ->andWhere('performer.section = :section')
                ->setParameter('section', $site->getSection());
        }

        return $builder->getQuery()->getResult();
    }

    /**
     * {@inheritdoc}
     */
    public function getRecordingCountOfPerformers(): array
    {
        $builder = $this->createQueryBuilder('recording');
        $builder
            ->select('performer.id as performerId, COUNT(recording.id) as recordingCount')
            ->innerJoin('recording.performer', 'performer')
            ->groupBy('performer.id');

        return $builder->getQuery()->getArrayResult();
    }
}
