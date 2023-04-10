<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\InputFilter;

use Zend\InputFilter\InputFilter;

final class CreateAddressInputFilter extends InputFilter
{
    public function init()
    {
        $this->add([
            'name' => 'webhookId',
        ]);
    }
}
