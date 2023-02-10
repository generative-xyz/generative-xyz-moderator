package eth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/inc-backend/crypto-libs/eth/bridge"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

type Client struct {
	eth *ethclient.Client
}

func NewClient(eth *ethclient.Client) *Client {
	return &Client{eth}
}

func (c *Client) GetClient() *ethclient.Client {
	return c.eth
}

func (c *Client) GenerateAddress() (privKey, pubKey, address string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		err = errors.Wrap(err, "crypto.GenerateKey")
		return
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privKey = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("failed to cast public key to ECDSA")
		return
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	pubKey = hexutil.Encode(publicKeyBytes)[4:]

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return
}

func (c *Client) GenPubPriKeyFromIncPriKey(incPrivateKey []byte) (ecdsa.PrivateKey, ecdsa.PublicKey) {
	priKey := new(ecdsa.PrivateKey)
	priKey.Curve = crypto.S256()
	priKey.D = c.b2ImN(incPrivateKey)
	priKey.PublicKey.X, priKey.PublicKey.Y = priKey.Curve.ScalarBaseMult(priKey.D.Bytes())
	return *priKey, priKey.PublicKey
}

func (c *Client) b2ImN(bytes []byte) *big.Int {
	x := big.NewInt(0)
	x.SetBytes(crypto.Keccak256Hash(bytes).Bytes())
	for x.Cmp(crypto.S256().Params().N) != -1 {
		x.SetBytes(crypto.Keccak256Hash(x.Bytes()).Bytes())
	}
	return x
}

func (c *Client) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := c.eth.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, errors.Wrap(err, "c.eth.BalanceAt")
	}
	return balance, nil
}

func (c *Client) GetTransaction(ctx context.Context, txAddress string) (uint64, error) {
	hash := common.HexToHash(txAddress)
	receipt, err := c.eth.TransactionReceipt(ctx, hash)
	if err != nil {
		return 0, errors.Wrap(err, "c.eth.GetTransaction")
	}
	return receipt.Status, nil
}

func (c *Client) PendingNonceAt(ctx context.Context, address common.Address) (uint64, error) {
	return c.eth.PendingNonceAt(ctx, address)
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.eth.SuggestGasPrice(ctx)
}

func (c *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.eth.NetworkID(ctx)
}

func (c *Client) GetEstimatedGasPrice(preference string) (*big.Int, error) {
	minGasPrice := big.NewInt(int64(1e9)) // 1 GWei
	gasPrice, err := bridge.GetGasPriceFromUpvest(preference)

	if err != nil {
		gasPrice, err = c.SuggestGasPrice(context.Background())
		if err != nil {
			gasPrice = big.NewInt(0)
		}
	}

	if gasPrice.Cmp(minGasPrice) < 0 {
		gasPrice = minGasPrice
	}

	return gasPrice, nil
}

func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.eth.SendTransaction(ctx, tx)
}

func (c *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.eth.TransactionReceipt(ctx, txHash)
}

func (c *Client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return c.eth.BlockByHash(ctx, hash)
}

func (c *Client) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return c.eth.TransactionByHash(ctx, hash)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.eth.BlockByNumber(ctx, number)
}

// TransactionsByBlockNumber returns all transactions from the current canonical chain. If number is nil, the
// latest known block is returned.
func (c *Client) TransactionsByBlockNumber(ctx context.Context, number *big.Int) (types.Transactions, error) {
	block, err := c.eth.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}

const ADDRESS_0 = "0x0000000000000000000000000000000000000000"

