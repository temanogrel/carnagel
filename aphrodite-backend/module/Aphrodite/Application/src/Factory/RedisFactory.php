<?php
/**
 *
 *
 *  AB
 */

declare(strict_types = 1);

namespace Aphrodite\Application\Factory;

use Aphrodite\Application\Options\RedisOptions;
use Redis;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RedisFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return Redis
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        /* @var $options RedisOptions */
        $options = $serviceLocator->get(RedisOptions::class);

        $redis = new Redis();
        $redis->connect($options->getHost(), $options->getPort());

        return $redis;
    }
}
