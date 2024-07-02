# GO Expert - Weather API - Cloud Run

You can use the live demo on: https://go-cloud-run-23aifcgpxq-uc.a.run.app/

## Run locally

Setup your your [Weather API](https://www.weatherapi.com/) key on `docker-compose.yml`.

In the project root execute:
```shell
docker compose up
```

## APIs

### GET /weather/{zip_code}

Examples are available at `api/requests.http`

200:
```json
{"temp_c":16,"temp_f":60.8,"temp_k":289.1}
```

422:
```
invalid zipcode
```

404:
```
can not find zipcode
```

## Run tests

go test ./...
