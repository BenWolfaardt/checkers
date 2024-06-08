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

#### Let Players Set a Wager

- Update the protobuf after the addition of `wager` to `StoredGame` and `MsgCreateGame`

```bash
docker run --rm -it \
    ignite generate proto-go
```

- Reset and restart chain

```bash
docker exec -it checkers \
    ignite chain server --reset-once
```

- Check initial balances

```bash
docker exec -it checkers \
    checkersd query bank balances $alice
docker exec -it checkers \
    checkersd query bank balances $bob

    # balances:
    # - amount: "100000000"
        # denom: stake
    # - amount: "20000"
        # denom: token
    # pagination:
        # next_key: null
        # total: "0"

    # balances:
    # - amount: "100000000"
        # denom: stake
    # - amount: "10000"
        # denom: token
    # pagination:
        # next_key: null
        # total: "0"
```

- Create a game with a wager

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob 1000000 --from $alice
```

- Was the game stored correctly?

```bash
docker exec -it checkers \
    checkersd query checkers show-stored-game 1

    # storedGame:
    #   ...
    #   wager: "1000000"
```

#### Handle Wager Payments

- Restart chain


// TODO confirm if you get errors?

```bash
docker exec -it checkers \
    ignite chain serve --reset-once
```

- Update your `Dockerfile` with the addition of `gomock` and `make`

- Prepare mocks

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    mockgen \
    -source=x/checkers/types/expected_keepers.go \
    -package testutil \
    -destination=x/checkers/testutil/expected_keepers_mocks.go
```

- Create Makefile to automate the above as it will be used frequently

- Run Makefile

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    make mock-expected-keepers
```

- After writing all tests and updating others make sure all tests work

```bash
docker exec -it checkers \
    go test github.com/BenWolfaardt/checkers/x/checkers/keeper
```

- Check balances

```bash
docker exec -it checkers \
    checkersd query bank balances $alice
docker exec -it checkers \
    checkersd query bank balances $bob
```

- Create a game on which the wager will be refunded because the player playing red did not join

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob 1000000 --from $alice
```

- Confirm that the balances of both Alice and Bob are unchanged - as they have not played yet.

- Have Alice play

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 2 2 3 --from $alice
```

- Confirm that Alice has paid her wager

```bash
docker exec -it checkers \
    checkersd query bank balances $alice

    # balances:
    # - amount: "99000000" # <- 1,000,000 fewer
    # denom: stake
    # - amount: "20000"
    # denom: token
    # pagination:
    # next_key: null
    # total: "0"
```

- Wait 5 minutes for the game to expire and check again

```bash
docker exec -it checkers \
    checkersd query bank balances $alice

    # balances:
    # - amount: "100000000" # <- 1,000,000 are back
    #   denom: stake
    # - amount: "20000"
    #   denom: token
    # pagination:
    #   next_key: null
    #   total: "0"
```

- Now create a game in which both players only play once each, i.e. where the player playing black forfeits

```bash
docker exec -it checkers \
    checkersd tx checkers create-game $alice $bob 1000000 --from $alice
docker exec -it checkers \
    checkersd tx checkers play-move 2 1 2 2 3 --from $alice
docker exec -it checkers \
    checkersd tx checkers play-move 2 0 5 1 4 --from $bob
```

- Confirm that both Alice and Bob paid their wagers. Wait 5 minutes for the game to expire and check again

```bash
docker exec -it checkers \
    checkersd query bank balances $alice
docker exec -it checkers \
    checkersd query bank balances $bob

    # balances:
    # - amount: "99000000" # <- her 1,000,000 are gone for good
    # denom: stake
    # ...
    # balances:
    # - amount: "101000000" # <- 1,000,000 more than at the beginning
    # denom: stake
    # ...
```

#### Integration Tests

- Add integration tests

```bash
docker exec -it checkers \
    go test github.com/BenWolfaardt/checkers/tests/integration/checker
