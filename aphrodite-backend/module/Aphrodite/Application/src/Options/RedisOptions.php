<?php
/**
 *
 *
 *  AB
 */

declare(strict_types = 1);

namespace Aphrodite\Application\Options;

use Zend\Stdlib\AbstractOptions;

class RedisOptions extends AbstractOptions
{
    /**
     * @var string
     */
    protected $host = 'localhost';

    /**
     * @var int
     */
    protected $port = 6379;

    /**
     * @return string
     */
    public function getHost()
    {
        return $this->host;
    }

    /**
     * @param string $host
     */
    public function setHost($host)
    {
        $this->host = $host;
    }

    /**
     * @return int
     */
    public function getPort()
    {
        return $this->port;
    }

    /**
     * @param int $port
     */
    public function setPort($port)
    {
        $this->port = $port;
    }
}
