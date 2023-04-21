const express = require('express');
const { Issuer, generators } = require('openid-client');
const app = express();
const port = 9010;

const clientMetadata = {
  client_id: 'yourClientID',
  client_secret: 'yourClientSecret',
  redirect_uris: ['http://127.0.0.1:9010/callback'],
  response_types: ['code'],
  token_endpoint_auth_method: 'client_secret_basic',
};

const oauthStateString = 'peofjkwfjwieufhiu';

(async () => {
  const issuer = await Issuer.discover('https://oidc.edupool.cloud/');
  const client = new issuer.Client(clientMetadata);
  const scopes = 'openid offline profile antares.context';

  app.get('/', (req, res) => {
    const htmlIndex = `
      <html>
      <body>
        <a href="/login">Log In</a>
      </body>
      </html>
    `;
    res.send(htmlIndex);
  });

  app.get('/login', (req, res) => {
    const authUrl = client.authorizationUrl({
      scope: scopes,
      state: oauthStateString,
    });
    res.redirect(authUrl);
  });

  app.get('/callback', async (req, res) => {
    const params = client.callbackParams(req);
    const tokenSet = await client.callback('http://127.0.0.1:9010/callback', params, { state: oauthStateString });
    const idToken = tokenSet.id_token;
    const claims = tokenSet.claims();

    console.log(JSON.stringify(claims, null, 2));
    res.send('Logged in! Check the console for claims.');
  });

  app.listen(port, () => {
    console.log(`Server listening at http://127.0.0.1:${port}`);
  });
})();