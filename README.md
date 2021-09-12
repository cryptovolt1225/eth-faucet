# eth-faucet
The faucet is a web application with the goal of distributing small amounts of Ether in private and test networks.

## Get started

### Run with Docker

* Use private key as funder
```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER="rpc endpoint" -e PRIVATE_KEY="hex private key" chainflag/eth-faucet -apiport 8080
```

* Use keystore file as funder
```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER="rpc endpoint" -e KEYSTORE="keystore path" -e PASSWORD="keystore pass" -v `pwd`/keystore:/app/keystore chainflag/eth-faucet -apiport 8080
```

**Optional Flags**

| flag        | Description                                      | Default Value
| ----------- | ------------------------------------------------ | -------------
| -apiport    | Listener port to serve HTTP connection           | 8080
| -interval   | Number of minutes to wait between funding rounds | 1440
| -payout     | Number of Ethers to transfer per user request    | 1
| -proxycount | Count of reverse proxies in front of the server  | 0
| -queuecap   | Maximum transactions waiting to be sent          | 100

### Deploy to Heroku

```bash
heroku create
heroku buildpacks:add heroku/nodejs
heroku buildpacks:add heroku/go
heroku config:set WEB3_PROVIDER=rpc endpoint
heroku config:set PRIVATE_KEY=hex private key
git push heroku main
heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## Development

### Prerequisites

* Go (1.16 or later)
* Node.js

### Build

1. Clone the repository and navigate to the app’s directory
```bash
git clone https://github.com/chainflag/eth-faucet.git
cd eth-faucet
```

2. Bundle Frontend web with Rollup
```bash
npm run build
```
_For more details, please refer to the [web readme](https://github.com/chainflag/eth-faucet/blob/main/web/README.md)_  

3. Build binary application to run faucet
```bash
go build -o eth-faucet main.go
export WEB3_PROVIDER=https://ropsten.infura.io
export PRIVATE_KEY=secret
./eth-faucet
```

## License

This project is licensed under the MIT License
