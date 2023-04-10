<?php
/**
 *
 *
 * 
 */
declare(strict_types=1);

namespace Aphrodite\Logger\Service;

use Aphrodite\Logger\Adapter\AbstractAdapter;
use Aphrodite\Logger\Options\LogHandlerOptions;
use Throwable;
use Zend\Http\Request;
use Zend\Http\Response;
use Zend\Json\Exception\RuntimeException;
use Zend\Json\Json;
use Zend\Mvc\Router\RouteMatch;
use Exception;

final class LogHandlerService implements LogHandlerServiceInterface
{
    /**
     * @var array
     */
    private $adapters = [];

    /**
     * @var LogHandlerOptions
     */
    private $options;

    /**
     * LogHandlerService constructor.
     * @param LogHandlerOptions $options
     * @param array $adapters
     */
    public function __construct(LogHandlerOptions $options, array $adapters)
    {
        $this->options  = $options;
        $this->adapters = $adapters;
    }

    /**
     * {@inheritdoc}
     */
    public function handleException(Throwable $exception, array $context = [])
    {
        if ($exception instanceof Throwable) {
            $data = [
                '@timestamp'  => date(DATE_RFC3339),
                'environment' => $this->options->getEnvironment(),
                'host'        => $this->options->getHost(),
                'message'     => $exception->getMessage(),
                'class'       => get_class($exception),
                'stacktrace'  => $exception->getTraceAsString(),
                'context'     => json_encode($context)
            ];

            $this->recurseException($data, $exception->getPrevious());
        } else {
            $data = [
                '@timestamp'  => date(DATE_RFC3339),
                'environment' => $this->options->getEnvironment(),
                'host'        => $this->options->getHost(),
                'message'     => 'None exception error occurred',
                'class'       => get_class($exception),
                'payload'     => $exception,
                'context'     => json_encode($context)
            ];
        }

        $this->writeData($data, 'errors');
    }

    /**
     * {@inheritdoc}
     */
    public function handleRequestResponse(
        Request $request,
        Response $response,
        RouteMatch $routeMatch = null,
        float $duration = null,
        array $context = []
    ) {
        if ($this->options->isDebug() || $this->isAlwaysLogRoute($routeMatch)) {
            $data = [
                '@timestamp'  => date(DATE_RFC3339),
                'host'        => $this->options->getHost(),
                'environment' => $this->options->getEnvironment(),
                'duration'    => $duration ?? 0.0,
                'request'     => [
                    'headers'    => $request->getHeaders()->toString(),
                    'routeMatch' => $routeMatch ? $routeMatch->getMatchedRouteName() : '',
                    'url'        => $request->getUriString(),
                    'body'       => $this->blurKeywords($request->getContent()),
                    'query'      => $this->blurKeywords($request->getQuery()->toString()),
                    'method'     => $request->getMethod(),
                ],
                'response'    => [
                    'headers'    => $response->getHeaders()->toString(),
                    'body'       => $response->getContent(),
                    'statusCode' => $response->getStatusCode(),
                ],
                'context' => json_encode($context)
            ];

           $this->writeData($data, 'requests');
        }
    }

    /**
     * Write the data to adapters
     *
     * @param array $data
     * @param string $type
     */
    private function writeData(array $data, string $type)
    {
        /* @var AbstractAdapter $adapter */
        foreach ($this->adapters as $adapter) {
            try {
                $adapter->write($data, $type);
            } catch (Exception $e) {
                // To prevent application from crashing
            }
        }
    }

    /**
     * {@inheritdoc}
     */
    public function getAdapters(): array
    {
        return $this->adapters;
    }

    /**
     * Check if string contains keys that should be blurred out
     * This will only check top-level
     *
     * @param string $body
     * @return string
     */
    private function blurKeywords(string $body): string
    {
        try {
            $isJson = true;
            $bodyAsArray = Json::decode($body, true);
        } catch (RuntimeException $e) {
            $isJson = false;
            parse_str($body, $bodyAsArray);
        }

        if (!is_array($bodyAsArray)) {
            return $body;
        }

        foreach ($bodyAsArray as $key => $value) {
            if (in_array($key, $this->options->getBlurredKeys(), false)) {
                $bodyAsArray[$key] = $this->options->getBlurredKeysValue();
            }
        }

        if ($isJson) {
            return Json::encode($bodyAsArray);
        }

        return http_build_query($bodyAsArray);
    }

    /**
     * Check if matched route should always be logged for http request/response
     *
     * @param RouteMatch|null $routeMatch
     * @return bool
     */
    private function isAlwaysLogRoute(RouteMatch $routeMatch = null): bool
    {
        return $routeMatch && in_array($routeMatch->getMatchedRouteName(), $this->options->getAlwaysLogRoutes(), true);
    }

    /**
     * @param $data
     * @param Throwable|null $exception
     */
    private function recurseException(&$data, Throwable $exception = null)
    {
        if ($exception === null) {
            return;
        }

        $data['previous'] = [
            'class'      => get_class($exception),
            'message'    => $exception->getMessage(),
            'stacktrace' => $exception->getTraceAsString(),
        ];

        $this->recurseException($data['previous'], $exception->getPrevious());
    }
}