func (c *Client) GetProof(txHash common.Hash) (*big.Int, string, uint, []string, error) {

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	chainID, err := c.GetClient().ChainID(ctx)

	if err != nil {
		return nil, "", 0, nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	// Get tx content
	txReceipt, err := c.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, "", 0, nil, errors.Wrap(err, "c.TransactionReceipt")
	}

	// time.Sleep(time.Second * 1)

	// return error when:
	// 1. currentReceipt.Status != types.ReceiptStatusSuccessful
	// 2. currentReceipt.BlockNumber + 15 <= result returns from api eth_blockNumber

	blockHashStr := txReceipt.BlockHash.String()
	txIndex := txReceipt.TransactionIndex
	blockNumber := txReceipt.BlockNumber

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	blockHeader, err := c.BlockByHash(ctx, common.HexToHash(blockHashStr))
	if err != nil {
		return nil, "", 0, nil, errors.Wrap(err, "c.BlockByHash")
	}

	siblingTxs := blockHeader.Transactions()

	// Constructing the receipt trie (source: go-ethereum/core/types/derive_sha.go)
	keybuf := new(bytes.Buffer)
	receiptTrie := new(trie.Trie)
	receipts := make([]*types.Receipt, 0)

	for i, tx := range siblingTxs {

		log.Println("sleep 300ms...")
		time.Sleep(300 * time.Millisecond)

		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
		siblingReceipt, err := c.TransactionReceipt(ctx, tx.Hash())

		if err != nil {

			log.Println("get proof err at TransactionReceipt: ", err)

			// if polygon, ignore the last tx if not found!
			if (chainID.Uint64() == 137 || chainID.Uint64() == 80001) && i == len(siblingTxs)-1 {
				continue
			}
			return nil, "", 0, nil, errors.Wrap(err, "can't get tx vs TransactionReceipt: "+tx.Hash().String()+" at index "+fmt.Sprintf("%v", i))
		}
		if i == len(siblingTxs)-1 {
			ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
			txInfo, _, err := c.TransactionByHash(ctx, tx.Hash())

			if err != nil {
				return nil, "", 0, nil, err

			}

			s := types.NewLondonSigner(chainID)
			from, err := s.Sender(txInfo)

			if err != nil {
				return nil, "", 0, nil, err
			}

			// if txinfo.To() == nil {
			// 	return nil, "", 0, nil, errors.New("to nil")
			// }

			// log.Println("txinfo.To(): ", txinfo.To())
			// log.Println("from: ", from)

			// if txinfo.To().String() == ADDRESS_0 && from.String() == ADDRESS_0 {
			// 	break
			// }

			if from.String() == ADDRESS_0 {
				break
			}
		}
		receipts = append(receipts, siblingReceipt)
	}

	receiptList := types.Receipts(receipts)
	receiptTrie.Reset()

	valueBuf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(valueBuf)

	// StackTrie requires values to be inserted in increasing hash order, which is not the
	// order that `list` provides hashes in. This insertion sequence ensures that the
	// order is correct.
	var indexBuf []byte
	for i := 1; i < receiptList.Len() && i <= 0x7f; i++ {
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		value := c.encodeForDerive(receiptList, i, valueBuf)
		receiptTrie.Update(indexBuf, value)
	}
	if receiptList.Len() > 0 {
		indexBuf = rlp.AppendUint64(indexBuf[:0], 0)
		value := c.encodeForDerive(receiptList, 0, valueBuf)
		receiptTrie.Update(indexBuf, value)
	}
	for i := 0x80; i < receiptList.Len(); i++ {
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		value := c.encodeForDerive(receiptList, i, valueBuf)
		receiptTrie.Update(indexBuf, value)
	}

	trieHash := receiptTrie.Hash()
	fmt.Println("trieHash", trieHash)

	ethReceiptHash := types.DeriveSha(receiptList, trie.NewStackTrie(nil))
	fmt.Println("ethReceiptHash", ethReceiptHash)

	// Constructing the proof for the current receipt (source: go-ethereum/trie/proof.go)
	proof := light.NewNodeSet()
	keybuf.Reset()
	rlp.Encode(keybuf, uint(txIndex))
	err = receiptTrie.Prove(keybuf.Bytes(), 0, proof)
	if err != nil {
		return nil, "", 0, nil, err
	}

	nodeList := proof.NodeList()
	encNodeList := make([]string, 0)
	for _, node := range nodeList {
		str := base64.StdEncoding.EncodeToString(node)
		encNodeList = append(encNodeList, str)
	}

	ethHeader, err := c.getEthHeader(common.HexToHash(blockHashStr))
	if err != nil {
		fmt.Printf("Header Error +%v \n", err)
		return nil, "", 0, nil, err
	}

	val, err := trie.VerifyProof(ethHeader.ReceiptHash, keybuf.Bytes(), proof)
	if err != nil {
		fmt.Printf("WARNING: ETH proof verification failed: %v", err)
		return nil, "", 0, nil, err
	}
	fmt.Printf("Verify result: +%v", val)

	return blockNumber, blockHashStr, uint(txIndex), encNodeList, nil
}

