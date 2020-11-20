# tdex-feeder

Feeder allows to connect an external price feed to the TDex Daemon to determine the current market price

## Overview

tdex-feeder connects to exchanges and retrieves market prices in order to consume the gRPC 
interface exposed from tdex-deamon `UpdateMarketPrice`.

## ⬇️ Run  Standalone

### Install

1. [Download the latest release for MacOS or Linux](https://github.com/tdex-network/tdex-feeder/releases)

2. Move the feeder into a folder in your PATH (eg. `/usr/local/bin`) and rename the feeder as `feederd`

3. Give executable permissions. (eg. `chmod a+x /usr/local/bin/feederd`)

### Run
```sh
# Run with default config and default flags.
$ feederd

# Run with debug mode on.
$ feederd -debug true

# Run with debug mode and different config path.
$ feederd -debug true -conf ./config.json
```

## 📄 Usage

In-depth documentation for using the tdex-feeder is available at [docs.tdex.network](https://docs.tdex.network/tdex-feeder.html)

## 🖥 Local Development

Below is a list of commands you will probably find useful.

### Linux

`make build-linux`

### Mac

`make build-mac`

### Run Linux

`make run-linux`

#### Flags

```
-conf: Configuration File Path. Default: "./config/config.json"
-debug: Log Debug Informations Default: "false"
```

#### Config file

You can find a config file with an working example in `./config.example.json`
that will be loaded except if there are instructions otherwise.

```
daemon_endpoint: String with the address and port of gRPC host. Required.
daemon_macaroon: String with the daemon_macaroon necessary for authentication.
kraken_ws_endpoint: String with the address and port of kraken socket. Required.
markets: Json List with necessary markets informations. Required.
base_asset: String of the Hash of the base asset for gRPC request. Required.
quote_asset: String of the Hash of the quote asset for gRPC request. Required.
kraken_ticker: String with the ticker we want kraken to provide informations on. Required.
interval: Int with the time in secods between gRPC requests. Required.
```