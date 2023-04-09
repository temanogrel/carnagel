<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure;

use Doctrine\Common\Persistence\ObjectRepository;
use Doctrine\ORM\EntityManagerInterface;
use Doctrine\ORM\Mapping\ClassMetadata;
use Doctrine\ORM\Repository\RepositoryFactory as RepositoryFactoryInterface;
use Ultron\Domain\Service\CacheService;
use Ultron\Infrastructure\Repository\CacheServiceAwareInterface;

class RepositoryFactory implements RepositoryFactoryInterface
{
    /**
     * @var CacheService
     */
    private $cacheService;

    /**
     * @var ObjectRepository[]
     */
    private $repositories = [];

    /**
     * RepositoryFactory constructor.
     *
     * @param CacheService $cacheService
     */
    public function __construct(CacheService $cacheService)
    {
        $this->cacheService = $cacheService;
    }

    /**
     * Gets the repository for an entity class.
     *
     * @param \Doctrine\ORM\EntityManagerInterface $entityManager The EntityManager instance.
     * @param string                               $entityName    The name of the entity.
     *
     * @return \Doctrine\Common\Persistence\ObjectRepository
     */
    public function getRepository(EntityManagerInterface $entityManager, $entityName)
    {
        $hash = $entityManager->getClassMetadata($entityName)->getName() . spl_object_hash($entityManager);

        if (isset($this->repositories[$hash])) {
            return $this->repositories[$hash];
        }

        return $this->repositories[$hash] = $this->createRepository($entityManager, $entityName);
    }

    private function createRepository(EntityManagerInterface $entityManager, string $entityName):ObjectRepository
    {
        /* @var $metadata ClassMetadata */
        $metadata            = $entityManager->getClassMetadata($entityName);
        $repositoryClassName = $metadata->customRepositoryClassName
            ?: $entityManager->getConfiguration()->getDefaultRepositoryClassName();

        if (in_array(CacheServiceAwareInterface::class, class_implements($repositoryClassName), true)) {
            return new $repositoryClassName($entityManager, $metadata, $this->cacheService);
        }

        return new $repositoryClassName($entityManager, $metadata);
    }
}
