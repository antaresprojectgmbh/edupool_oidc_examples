# EduPool OAuth2 Client in PHP

This repository contains a PHP script that demonstrates how to implement a basic OAuth2 client using the `league/oauth2-client` library, specifically for the EduPool Identity Provider. The client authenticates users via the EduPool OAuth2 server, fetches their profile information, and displays it in the browser.

## Requirements

To run this script, you need:

-   PHP 7.4 or higher
-   Composer (to install dependencies)


## Usage
**Configure the script:**

Open the script in a text editor and replace `yourClientID` and `yourClientSecret` with your actual client ID and secret provided by EduPool.

**Start a PHP built-in web server:**

`php -S 127.0.0.1:9010` 

The script will be accessible at [http://127.0.0.1:9010](http://127.0.0.1:9010/).
1.  Open the URL [http://127.0.0.1:9010](http://127.0.0.1:9010/) in your browser.
2.  Click the "Log In" link.
3.  You will be redirected to the EduPool OIDC server. Log in with your credentials.
4.  After successful authentication, you will be redirected back to the client application, and your user information will be displayed on the screen.