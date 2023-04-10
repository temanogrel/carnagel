<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Hydrator;

use Aphrodite\Site\Entity\Site;
use Aphrodite\Stdlib\Hydrator\Strategy\DateTimeStrategy;
use DateTime;
use Zend\Hydrator\AbstractHydrator;

class SiteHydrator extends AbstractHydrator
{
    public function __construct()
    {
        parent::__construct();

        $this->addStrategy('createdAt', new DateTimeStrategy(DateTime::RFC3339));
        $this->addStrategy('updatedAt', new DateTimeStrategy(DateTime::RFC3339));
    }

    /**
     * Extract values from an object
     *
     * @param  object $object
     *
     * @return array
     */
    public function extract($object)
    {
        if (!$object instanceof Site) {
            throw new \DomainException;
        }

        return [
            'id'      => $this->extractValue('id', $object->getId()),
            'name'    => $this->extractValue('name', $object->getName()),
            'enabled' => $this->extractValue('enabled', $object->isEnabled()),
            'sources' => $this->extractValue('sources', $object->getSources()),

            // Api
            'apiUri'   => $this->extractValue('apiUri', $object->getApiUri()),
            'username' => $this->extractValue('username', $object->getUsername()),
            'password' => $this->extractValue('password', $object->getPassword()),

            // Dates
            'createdAt' => $this->extractValue('createdAt', $object->getCreatedAt()),
            'updatedAt' => $this->extractValue('updatedAt', $object->getUpdatedAt())
        ];
    }

    /**
     * Hydrate $object with the provided $data.
     *
     * @param  array  $data
     * @param  object $object
     *
     * @return object
     */
    public function hydrate(array $data, $object)
    {
        if (!$object instanceof Site) {
            throw new \DomainException;
        }

        $object->setName($this->hydrateValue('name', $data['name']));
        $object->setEnabled($this->hydrateValue('enabled', $data['enabled']));
        $object->setSources($this->hydrateValue('sources', $data['sources']));
        $object->setApiUri($this->hydrateValue('apiUri', $data['apiUri']));
        $object->setUsername($this->hydrateValue('username', $data['username']));
        $object->setPassword($this->hydrateValue('password', $data['password']));

        return $object;
    }
}
