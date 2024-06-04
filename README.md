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

export alice=$(checkersd keys show alice -a)
export bob=$(checkersd keys show bob -a)
```

- How much gas is needed?
  - Note this isn't working yet.

```bash
 docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice --dry-run
```

> `--from`: Name or address of private key with which to sign

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

### Create and Save a Game Properly

- Create a game

> Initiate game

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice --gas auto
```

- Show system info

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info
```

- List all stored games

```bash
docker exec -it checkers \
    checkersd query checkers list-stored-game
```

- Show the new game alone:

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1
```

- Pretty formatted checkers board

> View checkers board

```bash
docker exec -it checkers \
    bash -c "checkersd query checkers show-stored-game 1 --output json | jq \".storedGame.board\" | sed 's/\"//g' | sed 's/|/\n/g'"
```

### Add a Way to Make a Move

- Ignite command

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold message playMove gameIndex fromX:uint fromY:uint toX:uint toY:uint \
    --module checkers \
    --response capturedX:int,capturedY:int,winner
```

- Start chain

```bash
docker exec -it checkers ignite chain serve
```

- Bob plays out of turn

```bash
docker exec -it checkers \
    checkersd tx checkers play-move --help
```

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 0 5 1 4 --from $bob
    #                               ^ ^ ^ ^ ^
    #                               | | | | To Y
    #                               | | | To X
    #                               | | From Y
    #                               | From X
    #                               Game id
```

- Query txhas

```bash
docker exec -it checkers \
    checkersd query tx <txhas>
```

- Alice plays wrong move by taking a piece on the side and moving it just outside the board:

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 7 2 8 3 --from $alice
```

> `Error: toX out of range (8): position index is invalid` - The transaction never went into the mem pool. This mistake did not cost Alice any gas.

- Alice tries to move to an occupied place

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 0 0 1 --from $alice
```

- Alice plays correctly

> Next move, Alice's first turn

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 2 2 3 --from $alice
```

- Confirm the move visually

> See updated board

```bash
docker exec -it checkers \
    bash -c "checkersd query checkers show-stored-game 1 --output json | jq \".storedGame.board\" | sed 's/\"//g' | sed 's/|/\n/g'"

    # *b*b*b*b
    # b*b*b*b*
    # ***b*b*b
    # **b*****     <--- Here
    # ********
    # r*r*r*r*
    # *r*r*r*r
    # r*r*r*r*
```

### Emit Game Information

> Bob's turn

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 0 5 1 4 --from $bob
```

- View event data better

```bash
docker exec -it checkers \
    checkersd query tx <txhash>
```

- Bob's pawn is now ready to be captured by Alice

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1 --output json | jq ".storedGame.board" | sed 's/"//g' | sed 's/|/\n/g'

    # *b*b*b*b
    # b*b*b*b*
    # ***b*b*b
    # **b*****
    # *r******    <-- Ready to be captured
    # **r*r*r*
    # *r*r*r*r
    # r*r*r*r*
```

> The rules of the game included in this project mandate that the player captures a piece when possible.

- Alice captures the piece

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 2 3 0 5 --from $alice
```

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1 --output json | jq ".storedGame.board" | sed 's/"//g' | sed 's/|/\n/g'

    # *b*b*b*b
    # b*b*b*b*
    # ***b*b*b
    # ********
    # ********
    # b*r*r*r*
    # *r*r*r*r
    # r*r*r*r*
```

### Record the Game Winner

- Ignite command

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite generate proto-go
```

- Reset the chain as we have updated the default value of winner from `.Winner == ""` to `rules.PieceStrings[rules.NO_PLAYER]`.

```bash
docker run --rm -it \
    --name checkers \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite chain serve --reset-once
```

- Confirm there is no winner for a newly created game

> Note: need to re `export` Bob and Alice

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice

docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # ...
    #   winner: '*'
    # ...
```

- Confirm there is no winner after first move

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 2 2 3 --from $alice

docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # ...
    #   winner: '*'
    # ...
```
