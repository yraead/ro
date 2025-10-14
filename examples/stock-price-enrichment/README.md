
# Stock price enrichment

Connect to a websocket and aggregate 1min trade data to expose low, high, last price and volume.

## Run

Create an API key here: https://finnhub.io

```bash
export FINNHUB_API_KEY=xxxx
go run main.go
```

You will be asked to paste a stock symbol. 👇

```bash
Subscribe to a stock symbol: BINANCE:BTCUSDT

[BINANCE:BTCUSDT] -> Last:123070 Low:123055 High:123081 Volume:10.715
[BINANCE:BTCUSDT] -> Last:123089 Low:123070 High:123089 Volume:3.4653
[BINANCE:BTCUSDT] -> Last:123064 Low:123061 High:123105 Volume:5.218
[BINANCE:BTCUSDT] -> Last:123148 Low:123064 High:123148 Volume:8.5387
[BINANCE:BTCUSDT] -> Last:123160 Low:123148 High:123160 Volume:1.6681
[BINANCE:BTCUSDT] -> Last:123234 Low:123160 High:123234 Volume:7.9309
[BINANCE:BTCUSDT] -> Last:123219 Low:123217 High:123234 Volume:3.0173
[BINANCE:BTCUSDT] -> Last:123219 Low:123219 High:123219 Volume:1.1453
...
```

Examples of tickers you can track:
- BINANCE:BTCUSDT (24/7)
- AAPL (during trading hours)
- TSLA (during trading hours)
- NVDA (during trading hours)

For more informations, read the [Finnhub doc](https://finnhub.io/docs/).