s/keeper
```

#### Play with Cross-Chain Tokens

- Recompile the protobufs after adding `denom` to `MsgCreateGame` and `StoredGame`

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite generate proto-go
```  

- Update all tests accounting for the `denom`

```bash
go test github.com/BenWolfaardt/checkers/x/checkers/keeper
```

- Reset and restart chain

> Export addresses

```bash
docker exec -it checkers \
    checkersd query bank balances $alice

    # balances:
    # - amount: "100000000"
    #   denom: stake
    # - amount: "20000"
    #   denom: token
    # pagination:
    #   next_key: null
    #   total: "0"
```
   
- You can make use of this other token to create a new game

```bash
docker exec -it checkers \
    checkersd tx checkers create-game \
    $alice $bob 1 token \
    --from $alice

    # ...
    # - key: wager
    #   value: "1"
    # - key: denom
    #   value: token
    # ...
```
    
- Have Alice play once

```bash
docker exec -it checkers \
    checkersd tx checkers play-move 1 1 2 2 3 --from $alice

    # - attributes:
    #   - key: recipient
    #     value: cosmos16xx0e457hm8mywdhxtmrar9t09z0mqt9x7srm3
    #   - key: sender
    #     value: cosmos1d2ftstec7v6glkmtz820lnu5x929cm9vv2m99d
    #   - key: amount
    #     value: 1token
    #   type: transfer
```
    
- Check to see if Alice has been charged the wager?

```bash
docker exec -it checkers \
    checkersd query bank balances $alice

    # balances:
    # - amount: "100000000"
    #   denom: stake
    # - amount: "19999"
    #   denom: token
    # pagination:
    #   next_key: null
    #   total: "0"
```
 
- Now check the checkers module's balance

> `cosmos16xx0e457hm8mywdhxtmrar9t09z0mqt9x7srm3` is the checkers module's address.

```bash
docker exec -it checkers \
    checkersd query bank balances cosmos16xx0e457hm8mywdhxtmrar9t09z0mqt9x7srm3

    # balances:
    # - amount: "1"
    #   denom: token
    # pagination:
    #   next_key: null
    #   total: "0"
```

### CosmJS for Your Chain

> Skipped this section for now

### From Code to MVP to Production and Migrations

#### Simulate Production in Docker

- Update Makefile then run

```bash
docker exec -it checkers \
    make build-with-checksum
```

- Create `Dockerfile-checkersd-debian` and build the image

```bash
docker build -f prod-sim/Dockerfile-checkersd-debian . -t checkersd_i
```

- You may want to use a smaller `Alpine` image, so create `Dockerfile-checkersd-alpine` and build the image

```bash
docker build -f prod-sim/Dockerfile-checkersd-alpine . -t checkersd_i
```

- Add Alpine section into `Makefile` for `make build-with-checksum`

```bash
make build-with-checksum
```

- Now you can run your image

```bash
docker run --rm -it checkersd_i help
```

- Now create the `Dockerfile-tmkms-debian` and `Dockerfile-tmkms-alpine` files before building the one of your choice

```bash
docker build -f prod-sim/Dockerfile-tmkms-alpine . -t tmkms_i:v0.12.2
```

- Each container needs access to its private information, such as keys, genesis, and database. To facilitate data access and separation between containers, create folders that will map as a volume to the default `/root/.checkers` or `/root/tmkms` inside containers.

```bash
mkdir -p prod-sim/kms-alice
mkdir -p prod-sim/node-carol
mkdir -p prod-sim/sentry-alice
mkdir -p prod-sim/sentry-bob
mkdir -p prod-sim/val-alice
mkdir -p prod-sim/val-bob
```

- Also add the desktop computers of Alice and Bob, so that they never have to put keys on a server that should never see them:

```bash
mkdir -p prod-sim/desk-alice
mkdir -p prod-sim/desk-bob
```

- Before you can change the configuration you need to initialize it. Do it on all nodes with this one-liner:

