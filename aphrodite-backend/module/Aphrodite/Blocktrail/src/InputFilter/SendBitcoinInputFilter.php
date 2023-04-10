<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\InputFilter;

use Zend\Filter\Callback;
use Zend\InputFilter\InputFilter;

final class SendBitcoinInputFilter extends InputFilter
{
    public function init()
    {
        $this->add([
            'name' => 'address'
        ]);

        $this->add([
            'name' => 'bitcoinAmount',
            'filters' => [
                'name' => Callback::class,
                'options' => [
                    'callback' => function($value) {
                        return (float) $value;
                    }
                ]
            ]
        ]);
    }
}
