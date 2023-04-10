<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\User\Factory\Rbac;

use Aphrodite\User\Rbac\IdentityProvider;
use Zend\Authentication\AuthenticationService;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class IdentityProviderFactory implements FactoryInterface
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
        return new IdentityProvider(
            $serviceLocator->get(AuthenticationService::class),
            $serviceLocator->get('request'),
            $serviceLocator->get('config')['aphrodite']['serverAccessToken']
        );
    }
}
