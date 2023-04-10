<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Repository;

use Aphrodite\Recording\Entity\DeathFileEntity;
use Doctrine\Common\Collections\Selectable;

interface DeathFileRepositoryInterface extends Selectable
{
    /**
     * Retrieve a death file entity by it's primary id
     *
     * @param int $id
     *
     * @return DeathFileEntity|null
     */
    public function getById($id);

    /**+
     * Retrieve a death file entity by it's upsto.re url
     *
     * @param string $url
     *
     * @return DeathFileEntity|null
     */
    public function getByUrl($url);
}
