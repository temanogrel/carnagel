<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Site\Entity;

use DateTime;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;
use Doctrine\ORM\Mapping\GeneratedValue;
use Doctrine\ORM\Mapping\Id;
use Doctrine\ORM\Mapping\OneToMany;
use Doctrine\ORM\Mapping\Table;

/**
 * Class Site
 *
 * @Entity(repositoryClass="Aphrodite\Site\Repository\SiteRepository")
 * @Table(name="sites")
 */
class Site
{
    /**
     * @var int
     *
     * @Id
     * @Column(type="integer")
     * @GeneratedValue()
     */
    protected $id;

    /**
     * @var bool
     *
     * @Column(type="boolean")
     */
    protected $enabled = true;

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
    protected $apiUri;

    /**
     * @var string
     *
     * @Column()
     */
    protected $username;

    /**
     * @var string
     *
     * @Column()
     */
    protected $password;

    /**
     * @var array
     *
     * @Column(type="json_array")
     */
    protected $sources = [];

    /**
     * @var PostAssociation[]
     *
     * @OneToMany(targetEntity="PostAssociation", mappedBy="site")
     */
    protected $posts;

    /**
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $createdAt;

    /**
     * @var DateTime
     *
     * @Column(type="datetime")
     */
    protected $updatedAt;

    public function __construct()
    {
        $this->posts = new ArrayCollection();
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
        $this->id = $id;
    }

    /**
     * @return boolean
     */
    public function isEnabled()
    {
        return $this->enabled;
    }

    /**
     * @param boolean $enabled
     */
    public function setEnabled($enabled)
    {
        $this->enabled = (bool)$enabled;
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
        $this->name = $name;
    }

    /**
     * @return string
     */
    public function getApiUri()
    {
        return $this->apiUri;
    }

    /**
     * @param string $apiUri
     */
    public function setApiUri($apiUri)
    {
        $this->apiUri = $apiUri;
    }

    /**
     * @return string
     */
    public function getUsername()
    {
        return $this->username;
    }

    /**
     * @param string $username
     */
    public function setUsername($username)
    {
        $this->username = $username;
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
        $this->password = $password;
    }

    /**
     * @return array
     */
    public function getSources()
    {
        return $this->sources;
    }

    /**
     * @param string[] $sources
     */
    public function setSources(array $sources)
    {
        $this->sources = $sources;
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
}
