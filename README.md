# Checkers

## Docker Commands

- `docker build . -t checkers_i`
- `docker run --rm -it checkers_i ignite version`
- `docker run -it checkers_i ignite scaffold chain github.com/alice/checkers`

- Start chain in a terminal

```bash
docker create --name checkers -i \
    -v $(pwd):/checkers -w /checkers \
    -p 1317:1317 -p 3000:3000 -p 4500:4500 -p 5000:5000 -p 26657:26657 \
    checkers_i
docker start checkers
docker exec -it checkers ignite chain serve
```

- In another terminal

```bash
docker exec -it checkers bash -c "checkersd status 2>&1 | jq"  # checkers checkersd
docker exec -it checkers checkersd --help
docker exec -it checkers checkersd status --help
docker exec -it checkers checkersd query --help
```

- For the frontend

```bash
docker exec -it checkers bash -c "cd vue && npm install"
docker exec -it checkers bash -c "cd vue && npm run dev -- --host"
```

- Add counter to store in the blockchain

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold single systemInfo nextId:uint \
    --module checkers \
    --no-message
```

- Scaffold map using the StoredGame name

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold map storedGame board turn black red \
    --index index \
    --module checkers \
    --no-message
```

- Update the types/genesis/go automatically after updating the genesis.proto

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite generate proto-go
```

- Run the generated tests

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go test github.com/BenWolfaardt/checkers/x/checkers/keeper
```

- Run the test we added

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go test github.com/BenWolfaardt/checkers/x/checkers/types
```
