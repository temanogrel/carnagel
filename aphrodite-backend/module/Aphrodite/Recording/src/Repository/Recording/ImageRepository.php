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
use Doctrine\ORM\EntityRepository;

class ImageRepository extends EntityRepository  implements ImageRepositoryInterface
{
    /**
     * @inheritDoc
     */
    public function getById(int $id): ImageEntity
    {
        /* @var $image ImageEntity */
        $image = $this->findOneBy(['id' => $id]);
        if (!$image) {
            throw new ImageNotFoundException();
        }

        return $image;
    }

    /**
     * @inheritDoc
     */
    public function getForRecording(RecordingEntity $recording): array
    {
        return $this->findBy(['recording' => $recording->getId()]);
    }
}
