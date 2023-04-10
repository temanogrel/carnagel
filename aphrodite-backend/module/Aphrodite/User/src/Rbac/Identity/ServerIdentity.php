<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\User\Rbac\Identity;

use ZfcRbac\Identity\IdentityInterface;

class ServerIdentity implements IdentityInterface
{
    /**
     * {@inheritdoc}
     */
    public function getRoles()
    {
        return ['server'];
    }
}
