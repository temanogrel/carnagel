<?php
/**
 *
 *
 *
 */
declare(strict_types=1);

namespace Aphrodite\Logger\Adapter;

final class BlackHoleAdapter extends AbstractAdapter
{
    public function write(array $data, string $type = null): bool
    {
        return true;
    }
}
