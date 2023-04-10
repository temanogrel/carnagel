<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\InputFilter;

use Zend\InputFilter\FileInput;
use Zend\InputFilter\InputFilter;
use Zend\Validator\File\MimeType;

class DeathFileUploadInputFilter extends InputFilter
{
    public function init()
    {
        $this->add([
            'type'       => FileInput::class,
            'name'       => 'file',
            'validators' => [
                [
                    'name'    => MimeType::class,
                    'options' => [
                        'mimeType' => 'application/x-gzip'
                    ]
                ]
            ]
        ]);
    }
}
