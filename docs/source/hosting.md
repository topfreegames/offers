Hosting Offers
==============

## Docker

Running Offers with docker is rather simple. Our docker container image comes bundled with the API binary. All you need to do is load balance all the containers and you're good to go. The API runs at port `8888` in the docker image.

Offers uses PostgreSQL to store offers information. The container takes environment variables to specify this connection:

* `OFFERS_POSTGRES_HOST` - PostgreSQL host to connect to;
* `OFFERS_POSTGRES_PORT` - PostgreSQL port to connect to;
* `OFFERS_POSTGRES_DBNAME` - PostgreSQL database to connect to;
* `OFFERS_POSTGRES_PASSWORD` - Password of the PostgreSQL Server to connect to;
* `OFFERS_POSTGRES_USER` - PostgreSQL user;

Offers uses basic auth to restrict access to routes that are not used directly by a client consuming the offers.

* `OFFERS_BASICAUTH_USERNAME` - Basic Auth user;
* `OFFERS_BASICAUTH_PASSWORD` - Basic Auth password;

When a client requests the available offers the API returns a `max-age` header. The cache TTL (in seconds) can be defined using the following variable:

* `OFFERS_CACHE_MAXAGESECONDS` - Max age in seconds;

Other than that, there are a couple more configurations you can pass using environment variables:

* `OFFERS_NEWRELIC_KEY` - If you have a [New Relic](https://newrelic.com/) account, you can use this variable to specify your API Key to populate data with New Relic API;
* `OFFERS_NEWRELIC_APP` - Name of the NewRelic app ;
* `OFFERS_SENTRY_URL` - If you have a [sentry server](https://docs.getsentry.com/hosted/) you can use this variable to specify your project's URL to send errors to.
