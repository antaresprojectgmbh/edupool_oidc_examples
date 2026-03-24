<?php

require_once 'vendor/autoload.php';

use League\OAuth2\Client\Provider\GenericProvider;
use League\OAuth2\Client\Provider\Exception\IdentityProviderException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

session_start();

// the urls can be found here: https://oidc.edupool.cloud/.well-known/openid-configuration

$provider = new GenericProvider([
    'clientId' => 'yourClientID',
    'clientSecret' => 'yourClientSecret',
    'redirectUri' => 'http://127.0.0.1:9010/callback',
    'urlAuthorize' => 'https://oidc.edupool.cloud/oauth2/auth',
    'urlAccessToken' => 'https://oidc.edupool.cloud/oauth2/token',
    'urlResourceOwnerDetails' => 'https://oidc.edupool.cloud/userinfo',
    'scopes' => 'openid offline profile antares.context'
]);

$dispatcher = FastRoute\simpleDispatcher(function (FastRoute\RouteCollector $r) use ($provider) {
    $r->addRoute('GET', '/', function (ServerRequestInterface $request, ResponseInterface $response) {
        $response->getBody()->write('<html><body>'
            . '<a href="/login">Log In (default)</a><br>'
            . '<a href="/login?land=nrw">Log In (optional shortcut: land=nrw)</a><br>'
            . '<a href="/login?context=hh/HH">Log In (optional shortcut: context=hh/HH)</a>'
            . '</body></html>');
        return $response;
    });

    $r->addRoute('GET', '/login', function (ServerRequestInterface $request, ResponseInterface $response) use ($provider) {
        $queryParams = $request->getQueryParams();
        $authOptions = [];

        // Optional shortcut params for Hub preselection.
        // If omitted, default OIDC login behavior is unchanged.
        if (!empty($queryParams['land'])) {
            $authOptions['land'] = $queryParams['land'];
        }
        if (!empty($queryParams['context'])) {
            $authOptions['context'] = $queryParams['context'];
        }

        $authorizationUrl = $provider->getAuthorizationUrl($authOptions);
        $_SESSION['oauth2state'] = $provider->getState();
        header('Location: ' . $authorizationUrl);
        exit();
    });

    $r->addRoute('GET', '/callback', function (ServerRequestInterface $request, ResponseInterface $response) use ($provider) {
        $queryParams = $request->getQueryParams();
        if (!isset($queryParams['state']) || $queryParams['state'] !== $_SESSION['oauth2state']) {
            unset($_SESSION['oauth2state']);
            exit('Invalid state');
        }

        try {
            $accessToken = $provider->getAccessToken('authorization_code', [
                'code' => $queryParams['code']
            ]);
            $resourceOwner = $provider->getResourceOwner($accessToken);
            echo '<pre>';
            var_dump($resourceOwner->toArray());
            echo '</pre>';
        } catch (IdentityProviderException $e) {
            exit($e->getMessage());
        }
    });
});

$request = Laminas\Diactoros\ServerRequestFactory::fromGlobals();
$response = new Laminas\Diactoros\Response();

$routeInfo = $dispatcher->dispatch($request->getMethod(), $request->getUri()->getPath());

switch ($routeInfo[0]) {
    case FastRoute\Dispatcher::FOUND:
        $response = $routeInfo[1]($request, $response);
        break;
    default:
        $response = $response->withStatus(404);
}

$emitter = new Laminas\Diactoros\Response\SapiEmitter();
$emitter->emit($response);
