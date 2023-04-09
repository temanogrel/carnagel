<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Service\ValueObject;

final class PageInformation
{
    /**
     * @var int
     */
    private $maxId;

    /**
     * @var int
     */
    private $count;

    /**
     * PageInformation constructor.
     * @param int $maxId
     * @param int $count
     */
    public function __construct(int $maxId, int $count)
    {
        $this->maxId = $maxId;
        $this->count = $count;
    }

    /**
     * @return int
     */
    public function getMaxId(): int
    {
        return $this->maxId;
    }

    /**
     * @return int
     */
    public function getCount(): int
    {
        return $this->count;
    }
}