```bash
echo -e desk-alice'\n'desk-bob'\n'node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/prod-sim/{}:/root/.checkers \
    checkersd_i \
    init checkers
```

- Update `desk-alice/config/gensis.json` to use `upawn` denom

```bash
docker run --rm -it \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -i 's/"stake"/"upawn"/g' /root/.checkers/config/genesis.json
```

- In all seven `config/app.toml`

```bash
echo -e desk-alice'\n'desk-bob'\n'node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/prod-sim/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/([0-9]+)stake/\1upawn/g' /root/.checkers/config/app.toml
```

- Make sure that `config/client.toml` mentions `checkers-1`, the chain's name:

```bash
echo -e desk-alice'\n'desk-bob'\n'node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/prod-sim/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/^chain-id = .*$/chain-id = "checkers-1"/g' \
    /root/.checkers/config/client.toml
```

- Create on `desk-alice` the operator key for `val-alice`:

```bash
docker run --rm -it \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    add alice
```

> Use a passphrase you can remember. It does not need to be exceptionally complex as this is all a local simulation. This exercise uses `password` and stores this detail on file, which will become handy.

- Save so that we don't forget it

```bash
echo -n password > prod-sim/desk-alice/keys/passphrase.txt
```

> Because with this prod simulation you care less about safety, so much less in fact, you can even keep the mnemonic on file too.

- Do the same for val-bob

```bash
docker run --rm -it \
    -v $(pwd)/prod-sim/desk-bob:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    add bob

echo -n password > prod-sim/desk-bob/keys/passphrase.txt
```

- 

