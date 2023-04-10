<?php
/**
 *
 *
 *
 */
declare(strict_types=1);

namespace Aphrodite\Logger;

use Aphrodite\Logger\Listener\ErrorListener;
use Aphrodite\Logger\Listener\RequestResponseDataListener;
use Zend\Loader\StandardAutoloader;
use Zend\ModuleManager\Feature\AutoloaderProviderInterface;
use Zend\ModuleManager\Feature\ConfigProviderInterface;
use Zend\Mvc\MvcEvent;

class Module implements ConfigProviderInterface, AutoloaderProviderInterface
{
    /**
     * @inheritDoc
     */
    public function getAutoloaderConfig(): array
    {
        return [
            StandardAutoloader::class => [
                StandardAutoloader::LOAD_NS => [
                    __NAMESPACE__ => __DIR__ . '/src',
                ],
            ],
        ];
    }

    public function onBootstrap(MvcEvent $e)
    {
        $eventManager = $e->getApplication()->getEventManager();
        $sm           = $e->getApplication()->getServiceManager();

        /* @var RequestResponseDataListener $requestResponseListener */
        $requestResponseListener = $sm->get(RequestResponseDataListener::class);
        $requestResponseListener->attach($eventManager);

        /* @var ErrorListener $errorListener */
        $errorListener = $sm->get(ErrorListener::class);
        $errorListener->attach($eventManager);
    }

    /**
     * @inheritDoc
     */
    public function getConfig()
    {
        return include __DIR__ . '/config/module.config.php';
    }
}
