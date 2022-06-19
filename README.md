# go-dydx

## Documentation

See the documentation on the package page:

[![Go Reference](https://pkg.go.dev/badge/github.com/fardream/go-dydx.svg)](https://pkg.go.dev/github.com/fardream/go-dydx)

golang client for [dydx.exchange](https://dydx.exchange), supports:

- onboarding

  - create user.
  - deterministic recover api key and stark key.

- private api

  - get user, accounts, positions, orders, withdrawals, fills, funding, and pnl.
  - create, cancel orders and active orders.
  - subscription to account updates.

- public api

  - get markets, orderbooks, trades, candles, historical fundings.
  - subscription to markets, orderbooks, trades.

## Prior Art

This is based on the work from [go-numb](https://github.com/go-numb) at [here](https://github.com/go-numb/go-dydx) with some go idiomatic modifications.

There is also another version from [verichenn](https://github.com/verichenn) [here](https://github.com/verichenn/dydx-v3-go).

## Command Line Interface

A command line interface is provided in [`dydx-cli`](dydx-cli/). To install:

```shell
go install github.com/fardream/go-dydx/dyx-cli@latest
```

A command line interface is provided to replay the orderbook updates [`dydx-replay-orderbook`](dydx-replay-orderbook/). To install:

```shell
go install github.com/fardream/go-dydx/dyx-replay-orderbook@latest
```
