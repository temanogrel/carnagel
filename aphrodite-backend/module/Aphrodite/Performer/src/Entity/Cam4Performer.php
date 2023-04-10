<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Entity;

use Aphrodite\Performer\Stdlib\SectionAwareInterface;
use Aphrodite\Performer\Stdlib\SectionAwareTrait;
use Doctrine\ORM\Mapping\Entity;

/**
 * Class Cam4Performer
 *
 * @Entity()
 */
class Cam4Performer extends AbstractPerformerEntity implements SectionAwareInterface
{
    use SectionAwareTrait;

    /**
     * {@inheritdoc}
     */
    public function getService()
    {
        return 'cam4';
    }
}
