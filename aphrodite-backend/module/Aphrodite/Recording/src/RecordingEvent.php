<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording;

use Zend\EventManager\Event;

final class RecordingEvent extends Event
{
    const NEW_RECORDING = 'aphrodite:recording.new';
}
