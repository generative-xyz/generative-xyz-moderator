# rendering-hub

## How to run
```
make run-server
```

```
 abigen --pkg generative_nft_contract --abi  ./artifacts/contracts/nfts/GenerativeNFT.sol/GenerativeNFT.json --out /Users/autonomous/go/src/rendering-hub/utils/contracts/generative_nft_contract/GenerativeNFT.go
```

```
 abigen --pkg generative_project_contract --abi  ./artifacts/contracts/nfts/GenerativeProject.sol/GenerativeProject.json --out /Users/autonomous/go/src/rendering-hub/utils/contracts/generative_project_contract/GenerativeProject.go
```

```
abigen --pkg generative_project_data --abi  ./artifacts/contracts/data/GenerativeProjectData.sol/GenerativeProjectData.json --out /Users/autonomous/go/src/rendering-hub/utils/contracts/generative_project_data/GenerativeProjectData.go
```

```
abigen --pkg generative_marketplace --abi  ./artifacts/contracts/services/AdvanceMarketplaceService.sol/AdvanceMarketplaceService.json --out /Users/autonomous/go/src/rendering-hub/utils/contracts/generative_marketplace/GenerativeMarketplace.go
```

```
abigen --pkg generative_marketplace_lib --abi  ./artifacts/contracts/libs/structs/Marketplace.sol/Marketplace.json --out /Users/autonomous/go/src/rendering-hub/utils/contracts/generative_marketplace_lib/GenerativeMarketplaceLib.go
```
 