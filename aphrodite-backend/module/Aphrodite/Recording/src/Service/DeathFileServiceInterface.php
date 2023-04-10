<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Service;

use Aphrodite\Recording\Entity\DeathFileEntity;
use Zend\Stdlib\ParametersInterface;

interface DeathFileServiceInterface
{
    /**
     * Create a new death file entry based on the uploaded file
     *
     * @param ParametersInterface $file
     *
     * @return DeathFileEntity
     */
    public function createFromFile(ParametersInterface $file);

    /**
     * Update a death file entity
     *
     * @param DeathFileEntity $file
     *
     * @return void
     */
    public function update(DeathFileEntity $file);

    /**
     * Delete a death file entity and it's file
     *
     * @param DeathFileEntity $file
     *
     * @return void
     */
    public function delete(DeathFileEntity $file);
}
