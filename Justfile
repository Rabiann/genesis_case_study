run:
    docker compose -f compose.yaml up

build:
    docker compose -f compose.yaml build

clean:
    docker compose down -v

start:
    docker compose down -v
    docker compose -f compose.yaml build
    docker compose -f compose.yaml up