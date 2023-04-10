<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Repository;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use DateTime;
use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\EntityRepository;
use Generator;
use Zend\Paginator\Adapter\Callback;
use Zend\Paginator\Paginator;

class PerformerRepository extends EntityRepository implements PerformerRepositoryInterface
{
    /**
     * {@inheritdoc}
     */
    public function getById($id)
    {
        return $this->findOneBy(['id' => $id]);
    }

    /**
     * {@inheritdoc}
     */
    public function getByServiceId($service, $id)
    {
        $builder = $this->getEntityManager()->createQueryBuilder();
        $builder
            ->select('performer')
            ->from(AbstractPerformerEntity::serviceToEntityClassName($service), 'performer')
            ->where('performer.serviceId = :serviceId')
            ->setParameter('serviceId', (string)$id);

        return $builder->getQuery()->getOneOrNullResult();
    }

    /**
     * {@inheritdoc}
     */
    public function getBlacklisted()
    {
        return $this->findBy(['blacklisted' => true]);
    }

    /**
     * {@inheritdoc}
     */
    public function markAllAsOffline($service = null)
    {
        $entity = $service ?
            AbstractPerformerEntity::serviceToEntityClassName($service) :
            AbstractPerformerEntity::class;

        $builder = $this->getEntityManager()->createQueryBuilder();
        $builder
            ->update($entity, 'performer')
            ->where('performer.online = true');

        // Nuke the current viewer count
        $builder
            ->set('performer.online', 'false')
            ->set('performer.currentViewers', 0)
            ->set('performer.peakViewerCount', 0);

        return $builder->getQuery()->execute();
    }

    /**
     * {@inheritdoc}
     */
    public function getPerformerCount($state = null)
    {
        $builder = $this->createQueryBuilder('performer');
        $builder->select('COUNT(performer.id)');

        switch ($state) {
            case 'pending':
                $builder->andWhere('performer.isPendingRecording = true');
                break;

            case 'recording':
                $builder->andWhere('performer.isRecording = true');
                break;

            default:
                $state = null;
        }

        $query = $builder->getQuery();
        $query->useResultCache(true, $state === null ? 300 : 5);

        return (int)$query->getSingleScalarResult();
    }

    /**
     * {@inheritdoc}
     */
    public function search(Criteria $criteria, $service = null, $indexBy = null)
    {
        if ($service) {

            $builder = $this->getEntityManager()->createQueryBuilder();
            $builder->select('performer')
                ->from(AbstractPerformerEntity::serviceToEntityClassName($service), 'performer');

        } else {
            $builder = $this->createQueryBuilder('performer');
        }

        if ($indexBy) {
            $builder->indexBy('performer', sprintf('performer.%s', $indexBy));
        }

        $builder->addCriteria($criteria);

        $countBuilder = clone $builder;
        $countBuilder
            ->select('COUNT(performer.id)')
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
    public function matching(Criteria $criteria, $service = null, $indexBy = null)
    {
        if ($service) {

            $builder = $this->getEntityManager()->createQueryBuilder();
            $builder->select('performer')
                ->from(AbstractPerformerEntity::serviceToEntityClassName($service), 'performer');

        } else {
            $builder = $this->createQueryBuilder('performer');
        }

        if ($indexBy) {
            $builder->indexBy('performer', sprintf('performer.%s', $indexBy));
        }

        $builder->addCriteria($criteria);

        return $builder->getQuery()->getResult();
    }

    /**
     * {@inheritdoc}
     */
    public function updateAllOnline($service)
    {
        $entity = AbstractPerformerEntity::serviceToEntityClassName($service);

        $builder = $this->getEntityManager()->createQueryBuilder();
        $builder
            ->update($entity, 'performer')
            ->set('performer.updatedAt', ':timestamp')
            ->where('performer.online = true')
            ->setParameter('timestamp', (new DateTime())->format('Y-m-d H:i:s'));

        return $builder->getQuery()->execute();
    }

    /**
     * {@inheritdoc}
     */
    public function getAvailableServices()
    {
        return array_keys($this->getClassMetadata()->discriminatorMap);
    }

    /**
     * Get the number of performers with is_recording / is_pending_recording per service
     *
     * @return Generator
     */
    public function getPerformerStats(): Generator
    {
        $conn = $this->getEntityManager()->getConnection();
        $conn->exec('SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED');

        $stmt = $conn->executeQuery('SELECT count(*) AS c, is_recording, is_pending_recording, service  FROM performers WHERE is_pending_recording = TRUE OR is_recording = TRUE  GROUP BY service, is_recording, is_pending_recording;');

        foreach ($stmt->fetchAll() as $row) {

            $state = $row['is_recording'] ? 'recording' : 'pending';

            yield $row['c'] => [$state, $row['service']];
        }

        // NOOP operation that cancels the transaction isolation level
        $conn->exec('COMMIT');
    }
}
