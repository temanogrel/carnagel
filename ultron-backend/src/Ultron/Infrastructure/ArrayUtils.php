<?php
/**
 *
 *
 *
 */

namespace Ultron\Infrastructure;

use Zend\Stdlib\ArrayUtils as ArrayUtilsBase;

final class ArrayUtils extends ArrayUtilsBase
{
    /**
     * Case insensitive array unique
     *
     * @param string[]|\Traversable $array
     *
     * @return array
     */
    public static function iunique($array)
    {
        $array = static::iteratorToArray($array);

        return array_intersect_key(
            $array,
            array_unique(array_map('strtolower',$array))
        );
    }
}
