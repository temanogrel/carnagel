<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\InputFilter;

use Zend\I18n\Validator\IsInt;
use Zend\InputFilter\InputFilter;
use Zend\Validator\LessThan;

class DeathFileUpdateInputFilter extends InputFilter
{
    public function init()
    {
        $this->add([
            'name'       => 'ignored',
            'required'   => false,
            'validators' => [
                [
                    'name' => IsInt::class
                ],
            ]
        ]);

        $this->add([
            'name'       => 'pending',
            'required'   => false,
            'validators' => [
                [
                    'name' => IsInt::class
                ],
            ]
        ]);

        $this->add([
            'name'       => 'entries',
            'validators' => [
                [
                    'name' => IsInt::class
                ]
            ]
        ]);
    }
}
