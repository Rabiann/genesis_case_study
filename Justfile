run:
    go run ./src

test:
    docker compose down -v
    docker compose -f compose.yaml build
    docker compose -f compose.yaml up