<?php
/**
 *
 *
 */

declare(strict_types=1);

namespace Aphrodite\Blocktrail\Factory\Controller;

use Aphrodite\Blocktrail\Controller\RpcController;
use Aphrodite\Blocktrail\Service\BlocktrailService;
use Zend\Mvc\Controller\ControllerManager;

final class RpcControllerFactory
{
    /**
     * @param ControllerManager $controllerManager
     * @return RpcController
     */
    public function __invoke(ControllerManager $controllerManager): RpcController
    {
        $container = $controllerManager->getServiceLocator();

        return new RpcController(
            $container->get(BlocktrailService::class)
        );
    }
}
