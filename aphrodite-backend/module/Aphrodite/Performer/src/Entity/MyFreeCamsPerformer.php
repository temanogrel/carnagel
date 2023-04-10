<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Performer\Entity;

use Doctrine\ORM\Mapping\Column;
use Doctrine\ORM\Mapping\Entity;

/**
 * Class MyFreeCamsPerformer
 *
 * @Entity()
 */
class MyFreeCamsPerformer extends AbstractPerformerEntity
{
    /**
     * @var int
     *
     * @Column(type="smallint")
     */
    protected $videoState;

    /**
     * @var int
     *
     * @Column(type="integer")
     */
    protected $camScore;

    /**
     * @var int
     *
     * @Column(type="smallint")
     */
    protected $camServer;

    /**
     * @var int
     *
     * @Column(type="smallint")
     */
    protected $missMfcRank;

    /**
     * @var int
     *
     * @Column(type="smallint")
     */
    protected $accessLevel;

    /**
     * {@inheritdoc}
     */
    public function getService()
    {
        return 'mfc';
    }

    /**
     * @return int
     */
    public function getVideoState()
    {
        return $this->videoState;
    }

    /**
     * @param int $videoState
     */
    public function setVideoState($videoState)
    {
        $this->videoState = (int)$videoState;
    }

    /**
     * @return int
     */
    public function getCamScore()
    {
        return $this->camScore;
    }

    /**
     * @param int $camScore
     */
    public function setCamScore($camScore)
    {
        $this->camScore = (int)$camScore;
    }

    /**
     * @return int
     */
    public function getCamServer()
    {
        return $this->camServer;
    }

    /**
     * @param int $camServer
     */
    public function setCamServer($camServer)
    {
        $this->camServer = (int)$camServer;
    }

    /**
     * @return int
     */
    public function getMissMfcRank()
    {
        return $this->missMfcRank;
    }

    /**
     * @param int $missMfcRank
     */
    public function setMissMfcRank($missMfcRank)
    {
        $this->missMfcRank = (int)$missMfcRank;
    }

    /**
     * @return int
     */
    public function getAccessLevel()
    {
        return $this->accessLevel;
    }

    /**
     * @param int $accessLevel
     */
    public function setAccessLevel($accessLevel)
    {
        $this->accessLevel = (int)$accessLevel;
    }
}
