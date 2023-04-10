<?php
/**
 *
 *
 *
 */
declare(strict_types=1);

namespace Aphrodite\Logger\Adapter;

abstract class AbstractAdapter implements AdapterInterface
{
    abstract public function write(array $data, string $type = null): bool;
}
