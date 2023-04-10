<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\Site\Service;

use Aphrodite\Site\Entity\Site;

interface SiteServiceInterface
{
    public function remove(Site $site);
    public function update(Site $site);
    public function create(Site $site);
}
