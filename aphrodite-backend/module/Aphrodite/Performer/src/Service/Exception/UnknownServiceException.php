<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Service\Exception;

use DomainException;

class UnknownServiceException extends DomainException
{
    /**
     * Wrapper
     *
     * @param string $service
     *
     * @return static
     */
    public static function unknownService($service)
    {
        return new static(sprintf('Unknown service %s provided', $service));
    }
}
