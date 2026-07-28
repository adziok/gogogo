App will be used to manage feature flags.
So I will expose the endpoint to managing the flags as also gettting the flags.
Postgress will be resposible for storing the current state.
Clickhouse will be responsible for analytics. When, who, what ask for.
Maybe also redis as cache will be additional point? // Canceled
Also i would like to add grafana-lgtm and otel support. // Canceled
Also i will use auth0 for auth purpose. For the managing i will use web app, for "flag consumer" machine to machine.
One additional thing, instead of using some buffer table in clickhouse, I will add kafka for that :)
