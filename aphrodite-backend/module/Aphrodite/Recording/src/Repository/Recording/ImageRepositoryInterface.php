<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Aphrodite\Recording\Repository\Recording;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Recording\Entity\Recording\ImageEntity;
use Aphrodite\Recording\Repository\Exception\ImageNotFoundException;
use Doctrine\Common\Collections\Selectable;

interface ImageRepositoryInterface extends Selectable
{
    /**
     * @param int $id
     *
     * @throws ImageNotFoundException
     *
     * @return ImageEntity
     */
    public function getById(int $id): ImageEntity;

    /**
     * @param RecordingEntity $recording
     *
     * @return array
     */
    public function getForRecording(RecordingEntity $recording): array;
}
