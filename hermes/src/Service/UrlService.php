<?php
/**
 *
 *
 *
 */

namespace Hermes\Service;

use DateTime;
use Doctrine\Common\Persistence\ObjectManager;
use Hermes\Entity\UrlEntity;
use Hermes\Service\Exception\FailedToMatchUrlException;

class UrlService implements UrlServiceInterface
{
    const PREFIX = 'hrm-';

    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * @var UpstoreServiceInterface
     */
    private $upstoreService;

    /**
     * @param ObjectManager           $objectManager
     * @param UpstoreServiceInterface $upstoreService
     */
    public function __construct(ObjectManager $objectManager, UpstoreServiceInterface $upstoreService)
    {
        $this->objectManager  = $objectManager;
        $this->upstoreService = $upstoreService;
    }

    /**
     * Convert a entity id into a string
     *
     * @param int $id
     *
     * @return string
     */
    private function convertIdToString($id)
    {
        $chars  = array_merge(range(0, 9), range('a', 'z'), range('A', 'Z'));
        $length = count($chars);

        $code = '';

        while ($id > $length - 1) {

            // determine the value of the next higher character
            // in the short code should be and prepend
            $code = $chars[(int)fmod($id, $length)] . $code;

            // reset $id to remaining value to be converted
            $id = floor($id / $length);
        }

        return static::PREFIX . $chars[$id] . $code;
    }

    /**
     * {@inheritdoc}
     */
    public function create($url, $host)
    {
        $entity = new UrlEntity();
        $entity->setOriginalUrl($url);
        $entity->setHostname($host);
        $entity->setIsUpstore(true);
        $entity->setUpdatedAt(new DateTime());
        $entity->setCreatedAt(new DateTime());

        $this->objectManager->persist($entity);
        $this->objectManager->flush();

        // Once persisted it gets an id, we use that to generate the short code.
        $entity->setKey($this->convertIdToString($entity->getId()));

        $this->objectManager->flush();

        return $entity;
    }

    /**
     * {@inheritdoc}
     */
    public function update(UrlEntity $url)
    {
        $url->setUpdatedAt(new DateTime());

        $this->objectManager->flush();
    }

    /**
     * {@inheritdoc}
     */
    public function delete(UrlEntity $url)
    {
        $this->objectManager->remove($url);
        $this->objectManager->flush();
    }

    /**
     * {@inheritdoc}
     */
    public function incrementTransmissions(UrlEntity $url)
    {
        $transmissions = $url->getTransmissions();
        $transmissions++;

        $url->setTransmissions($transmissions);

        $this->update($url);
    }
}
