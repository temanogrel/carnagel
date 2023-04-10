<?php
/**
 *
 *
 *  AB
 */

declare(strict_types = 1);

namespace Aphrodite\Recording\Service;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Performer\Stdlib\SectionAwareInterface;
use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\RecordingEvent;
use Aphrodite\Recording\Service\Listener\UpdateRecordingCountListener;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;
use Zend\EventManager\EventManagerAwareTrait;

class RecordingService implements RecordingServiceInterface
{
    use EventManagerAwareTrait;

    const PERMISSION_READ = 'aphrodite:recording:read';
    const PERMISSION_UPDATE = 'aphrodite:recording:update';
    const PERMISSION_CREATE = 'aphrodite:recording:create';
    const PERMISSION_DELETE = 'aphrodite:recording:delete';

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
     * This is called by the event manager aware trait and hooks in basic listeners.
     *
     * @return void
     */
    private function attachDefaultListeners()
    {
        $this
            ->getEventManager()
            ->attach(new UpdateRecordingCountListener());
    }

    /**
     * {@inheritdoc}
     */
    public function create(RecordingEntity $recording, AbstractPerformerEntity $performer = null): void
    {
        if ($recording->getCreatedAt() === null) {
            $recording->setCreatedAt(new DateTime());
            $recording->setUpdatedAt(new DateTime());
        }

        $recording->setPerformer($performer);

        // Associate this for performance!
        $recording->setService($performer->getService());
        $recording->setStageName($performer->getStageName());

        if ($performer instanceof SectionAwareInterface) {
            $recording->setSection($performer->getSection());
        }

        $event = new RecordingEvent(RecordingEvent::NEW_RECORDING, $this, ['recording' => $recording]);

        $this
            ->getEventManager()
            ->trigger($event);

        $this->objectManager->persist($recording);
        $this->objectManager->flush();
    }

    /**
     * {@inheritdoc}
     */
    public function update(RecordingEntity $recording): void
    {
        $recording->setUpdatedAt(new DateTime());

        $this->objectManager->flush();
    }

    /**
     * {@inheritdoc}
     */
    public function delete(RecordingEntity $recording): void
    {
        $this->objectManager->remove($recording);
        $this->objectManager->flush();
    }
}
