<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Service;

interface IntersectionServiceInterface
{
    /**
     * Process the incoming data
     *
     * @param string $service
     * @param array  $data
     *
     * @return int
     */
    public function process($service, array $data);
}
