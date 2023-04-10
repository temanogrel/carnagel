<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Aphrodite\Recording\Repository\Exception;

use DomainException;

final class NonUniqueRecordingResultException extends DomainException
{
    /**
     * @var array
     */
    private $identifiers;

    /**
     * NonUniqueRecordingResultException constructor.
     *
     * @param array  $identifiers
     */
    public function __construct(array $identifiers = [])
    {
        parent::__construct(sprintf('%d recordings where found', count($identifiers)));

        $this->identifiers = $identifiers;
    }

    public function getIdentifiers():array
    {
        return $this->identifiers;
    }
}
