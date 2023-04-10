<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Options;

use Zend\Stdlib\AbstractOptions;

final class BlocktrailOptions extends AbstractOptions
{
    /**
     * @var string
     */
    private $walletId;

    /**
     * @var string
     */
    private $apiKey;

    /**
     * @var boolean
     */
    private $testNet = false;

    /**
     * @var string
     */
    private $testAddress;

    /**
     * @var string
     */
    private $apiSecret;

    /**
     * @var string
     */
    private $walletPassword;

    /**
     * @var int
     */
    private $numberOfConfirmationsToTriggerWebhook = 1;

    /**
     * @return string
     */
    public function getWalletId(): string
    {
        return $this->walletId;
    }

    /**
     * @param string $walletId
     */
    public function setWalletId(string $walletId)
    {
        $this->walletId = $walletId;
    }

    /**
     * @return string
     */
    public function getApiKey(): string
    {
        return $this->apiKey;
    }

    /**
     * @param string $apiKey
     */
    public function setApiKey(string $apiKey)
    {
        $this->apiKey = $apiKey;
    }

    /**
     * @return string
     */
    public function getApiSecret(): string
    {
        return $this->apiSecret;
    }

    /**
     * @param string $apiSecret
     */
    public function setApiSecret(string $apiSecret)
    {
        $this->apiSecret = $apiSecret;
    }

    /**
     * @return string
     */
    public function getWalletPassword(): string
    {
        return $this->walletPassword;
    }

    /**
     * @param string $walletPassword
     */
    public function setWalletPassword(string $walletPassword)
    {
        $this->walletPassword = $walletPassword;
    }

    /**
     * @return int
     */
    public function getNumberOfConfirmationsToTriggerWebhook(): int
    {
        return $this->numberOfConfirmationsToTriggerWebhook;
    }

    /**
     * @param int $numberOfConfirmationsToTriggerWebhook
     */
    public function setNumberOfConfirmationsToTriggerWebhook(int $numberOfConfirmationsToTriggerWebhook)
    {
        $this->numberOfConfirmationsToTriggerWebhook = $numberOfConfirmationsToTriggerWebhook;
    }

    /**
     * @return bool
     */
    public function getTestNet(): bool
    {
        return $this->testNet;
    }

    /**
     * @param bool $testNet
     */
    public function setTestNet(bool $testNet)
    {
        $this->testNet = $testNet;
    }

    /**
     * @return string
     */
    public function getTestAddress()
    {
        return $this->testAddress;
    }

    /**
     * @param string $testAddress
     */
    public function setTestAddress(string $testAddress = null)
    {
        $this->testAddress = $testAddress;
    }
}
