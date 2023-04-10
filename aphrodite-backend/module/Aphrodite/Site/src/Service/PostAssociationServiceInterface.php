<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\Site\Service;

use Aphrodite\Recording\Entity\RecordingEntity;
use Aphrodite\Site\Entity\Site;

interface PostAssociationServiceInterface
{
    /**
     * Create a new post association to the recording
     *
     * @param RecordingEntity $recording
     * @param Site            $site
     * @param int             $post
     *
     * @return void
     */
    public function create(RecordingEntity $recording, Site $site, $post);
}
