# Checkers

## Tutorial / Commands

### Rebuild Your Cosmos Chain With Ignite

#### Ignite CLI

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

#### Store Object

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

#### Create and Save a Game Properly

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

#### Create Custom Messages

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

#### Add a Way to Make a Move

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

#### Emit Game Information

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

#### Record the Game Winner

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

### Continue Developing Your Chain

#### Keep an Up-To-Date Game Deadline

- Update the proto after adding `deadline` to StoredGame

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite generate proto-go
```

- Confirm the project still compiles

```bash
docker run --rm -it \
    --name checkers \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite chain build
```

- Test to see new field

> If you still have state in your blockchain a migration will be necessary (we learn about that later)

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1
```

- Try making a new move to observe state

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 2 2 3 --from $alice
```

- Observe new move's `storedGame`

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1
```

- Or you could restart your blockchain deleting state

```bash
docker exec -it checkers ignite chain serve --reset-once
```

- Start a new game

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice
```

- Start a new game

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # storedGame:
    # black: cosmos1ujnssfn87l4zu0m6f2dvrnn5t270gc3nftvml2
    # board: '*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*'
    # deadline: 2024-06-07 14:27:40.235468194 +0000 UTC             <------ New field
    # index: "1"
    # red: cosmos1hpks8k8azlwdja9xnhw56gtrcpcdm4l0pnc9gq
    # turn: b
    # winner: '*'
```

> Observe the new `deadline` field in the `storedGame` output

#### Keep Track Of How Many Moves Have Been Played

- Add `moveCount` to `stored_game.proto` and update proto definition

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite generate proto-go
```

- View newly added field in `storedGame`

> If your blockchain state already contains games then they are missing the new field they'll thus get the default value `0`:

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # moveCount: "0"
```

> Our blockchain state is technically broken, hence a migration being necessary, something we will learn about later.

- Reset chain and create new game to test functionality of `moveCount`

```bash
docker exec -it checkers ignite chain serve --reset-once
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 2 2 3 --from $alice
docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # moveCount: "1"
```

#### Put Your Games in Order

- Update protobuf files are FIFO additonal info added to them.

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite generate proto-go
```

- Run updated tests

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go test github.com/alice/checkers/x/checkers/keeper
```

- Reset and restart chain

```bash
docker exec -it checkers \
    ignite chain serve --reset-once
```

- Is the genesis FIFO information correctly saved?

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info

    # SystemInfo:
    #   fifoHeadIndex: "-1"
    #   fifoTailIndex: "-1"
    #   nextId: "1"
```

- If you create a game, is the game as expected?

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $bob
docker exec -it checkers \
    checkersd query checkers show-system-info
    
    # SystemInfo:
    #   fifoHeadIndex: "1"
    #   fifoTailIndex: "1"
    #   nextId: "2"
```

- What about the information saved in the game?

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # storedGame:
    #   afterIndex: "-1"
    #   beforeIndex: "-1"
```

- And if you create another game?

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $bob
docker exec -it checkers \
    checkersd query checkers show-system-info

    # SystemInfo:
    #   fifoHeadIndex: "1"
    #   fifoTailIndex: "2"
    #   nextId: "3"
```

- Did the games also store the correct values? For the first game:

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # afterIndex: "2" # The second game you created
    # beforeIndex: "-1" # No game
```

- For the second game:

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 2

    # afterIndex: "-1" # No game
    # beforeIndex: "1" # The first game you created
```

> The FIFO in effect has the game IDs [1, 2]. 

- Add a third game, and confirm that your FIFO is [1, 2, 3].

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $bob
$ docker exec -it checkers \
    checkersd query checkers show-system-info
```

- What happens if Alice plays a move in game 2, the game in the middle?

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 2 1 2 2 3 --from $alice
docker exec -it checkers \
    checkersd query checkers show-system-info
```

- Is game 3 in the middle now?

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 3
```

#### Auto-Expiring Games

- Reset and restart the chain

```bash
docker exec -it checkers ignite chain serve --reset-once
```

> Export your aliases again

- Create three games one minute apart. Have Alice play the middle one, and both Alice and Bob play the last one:

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice```

> Wait one minute

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $bob
docker exec -it checkers \
    checkersd tx checkers play-move 2 1 2 2 3 --from $alice
```

> Wait another minute

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob --from $alice
docker exec -it checkers \
    checkersd tx checkers play-move 3 1 2 2 3 --from $alice
docker exec -it checkers \
    checkersd tx checkers play-move 3 0 5 1 4 --from $bob
```

- With three games in, confirm that you see them all

```bash
docker exec -it checkers \
    checkersd query checkers list-stored-game
```

- List them again after two, three, four, and five minutes. You should see games 1 and 2 disappear, and game 3 being forfeited by Alice, i.e. red Bob wins:

```bash
docker exec -it checkers \
    bash -c "checkersd query checkers show-stored-game 3 --output json | jq '.storedGame.winner'"

    # "r"
```

- Confirm that the FIFO no longer references the removed games nor the forfeited game

```bash
docker exec -it checkers \
    checkersd query checkers show-system-info

    # SystemInfo:
        # fifoHeadIndex: "-1"
        # fifoTailIndex: "-1"
        # nextId: "4"
```

- 

```bash

```


## Notes

- Sometimes full docker use other times quick exec
- Oftentimes bob and alice need to be exported when reseting chain
- Dev container means you can do all locally
- Can't run checkersd commands with spaces in front or you get: `Error: rpc error: code = NotFound desc = rpc error: code = NotFound desc = not found: key not found`
