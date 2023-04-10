<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\Site\Service;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Site\Entity\PostAssociation;
use Aphrodite\Site\Entity\Site;
use Doctrine\Common\Persistence\ObjectManager;

class PostAssociationService implements PostAssociationServiceInterface
{
    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * @param ObjectManager $objectManager
     */
    public function __construct(ObjectManager $objectManager)
    {
        $this->objectManager = $objectManager;
    }

    /**
     * {@inheritdoc}
     */
    public function create(RecordingEntity $recording, Site $site, $post)
    {
        $association = new PostAssociation();
        $association->setSite($site);
        $association->setPostId($post);
        $association->setRecording($recording);

        $this->objectManager->persist($association);
        $this->objectManager->flush();
    }

    public function delete(PostAssociation $association)
    {
        $this->objectManager->remove($association);
        $this->objectManager->flush();
    }
}
