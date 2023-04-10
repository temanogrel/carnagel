<?php
/**
 *
 * 
 */

declare(strict_types=1);

namespace Aphrodite\Logger\Options;

use Zend\Stdlib\AbstractOptions;

final class LogHandlerOptions extends AbstractOptions
{
    /**
     * @var bool
     */
    protected $debug = false;

    /**
     * @var array
     */
    protected $alwaysLogRoutes = [];

    /**
     * @var array
     */
    protected $blurredKeys = ['password'];

    /**
     * @var string
     */
    protected $blurredKeysValue = '***FILTERED***';

    /**
     * @var array
     */
    protected $adapters = [];

    /**
     * @var string
     */
    protected $environment = 'dev';

    /**
     * @var string
     */
    protected $host;

    /**
     * @return bool
     */
    public function isDebug(): bool
    {
        return $this->debug;
    }

    /**
     * @param bool $debug
     */
    public function setDebug(bool $debug)
    {
        $this->debug = $debug;
    }

    /**
     * @return array
     */
    public function getAlwaysLogRoutes(): array
    {
        return $this->alwaysLogRoutes;
    }

    /**
     * @param array $alwaysLogRoutes
     */
    public function setAlwaysLogRoutes(array $alwaysLogRoutes)
    {
        $this->alwaysLogRoutes = $alwaysLogRoutes;
    }

    /**
     * @return array
     */
    public function getAdapters(): array
    {
        return $this->adapters;
    }

    /**
     * @param mixed $adapters
     */
    public function setAdapters($adapters)
    {
        $this->adapters = $adapters;
    }

    /**
     * @return string
     */
    public function getEnvironment(): string
    {
        return $this->environment;
    }

    /**
     * @param string $environment
     */
    public function setEnvironment(string $environment)
    {
        $this->environment = $environment;
    }

    /**
     * @return string
     */
    public function getHost(): string
    {
        return $this->host;
    }

    /**
     * @param string $host
     */
    public function setHost(string $host)
    {
        $this->host = $host;
    }

    /**
     * @return array
     */
    public function getBlurredKeys(): array
    {
        return $this->blurredKeys;
    }

    /**
     * @param array $blurredKeys
     */
    public function setBlurredKeys(array $blurredKeys)
    {
        $this->blurredKeys = $blurredKeys;
    }

    /**
     * @return string
     */
    public function getBlurredKeysValue(): string
    {
        return $this->blurredKeysValue;
    }

    /**
     * @param string $blurredKeysValue
     */
    public function setBlurredKeysValue(string $blurredKeysValue)
    {
        $this->blurredKeysValue = $blurredKeysValue;
    }
}