> As per the [documentation](https://github.com/iqlusioninc/tmkms/tree/v0.12.2#configuration-tmkms-init) initialize the KMS folder

```bash
docker run --rm -it \
    -v $(pwd)/prod-sim/kms-alice:/root/tmkms \
    tmkms_i:v0.12.2 \
    init /root/tmkms
```

- In the newly-created `kms-alice/tmkms.toml` file make sure that you use the right protocol version.

```bash
docker run --rm -i \
  -v $(pwd)/prod-sim/kms-alice:/root/tmkms \
  --entrypoint sed \
  tmkms_i:v0.12.2 \
  -Ei 's/^protocol_version = .*$/protocol_version = "v0.34"/g' \
  /root/tmkms/tmkms.toml
```

- Pick an expressive name for the file that will contain the `softsign` key for `val-alice`

```bash
docker run --rm -i \
  -v $(pwd)/prod-sim/kms-alice:/root/tmkms \
  --entrypoint sed \
  tmkms_i:v0.12.2 \
  -Ei 's/path = "\/root\/tmkms\/secrets\/cosmoshub-3-consensus.key"/path = "\/root\/tmkms\/secrets\/val-alice-consensus.key"/g' \
  /root/tmkms/tmkms.toml
```

- Replace cosmoshub-3 with checkers-1, the name of your blockchain, wherever the former appears

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/kms-alice:/root/tmkms \
    --entrypoint sed \
    tmkms_i:v0.12.2 \
    -Ei 's/cosmoshub-3/checkers-1/g' /root/tmkms/tmkms.toml
```

- Now you need to import `val-alice`'s consensus key in `secrets/val-alice-consensus.key`.

> The private key will no longer be needed on `val-alice`. However, during the genesis creation Alice will need access to her consensus public key. Save it in a new `pub_validator_key-val-alice.json` on Alice's desk without any new line

```bash
docker run --rm -t \
    -v $(pwd)/prod-sim/val-alice:/root/.checkers \
    checkersd_i \
    tendermint show-validator \
    | tr -d '\n' | tr -d '\r' \
    > prod-sim/desk-alice/config/pub_validator_key-val-alice.json
```

- The consensus private key should not reside on the validator. You can simulate that by moving it out

```bash
cp prod-sim/val-alice/config/priv_validator_key.json \
  prod-sim/desk-alice/config/priv_validator_key-val-alice.json
mv prod-sim/val-alice/config/priv_validator_key.json \
  prod-sim/kms-alice/secrets/priv_validator_key-val-alice.json
```

- Import it into the softsign "device" as defined in `tmkms.toml`

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/kms-alice:/root/tmkms \
    -w /root/tmkms \
    tmkms_i:v0.12.2 \
    softsign import secrets/priv_validator_key-val-alice.json \
    secrets/val-alice-consensus.key
```

- On start, `val-alice` may still recreate a missing private key file due to how defaults are handled in the code. To prevent that, you can instead copy it from `sentry-alice` where it has no value.

```bash
cp prod-sim/sentry-alice/config/priv_validator_key.json \
    prod-sim/val-alice/config/
```

- With the key created you now set up the connection from `kms-alice` to `val-alice`.

> Choose a port unused on val-alice, for instance 26659, and inform kms-alice

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/kms-alice:/root/tmkms \
    --entrypoint sed \
    tmkms_i:v0.12.2 \
    -Ei 's/^addr = "tcp:.*$/addr = "tcp:\/\/val-alice:26659"/g' /root/tmkms/tmkms.toml
```

- Do not forget, you must inform Alice's validator that it should indeed listen on port `26659` in `val-alice/config/config.toml`.

```bash
docker run --rm -i \
  -v $(pwd)/prod-sim/val-alice:/root/.checkers \
  --entrypoint sed \
  checkersd_i \
  -Ei 's/priv_validator_laddr = ""/priv_validator_laddr = "tcp:\/\/0.0.0.0:26659"/g' \
  /root/.checkers/config/config.toml
```

- Make sure it will not look for the consensus key on file

```bash
docker run --rm -i \
  -v $(pwd)/prod-sim/val-alice:/root/.checkers \
  --entrypoint sed \
  checkersd_i \
  -Ei 's/^priv_validator_key_file/# priv_validator_key_file/g' \
  /root/.checkers/config/config.toml
```

- Make sure it will not look for the consensus state file either, as this is taken care of by the KMS

```bash
docker run --rm -i \
  -v $(pwd)/prod-sim/val-alice:/root/.checkers \
  --entrypoint sed \
  checkersd_i \
  -Ei 's/^priv_validator_state_file/# priv_validator_state_file/g' \
  /root/.checkers/config/config.toml
```

- Earlier you chose checkers-1, so you adjust it here too `prod-sim/desk-alice/config/genesis.json`

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/"chain_id": "checkers"/"chain_id": "checkers-1"/g' \
    /root/.checkers/config/genesis.json
```

- In this setup, Alice starts with 1,000 PAWN and Bob with 500 PAWN, of which Alice stakes 60 and Bob 40. With these amounts, the network cannot start if either of them is offline.

> Get their respective addresses.

```bash
ALICE=$(echo password | docker run --rm -i \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    show alice --address)
```

- Have Alice add her initial balance in the genesis

```bash
docker run --rm -it \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i \
    add-genesis-account $ALICE 1000000000upawn
```

- Now move the genesis file to desk-bob. This mimics what would happen in a real-life setup.

```bash
mv prod-sim/desk-alice/config/genesis.json \
    prod-sim/desk-bob/config/
```

- Have Bob add his own initial balance

```bash
BOB=$(echo password | docker run --rm -i \
    -v $(pwd)/prod-sim/desk-bob:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    show bob --address)

docker run --rm -it \
    -v $(pwd)/prod-sim/desk-bob:/root/.checkers \
    checkersd_i \
    add-genesis-account $BOB 500000000upawn
```

- Bob is not using the Tendermint KMS but instead uses the validator key on file `priv_validator_key.json`. So, first make a copy of it on Bob's desktop.

```bash
cp prod-sim/val-bob/config/priv_validator_key.json \
    prod-sim/desk-bob/config/priv_validator_key.json
```

- Bob appears in second position in `app_state.accounts`, so his `account_number` ought to be `1`; but it is in fact written as `0`, so you use `0`:

```bash
echo password | docker run --rm -i \
    -v $(pwd)/prod-sim/desk-bob:/root/.checkers \
    checkersd_i \
    gentx bob 40000000upawn \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    --account-number 0 --sequence 0 \
    --chain-id checkers-1 \
    --gas 1000000 \
    --gas-prices 0.1upawn
```

- Again, insert Bob's chosen passphrase instead of password. Return the genesis to Alice.

```bash
mv prod-sim/desk-bob/config/genesis.json \
    prod-sim/desk-alice/config/
```

- Create Alice's genesis transaction using the specific validator public key that you saved on file, and not the key that would be taken from `priv_validator_key.json` by default (and is now missing).

```bash
echo password | docker run --rm -i \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i \
    gentx alice 60000000upawn \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    --account-number 0 --sequence 0 \
    --pubkey $(cat prod-sim/desk-alice/config/pub_validator_key-val-alice.json) \
    --chain-id checkers-1 \
    --gas 1000000 \
    --gas-prices 0.1upawn
```

- With the two initial staking transactions created, have Alice include both of them in the genesis.

```bash
cp prod-sim/desk-bob/config/gentx/gentx-* \
    prod-sim/desk-alice/config/gentx
docker run --rm -it \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i collect-gentxs
```

- As an added precaution, confirm that it is a valid genesis.

```bash
docker run --rm -it \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i \
    validate-genesis

    # File at /root/.checkers/config/genesis.json is a valid genesis file
```

- All the nodes that will run the executable need the final version of the genesis. Copy it across.

```bash
echo -e desk-bob'\n'node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    cp prod-sim/desk-alice/config/genesis.json prod-sim/{}/config
```

- What are the nodes' public keys? For instance, for `val-alice`, it is.

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/val-alice:/root/.checkers \
    checkersd_i \
    tendermint show-node-id
    # fd9dc0618f9e1d39c37dcd7a6f77d4fabc6184e7
```

- The nodes that have access to val-alice should know Alice's sentry by this identifier.

```bash
fd9dc0618f9e1d39c37dcd7a6f77d4fabc6184e7@val-alice:26656
```

> Where `val-alice` will be resolved via Docker's DNS.

- `26656` is the port as found in val-alice's configuration `prod-sim/val-alice/config/config.toml`

```bash
laddr = "tcp://0.0.0.0:26656"
```

- In the case of `val-alice`, only `sentry-alice` has access to it. Moreover, this is a persistent node. 

> `sentry-alice`'s `prod-sim/sentry-alice/config/config.toml`

```bash
persistent_peers = "fd9dc0618f9e1d39c37dcd7a6f77d4fabc6184e7@val-alice:26656"
```

- `sentry-alice` also has access to `sentry-bob` and `node-carol`, although these nodes should probably not be considered persistent. You will add them under `"seeds"`. First, collect the same information from these nodes.

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/sentry-bob:/root/.checkers \
    checkersd_i \
    tendermint show-node-id
    # 94d8f467437996e61f457efef198662a30f27d1c
docker run --rm -i \
    -v $(pwd)/prod-sim/node-carol:/root/.checkers \
    checkersd_i \
    tendermint show-node-id
    # 6501b18f1c75da22b3b657e630a117f4ce32c2f1
```

- Eventually, in `sentry-alice`'s `prod-sim/sentry-alice/config/config.toml`, you should have

```bash
seeds = "94d8f467437996e61f457efef198662a30f27d1c@sentry-bob:26656,6501b18f1c75da22b3b657e630a117f4ce32c2f1@node-carol:26656"
persistent_peers = "fd9dc0618f9e1d39c37dcd7a6f77d4fabc6184e7@val-alice:26656"
```

- Before moving on to other nodes, remember that `sentry-alice` should keep `val-alice` secret. Set

```bash
private_peer_ids = "fd9dc0618f9e1d39c37dcd7a6f77d4fabc6184e7"
```

- Repeat the procedure for the other nodes, taking into account their specific circumstances.

> Get `sentry-alice`'s node ID

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/sentry-alice:/root/.checkers \
    checkersd_i \
    tendermint show-node-id
    # 1918dd30042fcc6ab8bb5edc9f405aca32a92cd9
```

- `val-alice`'s `prod-sim/val-alice/config/config.toml`

```bash
persistent_peers = "1918dd30042fcc6ab8bb5edc9f405aca32a92cd9@sentry-alice:26656"
```

- Get `sentry-bob`'s node ID

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/sentry-bob:/root/.checkers \
    checkersd_i \
    tendermint show-node-id
    # 94d8f467437996e61f457efef198662a30f27d1c
```

- `val-bob`'s `prod-sim/val-bob/config/config.toml`

```bash
persistent_peers = "94d8f467437996e61f457efef198662a30f27d1c@sentry-bob:26656"
```

- Get `val-bob`'s node ID

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/val-bob:/root/.checkers \
    checkersd_i \
    tendermint show-node-id
    # e2b55e791a2b47bcbbcaa7cc3b7d78073c581d00
```

- Get `val-alice`'s node ID

```bash
# fd9dc0618f9e1d39c37dcd7a6f77d4fabc6184e7
```

- Get `node-carol`'s node ID

```bash
# 6501b18f1c75da22b3b657e630a117f4ce32c2f1
```

- `sentry-bob`'s `prod-sim/sentry-bob/config/config.toml`

```bash
seeds = "1918dd30042fcc6ab8bb5edc9f405aca32a92cd9@sentry-alice:26656,6501b18f1c75da22b3b657e630a117f4ce32c2f1@node-carol:26656"
persistent_peers = "e2b55e791a2b47bcbbcaa7cc3b7d78073c581d00@val-bob:26656"
private_peer_ids = "e2b55e791a2b47bcbbcaa7cc3b7d78073c581d00"
```

- `node-carols`'s `prod-sim/node-carols/config/config.toml`

```bash
seeds = "1918dd30042fcc6ab8bb5edc9f405aca32a92cd9@sentry-alice:26656,94d8f467437996e61f457efef198662a30f27d1c@sentry-bob:26656"
```

- Carol created her node to open it to the public. Make sure that her node's RPC listens on all IP addresses

```bash
docker run --rm -i \
    -v $(pwd)/prod-sim/node-carol:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei '0,/^laddr = .*$/{s/^laddr = .*$/laddr = "tcp:\/\/0.0.0.0:26657"/}' \
    /root/.checkers/config/config.toml
    # laddr = "tcp://0.0.0.0:26657"
```

- As a last step, you can disable CORS policies so that you are not surprised if you use a node from a Web browser.

```bash
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/prod-sim/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/^cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' \
    /root/.checkers/config/config.toml
    # cors_allowed_origins = ["*"]
```

- In app.toml, first location

```bash
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/prod-sim/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/^enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' \
    /root/.checkers/config/app.toml
    # enabled-unsafe-cors = true
```

- In app.toml, second location

```bash
echo -e node-carol'\n'sentry-alice'\n'sentry-bob'\n'val-alice'\n'val-bob \
    | xargs -I {} \
    docker run --rm -i \
    -v $(pwd)/prod-sim/{}:/root/.checkers \
    --entrypoint sed \
    checkersd_i \
    -Ei 's/^enable-unsafe-cors = false/enable-unsafe-cors = true/g' \
    /root/.checkers/config/app.toml
    # enable-unsafe-cors = true
```

- You are now ready to start your setup with a name other than the folder it is running in

```bash
docker compose \
    --file prod-sim/docker-compose.yml \
    --project-name checkers-prod up \
    --detach
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```

- 

```bash
```


## Notes

- Sometimes full docker use other times quick exec
- Oftentimes bob and alice need to be exported when reseting chain
- Dev container means you can do all locally
- Can't run checkersd commands with spaces in front or you get: `Error: rpc error: code = NotFound desc = rpc error: code = NotFound desc = not found: key not found`
- In going to prod I exited the dev container