<?php

return [
    'id'        => $this->file->getId(),
    'location'  => $this->file->getLocation(),
    'entries'   => $this->file->getEntries(),
    'ignored'   => $this->file->getIgnored(),
    'pending'   => $this->file->getPending(),
    'createdAt' => $this->file->getCreatedAt()->format(DateTime::RFC3339),
    'updatedAt' => $this->file->getUpdatedAt()->format(DateTime::RFC3339),
];
