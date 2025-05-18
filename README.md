# Weather API Application (test task for Genesis Software Engineering School 5.0)

## Quickstart

This is application for subscription on periodic weather updates. To keep it simple, it was written in a monolithic way.

## Requirements
- Docker
- Just

## Environment

In order to provide needed tokens and URI's, setup environment variables:
```env
WEATHER_API_KEY="your weatherapi key"
WEATHER_API_ADDR="http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no"
POSTGRES_DB="subscriptions"
POSTGRES_USER="user"
POSTGRES_PASSWORD="1234"
SENDER_MAIL="your sender email"
SENDGRID_API_KEY="your sendgrid api token"
BASE_URL="localhost:8000"
HTTPS=0
PROD=0
PROD_DB_URL="your prod db"
```

## Docker running

- Building
```console
just build
```

- Running
```console
just run
```

- Dropping running containers
```console
just clean
```

- Starting
```console
just start
```

## Accessing deployed
[Weather Subscription](https://genesiscasestudy-production.up.railway.app/)(Railway) 
