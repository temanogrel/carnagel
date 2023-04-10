<?php
/**
 * 
 *
 *  AB
 */

namespace Aphrodite\User\Rbac;

use Aphrodite\User\Rbac\Identity\ServerIdentity;
use Aphrodite\User\Rbac\Identity\UltronIdentity;
use Zend\Authentication\AuthenticationService;
use Zend\Crypt\Utils as CryptUtils;
use Zend\Http\Request as HttpRequest;
use Zend\Stdlib\RequestInterface;
use ZfcRbac\Identity\IdentityInterface;
use ZfcRbac\Identity\IdentityProviderInterface;

class IdentityProvider implements IdentityProviderInterface
{
    /**
     * @var RequestInterface
     */
    private $request;

    /**
     * @var AuthenticationService
     */
    private $authentication;

    /**
     * @var string
     */
    private $serverAccessToken;

    /**
     * @param AuthenticationService $authentication
     * @param RequestInterface      $request
     * @param string                $serverAccessToken
     */
    public function __construct(
        AuthenticationService $authentication,
        RequestInterface $request,
        string $serverAccessToken
    ) {
        $this->request           = $request;
        $this->authentication    = $authentication;
        $this->serverAccessToken = $serverAccessToken;
    }

    /**
     * @return null
     */
    private function tryStaticIdentityVerification()
    {
        if (!$this->request instanceof HttpRequest) {
            return null;
        }

        $authorization = $this->request->getHeader('Authorization');
        if (!$authorization) {
            return null;
        }

        list ($type, $token) = explode(' ', $authorization->getFieldValue());

        if ($type === 'server' && CryptUtils::compareStrings($token, $this->serverAccessToken)) {
            return new ServerIdentity();
        }

        return null;
    }

    /**
     * Get the identity
     *
     * @return IdentityInterface|null
     */
    public function getIdentity()
    {
        $identity = $this->tryStaticIdentityVerification();
        if ($identity instanceof IdentityInterface) {
            return $identity;
        }

        return $this->authentication->getIdentity();
    }
}