func (c *Client) getEthHeader(
	blockHash common.Hash,
) (*types.Header, error) {
	blockByHash, err := c.BlockByHash(context.Background(), blockHash)
	if err != nil {
		return nil, errors.Wrap(err, "c.BlockByHash")
	}

	blockByNumber, err := c.BlockByNumber(context.Background(), blockByHash.Number())
	if err != nil {
		return nil, errors.Wrap(err, "c.BlockByHash")
	}

	if blockByNumber.Hash().String() != blockByHash.Hash().String() {
		return nil, errors.New("the requested eth BlockHash is being on fork branch, rejected")
	}

	return blockByHash.Header(), nil
}

func (c *Client) encodeForDerive(list types.DerivableList, i int, buf *bytes.Buffer) []byte {
	buf.Reset()
	list.EncodeIndex(i, buf)
	// It's really unfortunate that we need to do perform this copy.
	// StackTrie holds onto the values until Hash is called, so the values
	// written to it must not alias.
	return common.CopyBytes(buf.Bytes())
}

func (c *Client) GetNonceByPrivateKey(senderPrivKey string) (uint64, error) {
	privateKey, err := crypto.HexToECDSA(senderPrivKey)
	if err != nil {
		return 0, errors.Wrap(err, "crypto.HexToECDSA")
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return 0, errors.Wrap(err, "s.ethClient.PendingNonceAt")
	}

	return nonce, nil
}

func (c *Client) GetNonce(txs []string, ctx context.Context) (*big.Int, error) {
	if len(txs) == 0 {
		return nil, nil // No tx => not retry => nonce = nil to auto estimate
	}

	for _, tx := range txs {
		t, _, err := c.TransactionByHash(ctx, common.HexToHash(tx))
		if err != nil {
			continue
		}
		return big.NewInt(int64(t.Nonce())), nil
	}
	return nil, fmt.Errorf("failed getting nonce %v", txs)
}

func (c *Client) GetNonceByTx(tx string, ctx context.Context) (*big.Int, error) {
	if len(tx) == 0 {
		return nil, nil // No tx => not retry => nonce = nil to auto estimate
	}

	t, _, err := c.TransactionByHash(ctx, common.HexToHash(tx))
	if err != nil {
		return nil, err
	}

	return big.NewInt(int64(t.Nonce())), nil
}

func (c *Client) GetMaxGasPrice(txs []string) (*big.Int, error) {
	if len(txs) == 0 {
		return big.NewInt(0), nil
	}

	maxGasPrice := big.NewInt(0)
	for _, tx := range txs {
		t, _, err := c.TransactionByHash(context.Background(), common.HexToHash(tx))
		if err != nil {
			continue
		}
		p := t.GasPrice()
		if p.Cmp(maxGasPrice) > 0 {
			maxGasPrice = p
		}
	}
	return maxGasPrice, nil
}

func (c *Client) ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}

// transfer:
func (c *Client) Transfer(senderPrivKey, receiverAddress string, amount *big.Int) (string, error) {
	privateKey, err := crypto.HexToECDSA(senderPrivKey)
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", errors.Wrap(err, "s.ethClient.PendingNonceAt")
	}

	gasLimit := uint64(100000)
	gasPrice, err := c.SuggestGasPrice(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "s.ethClient.SuggestGasPrice")
	}

	fee := new(big.Int)
	fee.Mul(big.NewInt(int64(gasLimit)), gasPrice)

	value := new(big.Int)
	value = amount

	toAddress := common.HexToAddress(receiverAddress)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := c.NetworkID(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "c.NetworkID")
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", errors.Wrap(err, "types.SignTx")
	}
	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", errors.Wrap(err, "c.SendTransaction")
	}
	return signedTx.Hash().Hex(), nil
}

