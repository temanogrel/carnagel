let environment = {
  apiUri: '',
  websocketUri: '',
  proxyScheme: '',
  cookieDomain: '',
};

switch (document.domain) {
  case 'camtube.dev':
    environment.apiUri = '//api.camtube.dev:9000';
    environment.websocketUri = 'ws://api.camtube.dev:9000/ws';
    environment.cookieDomain = '.camtube.dev';
    break;

  default:
    environment.apiUri = 'https://api.camtube.co';
    environment.websocketUri = 'wss://api.camtube.co/ws';
    environment.cookieDomain = '.camtube.co';
    break;
}

export default environment;
