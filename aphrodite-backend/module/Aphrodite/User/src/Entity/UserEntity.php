<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\User\Entity;

use DateTime;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\Table;
use ZfcRbac\Identity\IdentityInterface;
use ZfrOAuth2\Server\Entity\TokenOwnerInterface;

/**
 * Class UserEntity
 *
 * @Entity()
 * @Table(name="users")
 */
class UserEntity implements TokenOwnerInterface, IdentityInterface
{
    /**
     * @var int
     *
     * @Id()
     * @GeneratedValue()
     * @Column(type="integer")
     */
    protected $id;

    /**
     * @var string
     *
     * @Column(unique=true)
     */
    protected $email;

    /**
     * @var string
     *
     * @Column()
     */
    protected $name;

    /**
     * @var string
     *
     * @Column()
     */
    protected $surname;

    /**
     * @var string
     *
     * @Column()
     */
    protected $password;

    /**
     * @var boolean
     *
     * @Column(type="boolean")
     */
    protected $activated = false;

    /**
     * @var boolean
     *
     * @Column(type="boolean")
     */
    protected $blocked = false;

    /**
     * @var Datetime
     *
     * @Column(type="datetime")
     */
    protected $createdAt;

    /**
     * @var Datetime
     *
     * @Column(type="datetime")
     */
    protected $updatedAt;

    /**
     * @var Datetime
     */
    protected $deletedAt;

    /**
     * {@inheritdoc}
     */
    public function getTokenOwnerId()
    {
        return $this->id;
    }

    /**
     * {@inheritdoc}
     */
    public function getRoles()
    {
        return ['admin'];
    }

    /**
     * @return int
     */
    public function getId()
    {
        return $this->id;
    }

    /**
     * @param int $id
     */
    public function setId($id)
    {
        $this->id = (int)$id;
    }

    /**
     * @return string
     */
    public function getEmail()
    {
        return $this->email;
    }

    /**
     * @param string $email
     */
    public function setEmail($email)
    {
        $this->email = (string)$email;
    }

    /**
     * @return string
     */
    public function getName()
    {
        return $this->name;
    }

    /**
     * @param string $name
     */
    public function setName($name)
    {
        $this->name = (string)$name;
    }

    /**
     * @return string
     */
    public function getSurname()
    {
        return $this->surname;
    }

    /**
     * @param string $surname
     */
    public function setSurname($surname)
    {
        $this->surname = (string)$surname;
    }

    /**
     * @return string
     */
    public function getPassword()
    {
        return $this->password;
    }

    /**
     * @param string $password
     */
    public function setPassword($password)
    {
        $this->password = (string)$password;
    }

    /**
     * @return boolean
     */
    public function isActivated()
    {
        return $this->activated;
    }

    /**
     * @param boolean $activated
     */
    public function setActivated($activated)
    {
        $this->activated = (bool)$activated;
    }

    /**
     * @return boolean
     */
    public function isBlocked()
    {
        return $this->blocked;
    }

    /**
     * @param boolean $blocked
     */
    public function setBlocked($blocked)
    {
        $this->blocked = (bool)$blocked;
    }

    /**
     * @return DateTime
     */
    public function getCreatedAt()
    {
        return $this->createdAt;
    }

    /**
     * @param DateTime $createdAt
     */
    public function setCreatedAt(DateTime $createdAt)
    {
        $this->createdAt = $createdAt;
    }

    /**
     * @return DateTime
     */
    public function getUpdatedAt()
    {
        return $this->updatedAt;
    }

    /**
     * @param DateTime $updatedAt
     */
    public function setUpdatedAt(DateTime $updatedAt)
    {
        $this->updatedAt = $updatedAt;
    }

    /**
     * @return DateTime
     */
    public function getDeletedAt()
    {
        return $this->deletedAt;
    }

    /**
     * @param DateTime $deletedAt
     */
    public function setDeletedAt(DateTime $deletedAt = null)
    {
        $this->deletedAt = $deletedAt;
    }
}
