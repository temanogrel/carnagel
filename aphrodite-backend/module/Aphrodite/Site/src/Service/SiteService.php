<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\Site\Service;

use Aphrodite\Site\Entity\Site;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;

class SiteService implements SiteServiceInterface
{
    const PERMISSION_READ = 'aphrodite:site:read';
    const PERMISSION_UPDATE = 'aphrodite:site:update';
    const PERMISSION_CREATE = 'aphrodite:site:create';
    const PERMISSION_DELETE = 'aphrodite:site:delete';

    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * @param ObjectManager $objectManager
     */
    public function __construct(ObjectManager $objectManager)
    {
        $this->objectManager = $objectManager;
    }

    public function update(Site $site)
    {
        $site->setUpdatedAt(new DateTime());

        $this->objectManager->flush();
    }

    public function create(Site $site)
    {
        $site->setUpdatedAt(new DateTime());
        $site->setCreatedAt(new DateTime());

        $this->objectManager->persist($site);
        $this->objectManager->flush();
    }

    public function remove(Site $site)
    {
        $this->objectManager->remove($site);
        $this->objectManager->flush();
    }
}
