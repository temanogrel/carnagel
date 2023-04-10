<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Service;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Repository\PerformerRepositoryInterface;
use DateTime;
use Doctrine\Common\Collections\Collection;
use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Persistence\ObjectManager;
use Doctrine\ORM\EntityManager;
use Zend\Stdlib\Hydrator\ClassMethods;
use Zend\Stdlib\Hydrator\Filter\MethodMatchFilter;
use Zend\Stdlib\Hydrator\HydratorInterface;

class IntersectionService implements IntersectionServiceInterface
{
    /**
     * @var HydratorInterface
     */
    private $hydrator;

    /**
     * @var EntityManager
     */
    private $objectManager;

    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @param ObjectManager                $objectManager
     * @param PerformerRepositoryInterface $performerRepository
     */
    public function __construct(ObjectManager $objectManager, PerformerRepositoryInterface $performerRepository)
    {
        $this->objectManager       = $objectManager;
        $this->performerRepository = $performerRepository;
    }

    /**
     * Lazy-load the hydrator used to update the performer
     *
     * @return HydratorInterface
     */
    private function getHydrator()
    {
        if ($this->hydrator !== null) {
            return $this->hydrator;
        }

        // todo: not sure how well class methods works, but it's worth a try
        $this->hydrator = new ClassMethods();
        $this->hydrator->getFilter()->addFilter('createdAt', new MethodMatchFilter('setCreatedAt'));
        $this->hydrator->getFilter()->addFilter('updatedAt', new MethodMatchFilter('setUpdatedAt'));

        return $this->hydrator;
    }

    /**
     * Retrieve all performers that belong to the given service and are online.
     *
     * @param string $service
     *
     * @return Collection|AbstractPerformerEntity[]
     */
    private function getPerformers($service)
    {
        $criteria = Criteria::create();
        $criteria->andWhere($criteria->expr()->eq('online', 1));

        // Retrieve
        return $this->performerRepository->matching($criteria, $service, 'serviceId');
    }

    /**
     * {@inheritdoc}
     */
    public function process($service, array $data)
    {
        $found      = [];
        $hydrator   = $this->getHydrator();
        $performers = $this->getPerformers($service);

        foreach ($data as $serviceId => $raw) {

            // Performer is online and available
            if (isset($performers[$serviceId])) {

                $performer = $performers[$serviceId];

                // This is really important since the model server is totally unaware of recording updates.
                if ($performer->isRecording()) {

                    $raw['isRecording']        = true;
                    $raw['isPendingRecording'] = false;
                }

            } else {

                // Attempt to load an offline performer
                $performer = $this->performerRepository->getByServiceId($service, $serviceId);

                // Create the new performer
                if (!$performer) {

                    // We need to know the entity name to get the proper hydrator
                    $className = AbstractPerformerEntity::serviceToEntityClassName($service);

                    /* @var $performer AbstractPerformerEntity */
                    $performer = new $className;
                    $performer->setServiceId($serviceId);
                    $performer->setCreatedAt(new DateTime());
                    $performer->setUpdatedAt(new DateTime());

                    // Persist the bastard
                    $this->objectManager->persist($performer);
                }
            }

            // Don't cause un-required updates
            if (isset($raw['createdAt'])) {
                unset($raw['createdAt']);
            }

            // Don't cause un-required updates
            if (isset($raw['updatedAt'])) {
                unset($raw['updatedAt']);
            }

            // Hydrate the performer
            $hydrator->hydrate($raw, $performer);

            // Add the stage name as an alias
            $performer->setOnline(true);
            $performer->addAlias($raw['stageName']);

            // Add it to the list of found performers
            $found[$serviceId] = $performer;
        }

        // Compare the found with the all performers online and find those that went offline.
        $delete = array_diff(array_keys($performers), array_keys($found));

        // Update the performers that were not found
        foreach ($delete as $serviceId) {
            $performers[$serviceId]->setUpdatedAt(new DateTime());
            $performers[$serviceId]->setOnline(false);
        }

        $this->performerRepository->updateAllOnline($service);
        $this->objectManager->flush();

        return array_values($found);
    }
}
