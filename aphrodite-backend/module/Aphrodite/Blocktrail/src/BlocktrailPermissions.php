<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail;

final class BlocktrailPermissions
{
    const CREATE_ADDRESS = 'aphrodite:blocktrail:create-address';
    const SEND_BITCOIN   = 'aphrodite:blocktrail:send-bitcoin';

    private function __construct()
    {
    }
}
