<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Service;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;

class PerformerService implements PerformerServiceInterface
{
    const PERMISSION_READ = 'aphrodite:performer:read';
    const PERMISSION_UPDATE = 'aphrodite:performer:update';
    const PERMISSION_INTERSECT = 'aphrodite:performer:intersect';

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

    /**
     * {@inheritdoc}
     */
    public function update(AbstractPerformerEntity $performer)
    {
        // Check if the performers stageName has changed
        $stageName = $performer->getStageName();

        if (!$performer->hasAlias($stageName)) {
            $performer->addAlias($stageName);
        }

        // Mark as updated
        $performer->setUpdatedAt(new DateTime());

        // Persist changes
        $this->objectManager->flush();
    }
}
