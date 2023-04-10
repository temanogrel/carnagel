<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Service;

use Aphrodite\Blocktrail\Options\BlocktrailOptions;
use Aphrodite\Blocktrail\Service\Exception\InitializeBlocktrailWalletException;
use Blocktrail\SDK\BlocktrailSDK;
use Blocktrail\SDK\Wallet;
use Blocktrail\SDK\WalletInterface;
use Exception;

final class BlocktrailService implements BlocktrailServiceInterface
{
    /**
     * @var BlocktrailSDK
     */
    private $blocktrail;

    /**
     * @var BlocktrailOptions
     */
    private $options;

    /**
     * BlocktrailService constructor.
     * @param BlocktrailSDK $blocktrail
     * @param BlocktrailOptions $options
     */
    public function __construct(BlocktrailSDK $blocktrail, BlocktrailOptions $options)
    {
        $this->blocktrail = $blocktrail;
        $this->options    = $options;
    }

    /**
     * {@inheritdoc}
     */
    public function createNewPaymentAddress(): string
    {
        if ($this->options->getTestNet()) {
            return $this->options->getTestAddress();
        }

        return $this->getWallet()->getNewAddress();
    }

    /**
     * {@inheritdoc}
     */
    public function setupAddressTransactionWebhook(string $id, string $address)
    {
        $this->blocktrail->subscribeAddressTransactions(
            $id,
            $address,
            $this->options->getNumberOfConfirmationsToTriggerWebhook()
        );
    }

    /**
     * {@inheritdoc}
     */
    public function sendBitcoin(string $address, float $amount)
    {
        $this->getWallet()->pay([
            $address => BlocktrailSDK::toSatoshi($amount)
        ]);
    }

    /**
     * Get the blocktrail wallet object
     *
     * @return WalletInterface
     *
     * @throws InitializeBlocktrailWalletException
     */
    private function getWallet(): WalletInterface
    {
        try {
            return $this->blocktrail->initWallet([
                'identifier' => $this->options->getWalletId(),
                'passphrase' => $this->options->getWalletPassword(),
            ]);
        } catch (Exception $e) {
            throw new InitializeBlocktrailWalletException($e->getMessage());
        }
    }
}
