# hypermint

[![CircleCI](https://circleci.com/gh/bluele/hypermint.svg?style=svg)](https://circleci.com/gh/bluele/hypermint)

hypermint = wasm + tendermint

This is a very experimental project.

## Build

```
$ dep ensure
$ make build
```

## Getting started

### Run a validator node

First, you need to initialize genesis state.

```
# these outputs will be different per execution
$ make init
{
  "chain_id": "test-chain-6AHEow",
  "node_id": "67b4f60a2b371a908848af2d35e7816b55610115",
  "app_message": "success"
}
export ADDR1=0x1221a0726d56aEdeA9dBe2522DdAE3Dd8ED0f36c
export ADDR2=0xD8eba1f372b9e0D378259F150d52C2e6C2e4109a
```

Next, run a blockchain node:

```
$ make start
```

### Smart contract

hypermint supports wasm based smart contract.

Contract example project is [here](https://github.com/bluele/hypermint/tree/develop/example).

*If you don't have cargo and wasm-gc, you should install these.*

- https://doc.rust-lang.org/cargo/getting-started/installation.html
- https://github.com/alexcrichton/wasm-gc

To deploy [simple token project](https://github.com/bluele/hypermint/tree/develop/example/token), exec below commands:

```
# '0x1221a0726d56aEdeA9dBe2522DdAE3Dd8ED0f36c' should be replace with the value which was got by `make init`
$ export ADDR1=0x1221a0726d56aEdeA9dBe2522DdAE3Dd8ED0f36c

# To exec deploy cmd, cargo with wasm32 and wasm-gc
$ make -C ./example/token deploy
cargo build --target=wasm32-unknown-unknown
   Compiling hmc v0.1.0 (/Users/jun/go/src/github.com/bluele/hypermint/hmc)
   Compiling token v0.1.0 (/Users/jun/go/src/github.com/bluele/hypermint/example/token)
    Finished dev [unoptimized + debuginfo] target(s) in 2.08s
wasm-gc ./target/wasm32-unknown-unknown/debug/token.wasm -o ./token.min.wasm
contract address is 0xceD4629963CCc0549094e962a01f454EBFD80Cbd
```

Now you got the first contract address!
Next, try to check your balance.

```
$ ./build/hmcli contract call --address=$ADDR1 --contract=0xceD4629963CCc0549094e962a01f454EBFD80Cbd --func="get_balance" --password=password --simulate --gas=1
100
```

## Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
