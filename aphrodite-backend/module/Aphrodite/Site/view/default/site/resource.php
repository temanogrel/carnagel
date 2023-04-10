<?php
/**
 *
 *
 *
 */

return [
    'id'      => $this->site->getId(),
    'name'    => $this->site->getName(),
    'enabled' => $this->site->isEnabled(),
    'sources' => $this->site->getSources(),

    // Api
    'apiUri'   => $this->site->getApiUri(),
    'username' => $this->site->getUsername(),
    'password' => $this->site->getPassword(),

    // Dates
    'createdAt' => $this->site->getCreatedAt()->format(DateTime::RFC3339),
    'updatedAt' => $this->site->getUpdatedAt()->format(DateTime::RFC3339),
];
