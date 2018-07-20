# ASSET

An application which can create, manage and destroy asset's on bloackchain. It also provide a smart way to buy asset by creating smart contract at the time of asset creation.

## Specification

HyperLedger fabric 1.2.0 is used for creating the wallet. Chaincode is written in golang. All wallet related information will be stored on hyperledger blockchain and hyperLedger fabric will be used to interact  application.

## Getting Started

### Prerequisites

Make sure docker and docker-compose is installed on the system. You can check if they are installed using following commands

```bash
docker -v
docker-compose -v
```

### Terminal 1 - Build network & Copy project files
```bash
curl -sSL http://bit.ly/2ysbOFE | bash -s 1.2.0
# copy asset folder inside fabric-samples/chaincode directory
cp -rf asset  fabric-samples/chaincode/
cd fabric-samples/chaincode-docker-devmode/  
docker-compose -f docker-compose-simple.yaml up
```

[above commands will start the fabric network bare minimum requirements you can refer
[chaincode developer tutorial](https://hyperledger-fabric.readthedocs.io/en/release-1.2/chaincode4ade.html)]

### Terminal 2 - Build & start the chaincode
```bash
docker exec -it chaincode bash
cd asset
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=asset:0 ./asset
```

### Terminal 3 - Use the chaincode
```bash
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/asset -n asset -v 0

# I don't understand significance of -C myc if anyone do let me know
peer chaincode instantiate -n asset -v 0 -C myc
peer chaincode invoke -n asset  -c '{"Args":["createAsset","asset1","asset1pass","assetowner","assetowner Is awesome","[{\"fraction\":0.5,\"walletAddress\":\"wallet1\"},{\"fraction\":0.5,\"walletAddress\":\"wallet2\"}]"]}' -C myc
peer chaincode invoke -n asset  -c '{"Args":["getAsset","asset1","asset1pass"]}' -C myc
```

## Usage

### 1. Asset Structures

```golang
type stakeholders struct {
        Fraction float64 `json:"fraction"`
        WalletAdress string `json:"walletAddress"`
}

type asset struct {
        Name string `json:"name"`
        Password string `json:"password"`
        Owner string `json:"Owner"`
        OwnerInfo string `json:"OwnerInfo"`
        Stakeholders []stakeholders `json:"stakeholder"`
}
```

### 2. Create Asset

Create wallet is used to initialise a new wallet.

args:- 1."name of asset" 2."password for asset" 3."name of owner" 4."owner information" 5."stakeholders"

example:-
```bash
peer chaincode invoke -n asset  -c '{"Args":["createAsset","asset1","asset1pass","assetowner","assetowner Is awesome","[{\"fraction\":0.5,\"walletAddress\":\"wallet1\"},{\"fraction\":0.5,\"walletAddress\":\"wallet2\"}]"]}' -C myc
```

respose:- 
Chaincode invoke successful. result: status:200

error response:-
chaincode result: status:500 message:"asset already exists"

### 3. Get Asset

Get asset information

args:- 1. fromWallet 2. towallet 3. amount  4. password

example :-
```bash
peer chaincode invoke -n asset  -c '{"Args":["getAsset","asset1","asset1pass"]}' -C myc
```

response:-
Chaincode invoke successful. result: status:200 payload:

### 4. Buy Asset

whenever a person wants to buy a copy of asset this function is envoked

args:- 1. assetname 2. price 3. customerWallet 4. password

example:-
```bash
peer chaincode invoke -n asset  -c '{"Args":["buyAsset","asset1","20","masterWallet","masterpass"]}' -C myc
```

response: -
Chaincode invoke successful. result: status: 200
## Build With

[hyperledger-fabric](https://www.hyperledger.org/projects/fabric) - a blockchain framework implementation and one of the Hyperledger projects hosted by The Linux Foundation.
[golang 1.10](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software

## Contributing

Coming soon

## Authors

- Parag Rahangdale

## licence

This project is licensed under the Apache-2.0 Licence

## Future Works

- [ ] find a way to properly deploy [not in dev mode] the chaincode on hyperledger Fabric
- [ ] Create a Nodejs client to interact with chaincode
- [ ] Check for wallet existence while creating an asset
- [ ] function for deleting the assets