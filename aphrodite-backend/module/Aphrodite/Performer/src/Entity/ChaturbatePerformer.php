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
 * Class ChaturbatePerformer
 *
 * @Entity()
 */
class ChaturbatePerformer extends AbstractPerformerEntity implements SectionAwareInterface
{
    use SectionAwareTrait;

    /**
     * {@inheritdoc}
     */
    public function getService()
    {
        return 'cbc';
    }
}