func (c *Client) TransferWithNoneNumber(
	senderPrivKey,
	receiverAddress string,
	amount *big.Int,
	gasPrice *big.Int,
	nonce uint64) (string, error) {
	privateKey, err := crypto.HexToECDSA(senderPrivKey)
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}

	if gasPrice.Uint64() <= 0 {
		return "", errors.Wrap(err, "GasPrice is empty")
	}

	gasLimit := uint64(100000)

	fee := new(big.Int)
	fee.Mul(big.NewInt(int64(gasLimit)), gasPrice)

	value := new(big.Int)
	value = amount

	toAddress := common.HexToAddress(receiverAddress)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := c.NetworkID(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "c.NetworkID")
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", errors.Wrap(err, "types.SignTx")
	}

	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", errors.Wrap(err, "c.SendTransaction")
	}
	return signedTx.Hash().Hex(), nil
}

// transfer:
func (c *Client) TransferFastest(senderPrivKey, receiverAddress string, amount *big.Int) (string, error) {
	privateKey, err := crypto.HexToECDSA(senderPrivKey)
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", errors.Wrap(err, "s.ethClient.PendingNonceAt")
	}

	gasLimit := uint64(100000)
	logger, err := zap.NewProduction()
	if err != nil {
		return "", errors.Wrap(err, "zap.NewProduction")
	}
	gasPrice := bridge.GetSafeGasPrice(c.eth, logger)
	//gasPrice, err := c.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return "", errors.Wrap(err, "s.ethClient.SuggestGasPrice")
	//}

	fee := new(big.Int)
	fee.Mul(big.NewInt(int64(gasLimit)), gasPrice)

	value := new(big.Int)
	value = amount

	toAddress := common.HexToAddress(receiverAddress)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := c.NetworkID(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "c.NetworkID")
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", errors.Wrap(err, "types.SignTx")
	}
	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", errors.Wrap(err, "c.SendTransaction")
	}
	return signedTx.Hash().Hex(), nil
}

func (c *Client) TransferToken(senderPrivKey, receiverAddress, tokenContract string, amount *big.Int) (string, error) {
	privateKey, err := crypto.HexToECDSA(senderPrivKey)
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", errors.Wrap(err, "s.ethClient.PendingNonceAt")
	}

	value := big.NewInt(0) // eth amount
	logger, err := zap.NewProduction()
	if err != nil {
		return "", errors.Wrap(err, "zap.NewProduction")
	}
	gasPrice := bridge.GetSafeGasPrice(c.eth, logger)
	//gasPrice, err := c.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return "", errors.Wrap(err, "s.ethClient.SuggestGasPrice")
	//}

	toAddress := common.HexToAddress(receiverAddress)
	tokenAddress := common.HexToAddress(tokenContract)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit := uint64(100000)
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := c.NetworkID(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}

	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", errors.Wrap(err, "crypto.HexToECDSA")
	}

	return signedTx.Hash().Hex(), nil
}

func (c *Client) GetInfoErc20(txID string) (*big.Int, string, error) {
	tx, _, err := c.TransactionByHash(context.Background(), common.HexToHash(txID))
	if err != nil {
		return nil, "", errors.Wrap(err, "c.TransactionByHash")
	}
	amount := tx.Data()[len(tx.Data())-32:]
	hexAmount := common.Bytes2Hex(amount)
	intAmount := new(big.Int)
	intAmount.SetString(hexAmount, 16)

	address := tx.Data()[4 : len(tx.Data())-32]
	hexAddress := common.Bytes2Hex(address)
	add := common.HexToAddress(hexAddress)

	return intAmount, add.Hex(), nil
}
