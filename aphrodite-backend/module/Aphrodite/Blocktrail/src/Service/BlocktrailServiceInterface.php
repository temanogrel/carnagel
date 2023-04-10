<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Service;

use Aphrodite\Blocktrail\Service\Exception\InitializeBlocktrailWalletException;

interface BlocktrailServiceInterface
{
    /**
     * Create a new payment address on the blocktrail wallet
     *
     * @return string
     *
     * @throws InitializeBlocktrailWalletException
     */
    public function createNewPaymentAddress(): string;

    /**
     * Set up callback to infinity when a new transaction is received on address
     *
     * @param string $id
     * @param string $address
     *
     * @return void
     */
    public function setupAddressTransactionWebhook(string $id, string $address);

    /**
     * Send given amount of bitcoin to given address
     *
     * @param string $address
     * @param float $amount
     *
     * @return void
     */
    public function sendBitcoin(string $address, float $amount);
}
