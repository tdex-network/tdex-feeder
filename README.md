# tdex-feeder

Feeder allows to connect an external price feed to the TDex Daemon to determine the current market price

## ⬇️ Install

### Build 

#### Linux

`make build-linux`

### Run

`./build/feederd-linux-amd64`

## 📄 Usage

In-depth documentation for using the tdex-feeder is available at [docs.tdex.network](https://docs.tdex.network/tdex-feeder.html)

## 🖥 Local Development

Below is a list of commands you will probably find useful.

### Run 

`go run ./cmd/feederd/main.go`

#### Flags

`-conf: Configuration File Path. Default: "./config/config.json"`
`-debug: Log Debug Informations Default: "false"`