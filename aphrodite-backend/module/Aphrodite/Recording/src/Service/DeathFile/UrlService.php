<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Service\DeathFile;

use Aphrodite\Recording\Entity\DeathFile\UrlEntry;
use Aphrodite\Recording\Entity\DeathFileEntity;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;

class UrlService implements UrlServiceInterface
{
    const STATE_PENDING = 'pending';
    const STATE_REMOVED = 'removed';
    const STATE_IGNORED = 'ignored';
    const STATE_IN_PROGRESS = 'in-progress';

    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * UrlService constructor.
     *
     * @param ObjectManager $objectManager
     */
    public function __construct(ObjectManager $objectManager)
    {
        $this->objectManager = $objectManager;
    }

    /**
     * @inheritDoc
     */
    public function create(UrlEntry $url, DeathFileEntity $deathFile)
    {
        $url->setCreatedAt(new DateTime());
        $url->setUpdatedAt(new DateTime());
        $url->setDeathFile($deathFile);

        $this->objectManager->persist($url);
        $this->objectManager->flush();
    }

    /**
     * @inheritDoc
     */
    public function update(UrlEntry $url)
    {
        $url->setUpdatedAt(new DateTime());

        $this->objectManager->flush();
    }

    /**
     * @inheritDoc
     */
    public function remove(UrlEntry $url)
    {
        $this->objectManager->remove($url);
        $this->objectManager->flush();
    }
}
