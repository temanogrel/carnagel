<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\InputFilter\DeathFile;

use Aphrodite\Recording\Service\DeathFile\UrlService;
use Aphrodite\Recording\Validator\RecordingExists;
use Aphrodite\Recording\Validator\UrlDoesNotExist;
use Zend\Filter\StringTrim;
use Zend\Filter\ToInt;
use Zend\InputFilter\InputFilter;
use Zend\Validator\InArray;
use Zend\Validator\Uri;

class UrlAddInputFilter extends InputFilter
{
    public function init()
    {
        $this->add([
            'name'       => 'url',
            'validators' => [
                [
                    'name'    => Uri::class,
                    'options' => [
                        'allowRelative' => false,
                    ],
                ],
                [
                    'name' => UrlDoesNotExist::class,
                ],
            ],
        ]);

        $this->add([
            'name'        => 'recording',
            'required'    => false,
            'allow_empty' => true,
            'validators'  => [
                [
                    'name' => RecordingExists::class,
                ],
            ],

            'filters' => [
                [
                    'name' => ToInt::class,
                ],
            ],
        ]);

        $this->add([
            'name'       => 'state',
            'validators' => [
                [
                    'name'    => InArray::class,
                    'options' => [
                        'haystack' => [
                            UrlService::STATE_IGNORED,
                            UrlService::STATE_IN_PROGRESS,
                            UrlService::STATE_PENDING,
                            UrlService::STATE_REMOVED,
                        ],
                    ],
                ],
            ],
        ]);

        $this->add([
            'name'        => 'ignoreReason',
            'required'    => false,
            'allow_empty' => true,
            'filters'     => [
                [
                    'name' => StringTrim::class,
                ],
            ],
        ]);

        $this->add([
            'name'     => 'filename',
            'required' => true,
        ]);

        // todo: write a validator that check hermes, perhaps ?
        $this->add([
            'name'        => 'hermesId',
            'required'    => false,
            'allow_empty' => true,
            'filters'     => [
                [
                    'name' => ToInt::class,
                ],
            ],
        ]);
    }
}
