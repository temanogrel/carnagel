<?php
/**
 *
 *
 *
 */

use Aphrodite\Performer\Entity\ChaturbatePerformer;
use Aphrodite\Performer\Entity\MyFreeCamsPerformer;

$default = [
    'id'                 => $this->performer->getId(),
    'serviceId'          => $this->performer->getServiceId(),
    'stageName'          => $this->performer->getStageName(),
    'aliases'            => $this->performer->getAliases(),
    'service'            => $this->performer->getService(),

    // count
    'currentViewers'     => $this->performer->getCurrentViewers(),
    'peakViewerCount' => $this->performer->getPeakViewerCount(),
    'recordingCount'  => $this->performer->getRecordingCount(),

    // Boolean
    'online'             => $this->performer->isOnline(),
    'blacklisted'        => $this->performer->isBlacklisted(),
    'isRecording'        => $this->performer->isRecording(),
    'isPendingRecording' => $this->performer->isPendingRecording(),

    // Dates
    'createdAt' => $this->performer->getCreatedAt()->format(DateTime::RFC3339),
    'updatedAt' => $this->performer->getUpdatedAt()->format(DateTime::RFC3339),
];

if ($this->performer instanceof ChaturbatePerformer) {

    $default += [
        'section' => $this->performer->getSection()
    ];

} else if ($this->performer instanceof MyFreeCamsPerformer) {

    $default += [
        'videoState'  => $this->performer->getVideoState(),
        'camScore'    => $this->performer->getCamScore(),
        'camServer'   => $this->performer->getCamServer(),
        'missMfcRank' => $this->performer->getMissMfcRank(),
        'accessLevel' => $this->performer->getAccessLevel()
    ];
}

return $default;
