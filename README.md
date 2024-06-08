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

- Now you can connect to `node-carol` to start interacting with the blockchain as you would a normal node. For instance, to ask a simple status.

```bash
docker run --rm -it \
    --network checkers-prod_net-public \
    checkersd_i status \
    --node "tcp://node-carol:26657"
```

> Whenever you submit a transaction to `node-carol`, it will be propagated to the sentries and onward to the validators.

- Get your prod setup's respective addresses for Alice and Bob

```bash
alice=$(echo password | docker run --rm -i \
    -v $(pwd)/prod-sim/desk-alice:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    show alice --address)
    # cosmos1rcqaz0qcm73lt4raet3u56avulzlv68uqxv0wy
bob=$(echo password | docker run --rm -i \
    -v $(pwd)/prod-sim/desk-bob:/root/.checkers \
    checkersd_i \
    keys \
    --keyring-backend file --keyring-dir /root/.checkers/keys \
    show bob --address)
    # cosmos1hqca6pqtqcvk7razsydx9xugjc22706r7qkyfz
```

- To stop your whole setup, run

```bash
docker compose --project-name checkers-prod down
```

> Finished up until `Self-contained checkers blockchain` in this section.

#### Tally Player Info After Production

- Add a set of stats per player

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold map playerInfo \
    wonCount:uint lostCount:uint forfeitedCount:uint \
    --module checkers --no-message
```

- You can run the tests with the verbose `-v` flag to get the log:

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go test -v github.com/alice/checkers/x/checkers/migrations/cv3/keeper
```

- After updating tests, test them

```bash
go test -v github.com/BenWolfaardt/checkers/x/checkers/migrations/cv3/keeper
go test -v github.com/BenWolfaardt/checkers/tests/integration/checkers/migrations/cv3
```

- Checkout `v1` of the chain

```bash
git checkout b0e8ddafecb2f3a040c7463f544ce5c14a9e0341
```

> It is important to note that the `/root/.checkers/*` directory that was created from our `docker-compose.yml` setup and that it needs to be deleted before we can continue with the below and start a new chain instance. Be sure to backup any necessary files or have your state in your `docker-compose.yml` configuration. 

- Build the v1 executable for your platform

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go build -o release/v1/checkersd cmd/checkersd/main.go
    # Output: release/v1/checkersd
```

> Because this is an exercise, to avoid messing with your keyring you must always specify --keyring-backend test.

- Add two players

```bash
docker network create checkers-net
docker create -it \
    -v $(pwd):/checkers -w /checkers \
    -p 26657:26657 \
    --network checkers-net \
    --name checkers \
    checkers_i
docker start checkers
docker exec -t checkers \
    ./release/v1/checkersd keys add alice --keyring-backend test

# - name: alice
#   type: local
#   address: cosmos1rdy0f88hvjqh2k7pwksut25hqxc9r58lvhtc23
#   pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AoluoXr57GKFZWd5GYTYWEtBEBJzCmcbPEqzPj23+By4"}'
#   mnemonic: "enforce crash riot produce glide genius erupt umbrella chimney that unique robot stock degree great problem normal clutch bunker wood fuel museum decade speed"

docker exec -t checkers \
    ./release/v1/checkersd keys add bob --keyring-backend test

    # - name: bob
    #   type: local
    #   address: cosmos1uftxtuzz0xs64zzddfw5dm56t3qax5h3e3vut4
    #   pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A9PWDZj0z88It4TdAPhis25I8D4KBPWJGCXMYua49E3Z"}'
    #   mnemonic: "cattle help grass surprise broom debris ozone member crumble damage hotel repeat hockey lion polar code cradle leisure dry guard asthma pigeon naive thunder"
```

- Create a new genesis

```bash
docker exec -t checkers \
    ./release/v1/checkersd init checkers-1 --chain-id checkers-1
```

- Give your players the same token amounts that were added by Ignite, as found in `config.yml`:

```bash
docker exec -t checkers \
    ./release/v1/checkersd add-genesis-account \
    alice 200000000stake,20000token --keyring-backend test
docker exec -t checkers \
    ./release/v1/checkersd add-genesis-account \
    bob 100000000stake,10000token --keyring-backend test
```

- To be able to run a quick test, you need to change the voting period of a proposal. This is found in the genesis

```bash
docker exec -it checkers \
    jq '.app_state.gov.voting_params.voting_period' /root/.checkers/config/genesis.json
    # "172800s"
```

- That is two days, which is too long to wait for CLI tests. Choose another value, perhaps 2 minutes, i.e. `"120s"`. Update it in place in the genesis

```bash
docker exec -t checkers \
    bash -c "cat <<< \$(jq '.app_state.gov.voting_params.voting_period = \"120s\"' /root/.checkers/config/genesis.json) \
    > /root/.checkers/config/genesis.json"
```

> You can confirm that the value is in using the earlier command.

- Make Alice the chain's validator too by creating a genesis transaction modeled on that done by Ignite, as found in `config.yml`

```bash
docker exec -t checkers \
    ./release/v1/checkersd gentx alice 100000000stake \
    --keyring-backend test --chain-id checkers-1
