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

- Start and reset the chain

```bash
docker run --rm -it \
    --name checkers \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite chain serve --reset-once
```

- Query the CLI

```bash
docker exec -it checkers \
    checkersd query checkers --help
```

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info
```

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info --help
```

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info --output json
```

```bash
docker exec -it checkers \
    checkersd query checkers list-stored-game
```

> Remember how you wrote --no-message? That was to not create messages or transactions, which would directly update your checkers storage. Soft-confirm there are no commands available:

```bash
docker exec -it checkers \
    checkersd tx checkers --help
```

- Create custom messages

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold message createGame black red \
    --module checkers \
    --response gameIndex
```

- Run tests on new test

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go test github.com/alice/checkers/x/checkers/keeper
```

- Interact with the CLI

```bash
docker exec -it checkers \
    checkersd tx checkers --help
```

```bash
docker exec -it checkers \
    checkersd tx checkers create-game --help
```

- Export Alice and Bob's addresses

```bash
export alice=$(docker exec checkers checkersd keys show alice -a)
export bob=$(docker exec checkers checkersd keys show bob -a)
```

- How much gas is needed?
  - Note this isn't working yet.

```bash
 docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice --dry-run
```

- Keep gas set to auto

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice --gas auto
```

- If you are curious, the .events.attributes are encoded in Base64:

```bash
docker exec -it checkers \
    bash -c "echo YWN0aW9u | base64 -d"
> action%   
docker exec -it checkers \
    bash -c "echo Y3JlYXRlX2dhbWU= | base64 -d"
> create_game%   
```

- Check to see if anything changed

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info
```

- Check whether any game was created:

```bash
docker exec -it checkers \
    checkersd query checkers list-stored-game
```
