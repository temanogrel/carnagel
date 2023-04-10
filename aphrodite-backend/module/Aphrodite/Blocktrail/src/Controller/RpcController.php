<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Controller;

use Aphrodite\Blocktrail\BlocktrailPermissions;
use Aphrodite\Blocktrail\InputFilter\CreateAddressInputFilter;
use Aphrodite\Blocktrail\InputFilter\SendBitcoinInputFilter;
use Aphrodite\Blocktrail\Service\BlocktrailServiceInterface;
use Aphrodite\Blocktrail\Service\Exception\InitializeBlocktrailWalletException;
use Zend\Http\Exception\InvalidArgumentException;
use Zend\Http\Request;
use Zend\Http\Response;
use Zend\Mvc\Controller\AbstractActionController;
use Zend\View\Model\JsonModel;
use ZfrRest\Http\Exception\Client\ForbiddenException;
use ZfrRest\Http\Exception\Client\MethodNotAllowedException;
use ZfrRest\Http\Exception\Server\InternalServerErrorException;

/**
 * Class RpcController
 *
 * @method bool isGranted($permission, $context = null)
 * @method array validateIncomingData($inputFilter)
 * @method Request getRequest
 * @method Response getResponse
 */
final class RpcController extends AbstractActionController
{
    /**
     * @var BlocktrailServiceInterface
     */
    private $service;

    /**
     * RpcController constructor.
     * @param BlocktrailServiceInterface $service
     */
    public function __construct(BlocktrailServiceInterface $service)
    {
        $this->service = $service;
    }

    /**
     * Create a new address on the configured blocktrail wallet
     *
     * @return JsonModel
     *
     * @throws InternalServerErrorException
     * @throws ForbiddenException
     * @throws MethodNotAllowedException
     */
    public function createAddressAction(): JsonModel
    {
        if (!$this->getRequest()->isPost()) {
            throw new MethodNotAllowedException();
        }

        if (!$this->isGranted(BlocktrailPermissions::CREATE_ADDRESS)) {
            throw new ForbiddenException('Does not have permission to create new bitcoin address');
        }

        $values = $this->validateIncomingData(CreateAddressInputFilter::class);

        try {
            $address = $this->service->createNewPaymentAddress();
        } catch (InitializeBlocktrailWalletException $e) {
            throw new InternalServerErrorException('Unable to initialize blocktrail wallet', [
                'errors' => [
                    'failedToInitializeWallet' => $e->getMessage(),
                ],
            ]);
        }

        $this->service->setupAddressTransactionWebhook($values['webhookId'], $address);

        return new JsonModel([
            'address' => $address,
        ]);
    }

    /**
     * Send a given bitcoin amount to a given address
     *
     * @return Response
     *
     * @throws InvalidArgumentException
     * @throws ForbiddenException
     * @throws MethodNotAllowedException
     */
    public function sendBitcoinAction(): Response
    {
        if (!$this->getRequest()->isPost()) {
            throw new MethodNotAllowedException();
        }

        if (!$this->isGranted(BlocktrailPermissions::SEND_BITCOIN)) {
            throw new ForbiddenException('Does not have permission to send bitcoin');
        }

        $values = $this->validateIncomingData(SendBitcoinInputFilter::class);

        $this->service->sendBitcoin($values['address'], $values['bitcoinAmount']);

        return $this->getResponse()->setStatusCode(Response::STATUS_CODE_204);
    }
}