docker exec -t checkers \
    ./release/v1/checkersd collect-gentxs
```

- Now you can start the chain proper

```bash
docker exec -it checkers \
    ./release/v1/checkersd start \
    --rpc.laddr "tcp://0.0.0.0:26657"
```

- Add games

```bash
export alice=$(docker exec checkers ./release/v1/checkersd keys show alice -a --keyring-backend test)
export bob=$(docker exec checkers ./release/v1/checkersd keys show bob -a --keyring-backend test)
docker exec -t checkers \
    ./release/v1/checkersd tx checkers create-game \
    $alice $bob 10 stake \
    --from $alice --keyring-backend test --yes \
    --chain-id checkers-1 \
    --broadcast-mode block
```

<!-- TODO: Requires the integration tests from the CosmJS for Your Chain section -->

- For the software upgrade governance proposal, you want to make sure that it stops the chain not too far in the future but still after the voting period. With a voting period of 10 minutes, take 15 minutes. How many seconds does a block take?

```bash
docker exec -t checkers \
    bash -c 'jq -r ".app_state.mint.params.blocks_per_year" /root/.checkers/config/genesis.json'
    # 6311520  <-- 5 seconds per block. At this rate, 2 minutes mean 24 blocks.
```

- What is the current block height?

```bash
docker exec -t checkers \
    bash -c './release/v1/checkersd status \
        | jq -r ".SyncInfo.latest_block_height"'
    # 30
```

> Take this block heigh and add 180 therefore, `--upgrade-height <30 + 24 >`

- What is the minimum deposit for a proposal?

```bash
docker exec -t checkers \
    bash -c 'jq ".app_state.gov.deposit_params.min_deposit" \
        /root/.checkers/config/genesis.json'

    # [
    #     {
    #         "denom": "stake",
    #         "amount": "10000000"
    #     }
    # ]
```

> This is the minimum amount that Alice has to deposit when submitting the proposal. This will do: `--deposit 10000000stake`

- Submit your governance proposal upgrade

> Update to your desired `--upgrade-height <height>`

```bash
docker exec -t checkers \
    ./release/v1/checkersd tx gov submit-proposal software-upgrade v1tov1_1 \
    --title "v1tov1_1" \
    --description "Introduce a tally of games per player" \
    --from $alice --keyring-backend test --yes \
    --chain-id checkers-1 \
    --broadcast-mode block \
    --upgrade-height 54 \
    --deposit 10000000stake

    # - attributes:
    #   - key: proposal_id
    #     value: "1"
    #   - key: proposal_type
    #     value: SoftwareUpgrade
    #   - key: voting_period_start
    #     value: "1"
    #   type: submit_proposal
```

- Where `1` is the proposal ID you reuse. Have Alice and Bob vote yes on it

```bash
docker exec -t checkers \
    ./release/v1/checkersd tx gov vote 1 yes \
    --from $alice --keyring-backend test --yes \
    --chain-id checkers-1
docker exec -t checkers \
    ./release/v1/checkersd tx gov vote 1 yes \
    --from $bob --keyring-backend test --yes \
    --chain-id checkers-1
```

- Confirm that it has collected the votes

```bash
docker exec -t checkers \
    ./release/v1/checkersd query gov votes 1

    # votes:
    # - option: VOTE_OPTION_YES
    #   options:
    #   - option: VOTE_OPTION_YES
    #     weight: "1.000000000000000000"
    #   proposal_id: "1"
    #   voter: cosmos1zkmc72ac9zjy2gr5dfdmnkj6wqd072ef8ntses
    # - option: VOTE_OPTION_YES
    #   options:
    #   - option: VOTE_OPTION_YES
    #     weight: "1.000000000000000000"
    #   proposal_id: "1"
    #   voter: cosmos1394z74aprj2x8xwxlf7jyj3zp3ns79u6k9heal
```

- See how long you have to wait for the chain to reach the end of the voting period

```bash
docker exec -t checkers \
    ./release/v1/checkersd query gov proposal 1

    # ...
    # status: PROPOSAL_STATUS_VOTING_PERIOD
    # ...
    # voting_end_time: "2024-06-08T17:30:09.881755722Z"
    # ...
