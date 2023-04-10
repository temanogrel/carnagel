<?php

declare(strict_types=1);

namespace Aphrodite\Stdlib\Hydrator\Strategy;

use DateTime;
use Zend\Hydrator\Strategy\DefaultStrategy;

/**
 * https://juriansluiman.nl/article/125/strategies-for-hydrators-a-practical-use-case
 *
 * Class DateTimeStrategy
 * @package MyModule\Hydrator\Strategy
 */
class DateTimeStrategy extends DefaultStrategy
{
    /**
     * {@inheritdoc}
     *
     * Convert a string value into a DateTime object
     */
    public function hydrate($value)
    {
        if (is_string($value)) {
            $value = new DateTime($value);
        }

        return $value;
    }

    /**
     * {@inheritdoc}
     *
     * Convert a string value into a DateTime object
     */
    public function extract($value)
    {
        if (null == $value) {
            return null;
        }

        return $value->format(DateTime::ISO8601);
    }
}
