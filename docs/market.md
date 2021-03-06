# Market Module

## Overview

The market module enable accounts to sell an amount of *source* tokens in exchange for an fixed amount of *destination* tokens, where the price is derived from the amount of *source* and *destination* tokens.

This is a generalisation of the classic limit order for two-sided markets:

| Side | Base | Term | Price |
|------|------|------|-------|
| Buy  | Destination Denom | Source Denom | Source Amount / Destination Amount |
| Sell | Source Denom | Destination Denom | Destination Amount / Source Amount |

As the *destination* amount is fixed, less than *source* amount of tokens will be paid if a better price exist in the market. Having the *destination* amount fixed is useful for payments where a fixed amount of foreign currency needs to be delivered.

## Order Data

An order consists of the following data:

* Owner: a `AccAddress` which will be used for settlement and order modifications.
* OrderId: a `uint64` assigned by the market module, monotonically increasing.
* ClientOrderId: a `string` assigned by owner, which must not be a duplicate of an existing order.
* Source: a `Coin` representing the desired amount of tokens to sell.
* SourceFilled: `Int` that tracks the sold amount so far.
* SourceRemaining: a `Int` that is adjusted with *SourceFilled* and if the owner account balance change.
* Destination: a `Coin` representing the minimum amount of tokens to buy.
* DestinationFilled: `Int` that tracks the bought amount so far.
* Price: a `Dec` calculated as *Destination* / *Source*.

## Features

*No instrument listing required*. Any token is immediately tradeable against other tokens.

*No execution fees*. This applies for both makers and takers, which only need to pay the standard transaction costs.

*Optimized for liquidity*. Orders do not touch the account balance until they are matched, so that makers can place multiple orders based on the same *Source*.
When the balance of the owner account changes, SourceRemaining is adjusted accordingly and any untradable orders are canceled. 

*Takers always trade at the best price*. In case there is a better price in the market, price improvement is passed to the taker who pays less than the specified amount of *Source* tokens.

*Arbitrage-free*. Sophisticated order matching ensures that no arbitrage opportunities exist in the market. Orders always trade at the best price by considering synthetic instruments, e.g. a single USD->EUR order matched against EUR->GBP and GBP->USD simultaneously.

*Price/time priority matching*. Orders at the same price will be ordered by OrderId, with the lowest matched first.  

*Immediate settlement*. Matched orders are settled immediately with finality.

## Transaction Types

The transaction types mirror those of the [FIX trading specification](https://www.fixtrading.org/online-specification/business-area-trade/) for single order handling.

### NewOrderSingle

Adds a new order to the order book. The ClientOrderId must be unique among existing orders for the same owner account.

### CancelOrder

Cancels the remaining part of an existing order, referenced by it's ClientOrderId. In case the order has already been fully filled, an error will be returned. 

### CancelReplaceOrder

Cancels the remaining part of an existing order, referenced by it's ClientOrderId. The filled part of the cancelled order is carried over into the new order.
