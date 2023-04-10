<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\InputFilter\DeathFile;

use Aphrodite\Recording\Service\DeathFile\UrlService;
use Zend\InputFilter\InputFilter;
use Zend\Validator\InArray;

class UrlUpdateInputFilter extends InputFilter
{
    public function init()
    {
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
                            UrlService::STATE_REMOVED
                        ]
                    ]
                ]
            ]
        ]);
    }
}
