<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Service;

use Aphrodite\Performer\Entity\AbstractPerformerEntity;

interface PerformerServiceInterface
{
    /**
     * Update the performer
     *
     * @param AbstractPerformerEntity $performer
     *
     * @return void
     */
    public function update(AbstractPerformerEntity $performer);
}
