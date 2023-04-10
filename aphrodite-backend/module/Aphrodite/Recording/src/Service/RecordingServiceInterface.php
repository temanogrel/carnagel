<?php
/**
 *
 *
 *  AB
 */

declare(strict_types=1);

namespace Aphrodite\Recording\Service;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;
use Aphrodite\Recording\Entity\RecordingEntity;
use Zend\EventManager\EventsCapableInterface;

interface RecordingServiceInterface extends EventsCapableInterface
{
    /**
     * Create a associated recording
     *
     * @param RecordingEntity         $recording
     * @param AbstractPerformerEntity $performer
     *
     * @return void
     */
    public function create(RecordingEntity $recording, AbstractPerformerEntity $performer = null): void;

    /**
     * Update a recording
     *
     * @param RecordingEntity $recording
     *
     * @return void
     */
    public function update(RecordingEntity $recording): void;

    /**
     * Delete a recording
     *
     * @param RecordingEntity $recording
     *
     * @return void
     */
    public function delete(RecordingEntity $recording): void;
}
