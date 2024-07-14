# Consumer Testing Integration Guide

## Overview
This document covers the setup and integration of Ethos' tests with a given consumer chain to validate correct integration of `x/ccv/consumer` as well as `x/ccv/democracy/staking`, `x/ccv/democracy/governance`, and `x/ccv/democracy/distribution`.

## Tests
>Note: At the time of writing, Ethos E2E tests must be run individually via `go run`. Make targets for groups of like-tests and important milestones like `happy-path` are WIP.

### E2E Test Suite

#### Architecture
Currently Ethos' E2E configuration separates everything for `ethos-avs` and `ethos-chain` and builds two images, `cosmos-ics:local` for the shared security network, and `anvil` for the mock Ethereum network.

When setting up tests, two containers are run for each image, `interchain-security-instance` and `anvil`. Each test specifies a starting configuration, and predefined steps. Test configurations are initialized before any steps can run. Each step has an action that is executed, and a final expected chain state for both provider and consumer after each action.

<br>

#### Containers
**`anvil`** contains a mock Ethereum network with the Ethos AVS deployed, using Anvil. This container copies over the contracts and scripts from `ethos-avs`, and then executes the `start-anvil-chain-with-el-and-avs-deployed.sh` script to start anvil with the AVS deployed.

There are three actions in the E2E tests that interact with the Ethos AVS:
- `SetupAvsOperator`: Registers an operator with EigenLayer, deposits tokens into an EigenLayer strategy, whitelists the operator on the Ethos AVS, and registers the operator with the Ethos AVS
- `RegisterConsumerChain`: Registers a consumer chain on the Ethos AVS
- `BondAction`: Opts an operator into a consumer chain on the Ethos AVS, or can also be used to change their stake percentage on the consumer chain

**`interchain-security-instance`** contains the entire provider network (3 validator nodes, 1 query node), the (hermes) relayer, and consumer chain network (3 validator nodes), all running under separate network namespaces. This is set to change in the future to accomodate easier integration for consumer chain partners. Integrators can expect a separate container for the consumer chain network.

<br>

#### Test Cases
_**Happy Path**_

`happy-path` tests the basic functionality of bonding operators to consumer chain validators, as well as the avs-oracle through unbonding actions and changing an operators stake to its delegated consumer validator. Aftwards the test submits a governance proposal to remove the a consumer chain. The proposal is meant to fail and the consumer chain should not be removed. Further, this test ensures that relayers are working as expected, relaying packets between the ethos provider and consumer chains.


_**Democracy**_

Democracy tests cover behavior with respect to electing representatives (governators) and subsequent reward distribution for governators and their delegators given happy path conditions and misbehavior edge cases.

The only difference between `democracy` and `democracy-reward` is the expectation that consumer rewards are shared on the provider chain in the native token. This is enabled in the test configuration during setup.

First governators are registered on the consumer chain, and then delegated to. A change in only representative power and not validator powers is validated, as well as expected reward token distribution to representatives and delegators.

The consumer chain submits a governance proposal to update transfer params to set send and receive enabled to true, which is expected to pass. Before the consumer denom is registered on the provider chain, the test verifies that no rewards have been distributed on the provider. Subsequently a proposal on the provider chain to register the consumer denom.

Finally slashing is tested by triggering downtime for one of the validators on the consumer chain. Once downtime is registered on both the consumer and provider, the test validates that governators and their delegators were not slashed.


_**Misbehavior**_

Misbehavior tests cover expected behavior and relayer detection and propagation of evidence for both double signing and downtime infractions executed on the consumer chain. 

For both **`consumer-double-sign`** and **`consumer-downtime`**, misbehavior is triggered on the consumer chain via scripts found in [testnet-scripts](https://github.com/Ethos-Works/ethos/tree/main/ethos-chain/tests/e2e/testnet-scripts). 

Evidence detection is started on the relayer, which will poll for evidence on the consumer chain and submit to the provider. When registered on the provider, it is expected that the correpsonding operator is frozen on all consumer chains it is operating on.

Finally once all packets are relayed to the consumer chain, the tests validate that offending validators are slashed.

<br>

#### Debugging
Debugging tests can be challenging with the current setup. To pause at any step and inspect logs in the `ethos-ics` container, add the following line before any given action:

```go
func (tr *TestConfig) WaitTime(duration time.Duration) {
   tr.waitBlocks(chain, 500, 360*time.Second)
   ...
}
```
All chain binaries can be found in `/usr/local/bin`.

Navigate to the following for provider and consumer node home and logs

Provider
```
/provi
 ├── /validatoralice
 ├── /validatorbob
 └── /validatorcarol
```

Consumer
```
/consu
 ├── /validatoralice
 ├── /validatorbob
 └── /validatorcarol
```

All test scripts can be found in `/testnet-scripts`. Certain scripts for doublesigning and misbehavior will create validator directories in `/consu`.

<br>

#### Integration

At present, all E2E test must be run independently until make targets are working.

To run Ethos E2E from the [Ethos chain repo](https://github.com/Ethos-Works/ethos/tree/main/ethos-chain) with the example consumer and democracy apps, no additional setup is required. The Dockerfile will copy over the ethos directory and build each binary for ethos, consumer and consumer-democracy.

To run Ethos E2E tests with Ethos as the provider binary and a custom consumer chain binary, run the following in the main directory of the [`ethos-chain` repository](https://github.com/Ethos-Works/ethos/tree/main/ethos-chain).
```
go run ./tests/e2e/... -tc <test-case> -cv-git <consumer git url> -cv <consumer git commit>  -pv <ethos provider commit>   --verbose
```

#### Custom Cases

To write additional test cases, see the guide to [defining a new test case
](https://github.com/Ethos-Works/ethos/tree/main/ethos-chain/tests/e2e#defining-a-new-test-case) in the `ethos-chain` E2E directory.

<hr>

### Integration Test Suite

Ethos also provides an integration test suite that partner chains can integrate into their own repositories by wiring their application. 

See both the [example configuration in Ethos](https://github.com/Ethos-Works/ethos/blob/main/ethos-chain/tests/integration/instance_test.go) and the [consumer template](./integration/ccv_test.go).