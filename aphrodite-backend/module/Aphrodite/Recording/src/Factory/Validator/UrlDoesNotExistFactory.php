<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Factory\Validator;

use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\Validator\UrlDoesNotExist;
use Doctrine\Common\Persistence\ObjectRepository;
use Zend\File\Transfer\Adapter\ValidatorPluginManager;
use Zend\ServiceManager\FactoryInterface;
use Zend\ServiceManager\MutableCreationOptionsInterface;
use Zend\ServiceManager\ServiceLocatorInterface;

class UrlDoesNotExistFactory implements FactoryInterface, MutableCreationOptionsInterface
{
    /**
     * @var array
     */
    private $options = [];

    /**
     * Set creation options
     *
     * @param  array $options
     *
     * @return void
     */
    public function setCreationOptions(array $options)
    {
        $this->options = $options;
    }

    /**
     * Create service
     *
     * @param ValidatorPluginManager|ServiceLocatorInterface $validatorManager
     *
     * @return UrlDoesNotExist
     */
    public function createService(ServiceLocatorInterface $validatorManager)
    {
        $serviceLocator = $validatorManager->getServiceLocator();

        /* @var $repository ObjectRepository */
        $repository = $serviceLocator->get('Aphrodite\ObjectManager')->getRepository(UrlEntry::class);

        $options = array_merge($this->options, [
            'object_repository' => $repository,
            'fields'            => ['url'],
        ]);

        return new UrlDoesNotExist($options);
    }
}
