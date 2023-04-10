<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Repository;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Repository\Exception\NonUniqueRecordingResultException;
use Aphrodite\Recording\Repository\Exception\RecordingNotFoundException;
use DateTime;
use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityRepository;
use DomainException;
use Generator;
use Zend\Paginator\Adapter\Callback;
use Zend\Paginator\Paginator;

class RecordingRepository extends EntityRepository implements RecordingRepositoryInterface
{
    private function getBy(array $properties): RecordingEntity
    {
        $builder = $this->createQueryBuilder('recording');
        $builder
            ->addSelect('publishedOn')
            ->leftJoin('recording.publishedOn', 'publishedOn')
            ->orderBy('recording.id', 'desc');

        foreach ($properties as $prop => $value) {
            $normalizedProp = str_replace('.', '',$prop);
            
            $builder
                ->andWhere($builder->expr()->eq($prop, ':' . $normalizedProp))
                ->setParameter($normalizedProp, $value);
        }

        $result = $builder->getQuery()->getResult();

        if (count($result) === 0) {
            throw new RecordingNotFoundException();
        } else if (count($result) > 1) {

            $identifiers = array_map(function (RecordingEntity $recording) {
                return $recording->getId();
            }, $result);

            throw new NonUniqueRecordingResultException($identifiers);
        }

        return current($result);
    }

    /**
     * {@inheritdoc}
     */
    public function getById(int $id): RecordingEntity
    {
        return $this->getBy(['recording.id' => $id]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByOldId(int $id): RecordingEntity
    {
        $recording = $this->findOneBy(['oldId' => $id]);
        if (!$recording) {
            throw new RecordingNotFoundException();
        }

        return $recording;
    }

    /**
     * {@inheritdoc}
     */
    public function getByGalleryUrl(string $url): RecordingEntity
    {
        return $this->getBy(['recording.imageUrls.gallery' => $url]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByVideoUrl(string $url): RecordingEntity
    {
        return $this->getBy(['recording.videoUrl' => $url]);
    }

    public function getRecordingsPerState(): Generator
    {
        $conn = $this->getEntityManager()->getConnection();
        $conn->exec('SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED');

        $stmt = $conn->executeQuery('SELECT state, count(id) as c, service FROM recordings GROUP BY state, service');

        foreach ($stmt->fetchAll() as $row) {
           yield $row['c'] => [$row['state'], $row['service']];
        }

        // NOOP, but will clear the transaction isolation level set above
        $conn->exec('COMMIT');
    }

    /**
     * {@inheritdoc}
     */
    public function getForPerformer(AbstractPerformerEntity $performer): Paginator
    {
        $builder = $this->getEntityManager()->createQueryBuilder();
        $builder
            ->select('recording')
            ->from(RecordingEntity::class, 'recording')
            ->where('recording.performer = :performer')
            ->setParameter('performer', $performer);

        $countBuilder = clone $builder;
        $countBuilder
            ->select('COUNT(recording.id)')
            ->resetDQLPart('orderBy')
            ->setMaxResults(null)
            ->setFirstResult(null);

        $countQuery = $countBuilder->getQuery();
        $countQuery->useResultCache(true, 600);

        $count = function () use ($countQuery) {
            return $countQuery->getSingleScalarResult();
        };

        $data = function () use ($builder) {
            return $builder->getQuery()->getResult();
        };

        return new Paginator(new Callback($data, $count));
    }

    /**
     * {@inheritdoc}
     */
    public function getRecordingTimeLine(string $mode = self::MODE_HOURLY, $limit = 20): array
    {
        list ($pattern, $dateInterval) = $this->modeToDatePattern($mode, $limit);

        $since = new DateTime();
        $since->modify($dateInterval);

        // Get the services from the performer entity
        $metadata = $this->getEntityManager()->getClassMetadata(AbstractPerformerEntity::class);

        $result = [];

        foreach ($metadata->discriminatorMap as $service => $fqcn) {

            $builder = $this->createQueryBuilder('recording');
            $builder
                ->select(
                    sprintf('COUNT(recording.id) as hits, DATE_FORMAT(recording.createdAt, \'%s\') AS interval',
                        $pattern)
                )
                ->setMaxResults($limit)
                ->addGroupBy('interval')
                ->andWhere('recording.service = :service')
                ->andWhere('recording.createdAt >= :since')
                ->setParameter('service', $service)
                ->setParameter('since', $since->format('Y-m-d H:i:s'));

            $query = $builder->getQuery();
            $query->useResultCache(true, 60);

            $result[$service] = $query->getArrayResult();
        }

        return $result;
    }

    /**
     * {@inheritdoc}
     */
    public function getRecordingCount(Criteria $criteria): int
    {
        $builder = $this->createQueryBuilder('recording');
        $builder->select('COUNT(recording.id)');
        $builder->addCriteria($criteria);

        $query = $builder->getQuery();
        $query->useResultCache(true, 30);

        return (int)$query->getSingleScalarResult();
    }

    /**
     * Convert the mode to a date pattern and since limit to increase performance
     *
     * @param string $mode
     * @param int $limit
     *
     * @throws \InvalidArgumentException
     *
     * @return array
     */
    private function modeToDatePattern($mode, $limit)
    {
        switch ($mode) {
            case self::MODE_DAILY:
                return ['%Y-%m-%d', sprintf('-%d days', $limit)];

            case self::MODE_HOURLY:
                return ['%Y-%m-%d %H:00', sprintf('-%d hours', $limit)];

            case self::MODE_MINUTELY:
                return ['%Y-%m-%d %H:%i:00', sprintf('-%d minutes', $limit)];

            default:
                throw new \InvalidArgumentException('Invalid date pattern mode');
        }
    }

    /**
     * {@inheritdoc}
     */
    public function matching(Criteria $criteria): Paginator
    {
        $builder = $this->createQueryBuilder('recording');
        $builder->addCriteria($criteria);

        $expr = $builder->expr();

        $countBuilder = clone $builder;
        $countBuilder
            ->select('COUNT(recording.id)')
            ->resetDQLPart('orderBy')
            ->setMaxResults(null)
            ->setFirstResult(null);

        $countQuery = $countBuilder->getQuery();
        $countQuery->useResultCache(true, 600);

        $count = function () use ($countQuery) {
            return $countQuery->getSingleScalarResult();
        };

        $data = function () use ($builder) {
            return $builder->getQuery()->getResult();
        };

        return new Paginator(new Callback($data, $count));
    }

}
