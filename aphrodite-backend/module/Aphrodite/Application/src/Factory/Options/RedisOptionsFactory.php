<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Aphrodite\Application\Factory\Options;

use Aphrodite\Application\Options\RedisOptions;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class RedisOptionsFactory implements FactoryInterface
{
    /**
     * Create service
     *
     * @param ServiceLocatorInterface $serviceLocator
     *
     * @return mixed
     */
    public function createService(ServiceLocatorInterface $serviceLocator)
    {
        /* @var array $config */
        $config = $serviceLocator->get('config');

        return new RedisOptions($config['aphrodite']['options'][RedisOptions::class]);
    }
}
