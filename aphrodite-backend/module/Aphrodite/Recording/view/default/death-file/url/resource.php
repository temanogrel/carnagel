<?php
return [
    'id'           => $this->url->getId(),
    'url'          => $this->url->getUrl(),
    'state'        => $this->url->getState(),
    'filename'     => $this->url->getFilename(),
    'ignoreReason' => $this->url->getIgnoreReason(),
    'hermesId'     => $this->url->getHermesId(),
    'recordingId'  => $this->url->getRecording() ? $this->url->getRecording()->getId() : null,
    'createdAt'    => $this->url->getCreatedAt()->format(DateTime::RFC3339),
    'updatedAt'    => $this->url->getUpdatedAt()->format(DateTime::RFC3339),
];
