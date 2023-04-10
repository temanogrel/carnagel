<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Service\Listener;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\RecordingEvent;
use Zend\EventManager\EventManagerInterface;
use Zend\EventManager\ListenerAggregateInterface;
use Zend\EventManager\ListenerAggregateTrait;

class UpdateRecordingCountListener implements ListenerAggregateInterface
{
    use ListenerAggregateTrait;

    /**
     * Attach one or more listeners
     *
     * Implementors may add an optional $priority argument; the EventManager
     * implementation will pass this to the aggregate.
     *
     * @param EventManagerInterface $events
     *
     * @return void
     */
    public function attach(EventManagerInterface $events)
    {
        $this->listeners[] = $events->attach(RecordingEvent::NEW_RECORDING, [$this, 'update']);
    }

    /**
     * Simply increments the number of recordings.
     *
     * Since this does a very simple increment, it's possible that during a race condition the update is wrong.
     * But since this is not important information it's fine to handle it this way
     *
     * @param RecordingEvent $event
     */
    public function update(RecordingEvent $event)
    {
        /* @var $recording RecordingEntity */
        $recording = $event->getParam('recording');

        if ($recording->getPerformer() !== null) {
            $performer = $recording->getPerformer();
            $performer->setRecordingCount($performer->getRecordingCount() + 1);
        }
    }
}