```

- Wait for this period. Afterward, with the same command you should see

```bash
...
status: PROPOSAL_STATUS_PASSED
...
```

- Now, wait for the chain to reach the desired block height, which should take five more minutes, as per your parameters. When it has reached that height, the shell with the running `checkersd` should show something like

```bash
...
5:55PM INF finalizing commit of block hash=E6CBD7E5A7FA47ECB3259DDB73DFA69D815EA99101E2E1A1EF5B925FC0380ED5 height=75 module=consensus num_txs=0 root=2CE9C9B08BAE3D66710E3594C23645AA45F820CD2640C7D403D13F6F4DC90DEF
5:55PM ERR UPGRADE "v1tov1_1" NEEDED at height: 75: 
5:55PM ERR CONSENSUS FAILURE!!! err="UPGRADE \"v1tov1_1\" NEEDED at height: 75: " module=consensus stack="goroutine 23 [running]:\nruntime/debug.Stack
...
5:55PM INF Stopping baseWAL service impl={"Logger":{}} module=consensus wal=/root/.checkers/data/cs.wal/wal
5:55PM INF Stopping Group service impl={"Dir":"/root/.checkers/data/cs.wal","Head":{"ID":"pMLx3SbYdq5D:/root/.checkers/data/cs.wal/wal","Path":"/root/.checkers/data/cs.wal/wal"},"ID":"group:pMLx3SbYdq5D:/root/.checkers/data/cs.wal/wal","Logger":{}} module=consensus wal=/root/.checkers/data/cs.wal/wal
5:55PM INF Timed out dur=3000 height=75 module=consensus round=0 step=3
...
```

- At this point, run in another shell

```bash
docker exec -it checkers \
    bash -c './release/v1/checkersd status \
        | jq -r ".SyncInfo.latest_block_height"'
```

- You should always get the same value, no matter how many times you try. That is because the chain has stopped. For instance:

```bash
75
```

- Stop `checkersd` with CTRL-C. It has saved a new file:

```bash
docker exec -it checkers \
    cat /root/.checkers/data/upgrade-info.json
```

- This prints

```bash
{"name":"v1tov1_1","height":75}
```

- With your node (and therefore your whole blockchain) down, you are ready to move to v1.1

```bash
git checkout master
```

- Back in the first shell, build the v1.1 executable

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go build -o ./release/v1_1/checkersd ./cmd/checkersd/main.go
```

- Launch it

```bash
docker exec -it checkers \
    ./release/v1_1/checkersd start \
    --rpc.laddr "tcp://0.0.0.0:26657"
```

- It should start and display something like

```bash
...
6:02PM INF applying upgrade "v1tov1_1" at height: 75
6:02PM INF migrating module checkers from version 2 to version 3
6:02PM INF Start to compute checkers games to player info calculation...
6:02PM INF Checkers games to player info computation done
...
```

- After it has started, you can confirm in another shell that you have the expected player info with

```bash
docker exec -t checkers \
    ./release/v1_1/checkersd query checkers list-player-info
```

- This should print something like

<!-- TODO: fill in data -->

```bash
playerInfo:
- forfeitedCount: "0"
  index: cosmos1fx6qlxwteeqxgxwsw83wkf4s9fcnnwk8z86sql
  lostCount: "0"
  wonCount: "3"
- forfeitedCount: "0"
  index: cosmos1mql9aaux3453tdghk6rzkmk43stxvnvha4nv22
  lostCount: "3"
  wonCount: "0"
```

> Congratulations, you have upgraded your blockchain almost as if in production!

- You can stop Ignite CLI. If you used Docker that would be

```bash
$ docker stop checkers
$ docker rm checkers
$ docker network rm checkers-net
```

#### Add a Leaderboard Module

- Do so with ignite

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold module leaderboard \
    --params length:uint
```

- Add a structure for the leaderboard: you want a single stored leaderboard for the whole module

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold single leaderboard winners \
    --module leaderboard --no-message
```

- Update the various `.proto` files and then run

```bash
docker exec -it checkers \
    ignite generate proto-go
```

- Confirm all of your tests are passing

```bash
docker exec -it checkers \
    go test github.com/BenWolfaardt/checkers/x/leaderboard/keeper
docker exec -it checkers \
    go test github.com/BenWolfaardt/checkers/x/leaderboard/types
```

- After adding `Candidate` to `leaderboard.proto` run

```bash
docker exec -it checkers \
    ignite generate proto-go
```

- Create new mocks to test the `MultiHook` type we added

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    make mock-expected-keepers
```

> Re run your tests

- Your v2 blockchain is fully functioning. It will work as long as you start it from scratch

```bash
docker network create checkers-net
docker run --rm -it \
    -v $(pwd):/checkers -w /checkers \
    -p 4500:4500 -p 26657:26657 \
    --network checkers-net \
    --name checkers \
    checkers_i \
    ignite chain serve --reset-once
```

<!-- TODO: require the CosmJS integrations tests here which section was skipped -->

- After that, you can query your leaderboard

```bash
docker exec -it \
    checkers \
    checkersd query leaderboard show-leaderboard

# Leaderboard:
#   winners:
#   - addedAt: "1682373982"
#     address: cosmos1fx6qlxwteeqxgxwsw83wkf4s9fcnnwk8z86sql
#     wonCount: "1"
```

#### Migrate the Leaderboard Module After Production

- You introduced a new expected keeper. If you want to unit test your migration helpers properly, you have to mock this new expected interface. Update the `Makefile` and then

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    make mock-expected-keepers
```

- Write your tests and then test them

```bash
docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    go test github.com/BenWolfaardt/checkers/x/leaderboard/migrations/cv2/keeper
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
- I deleted all .checkers in root before the Tally player info 
