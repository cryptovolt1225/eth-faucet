# eth-faucet
The faucet is a web application with the goal of distributing small amounts of Ether in private and test networks.

## Get started
### Prerequisites

* Go (1.16 or later)
* Node.js

### Build

1. Clone the repository
```bash
git clone https://github.com/chainflag/eth-faucet.git
cd eth-faucet
```

2. Bundle web using rollup
```bash
cd web && npm install
npm run build
```
_For more details, please refer to the [web readme](https://github.com/chainflag/eth-faucet/blob/main/web/README.md)_  

3. Build binary application
```bash
cd ..
go build -o eth-faucet main.go
```

## Usage
First create config file `config.yml` based on the [example](https://github.com/chainflag/eth-faucet/blob/main/config.sample.yml)
* `provider` Ethereum json rpc endpoint
* `wallet` Funder account specified via hex private key or keystore. Faucet will use private key first if it is not empty.

then start the faucet server  
```
./eth-faucet -config config.yml
```
or run faucet using docker without the compiled binary  
```bash
docker run -d -p 8080:8080 -v `pwd`/config.yml:/app/config.yml -v `pwd`/keystore:/app/keystore chainflag/eth-faucet
```

**Optional Flags**

| flag      | Description                                      | Default Value
| --------- | ------------------------------------------------ | -------------
| -apiport  | Listener port to serve HTTP connection           | 8080
| -config   | Path of wallet config yaml file                  | config.yml
| -interval | Number of minutes to wait between funding rounds | 1440
| -payout   | Number of Ethers to transfer per user request    | 1
| -queuecap | Maximum transactions waiting to be sent          | 100

## License
This project is licensed under the MIT License
