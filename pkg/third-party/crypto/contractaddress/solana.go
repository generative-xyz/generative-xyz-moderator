package contractaddress

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptotransactionstatus"
)

type solonaImpl struct {
	RpcClient *rpc.Client
}

func (s *solonaImpl) GetNativeStrategy() Strategy {
	return s
}

func (s solonaImpl) GetTransactionReceipt(ctx context.Context, transactionID string) (*TransactionReceipt, error) {
	return &TransactionReceipt{
		Status:        cryptotransactionstatus.Success,
		TransactionID: transactionID,
		Coin:          cryptocurrency.Solana,
	}, nil
}

func (s solonaImpl) EstimateGas(req *TransferRequest) (uint64, error) {
	return 0, nil
}

func (s solonaImpl) SuggestGasPrice(context context.Context) (*big.Int, error) {
	return nil, nil
}

func (s solonaImpl) GenerateContractAddress() (*ContractAddress, error) {
	// test 1: address - HFEu22SYhrFMVKB8X4nUAzdgaQT3z1KL2rU3iz4whNnj, pri - Tj43PV6KS3AWYhGRPoxouxAzEXBysXdEFGchN6ShJXqrdvAgdTmCLmuFSx48G7jBgWiP49Lnh499zQXgz7tXvpP
	// test 2: address - BM2RaPWwdvumJ9bBuqe9BsR3NELvPQVnw6FiYco2pNZy, pri - vh5DT8RYMy65kJL5Eq7aHZJfLa6RUjqKeYWsdWqQd5AqFziUzi2HtQedZxYUbt8ARPkkbrWuWJR4db1TfJ7Xni7
	account := solana.NewWallet()
	address := account.PublicKey().String()
	privateKey := account.PrivateKey.String()
	cAddress := NewContractAddress().WithAddress(address).WithPrivateKey(privateKey)
	return cAddress, nil
}

func (s solonaImpl) CheckBalance(address string) (float64, error) {
	publicKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return 0, err
	}
	balance, err := s.RpcClient.GetBalance(
		context.TODO(),
		publicKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return 0, err
	}

	value := float64(balance.Value) / math.Pow10(cryptocurrency.SolonaRoundPlaces)
	return value, nil
}

func (s *solonaImpl) generateNativeAddress(feePlayerPrivkey string) (privKey string, pubKey string, address, txHash string, err error) {
	account := solana.NewWallet()

	privKey = account.PrivateKey.String()
	address = account.PublicKey().String()

	feePayer, err := solana.PrivateKeyFromBase58(feePlayerPrivkey) // account to create tx.
	shieldMaker, err := solana.PrivateKeyFromBase58(privKey)       // user fixed accout

	// find address:burnProofdAccounts
	shieldNativeTokenAcc, _, err := solana.FindAssociatedTokenAddress(
		shieldMaker.PublicKey(),
		solana.SolMint,
	)

	pubKey = shieldNativeTokenAcc.String()

	// check account exist
	needCreateAccount := false
	_, err = s.RpcClient.GetAccountInfo(context.TODO(), solana.MustPublicKeyFromBase58(shieldNativeTokenAcc.String()))
	if err != nil {
		if err.Error() == "not found" {
			fmt.Println("need init account")
			needCreateAccount = true
		} else {
			log.Println("GetAccountInfo err: ", err)
		}
	}

	if !needCreateAccount {
		err = nil
		return
	}

	// init account:
	recent, err := s.RpcClient.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			associatedtokenaccount.NewCreateInstruction(
				feePayer.PublicKey(),    // account fee to create tx.
				shieldMaker.PublicKey(), // owner of token.
				solana.SolMint,          // token id.
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(feePayer.PublicKey()),
	)
	signers := []solana.PrivateKey{
		feePayer,
	}
	sig, err := s.SignAndSendTx(tx, signers)
	if err != nil {
		return
	}

	txHash = sig.String()

	return
}

func (s *solonaImpl) ValidAddress(address string) bool {
	_, err := solana.PublicKeyFromBase58(address)
	return err == nil
}

func (s *solonaImpl) SignAndSendTx(tx *solana.Transaction, signers []solana.PrivateKey) (solana.Signature, error) {
	_, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		for _, candidate := range signers {
			if candidate.PublicKey().Equals(key) {
				return &candidate
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("unable to sign transaction: %v \n", err)
		return solana.Signature{}, err
	}
	// send tx
	signature, err := s.RpcClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Printf("unable to send transaction: %v \n", err)
		return solana.Signature{}, err
	}

	return signature, nil
}

func (s *solonaImpl) Transfer(req *TransferRequest) (string, error) {
	toPub, err := solana.PublicKeyFromBase58(req.ToAddress)
	if err != nil {
		return "", err
	}

	fromAddress, err := solana.PrivateKeyFromBase58(req.FromAddress.GetPrivateKey())
	if err != nil {
		return "", err
	}

	feePayer, err := solana.PrivateKeyFromBase58(req.Payer.GetPrivateKey())
	if err != nil {
		return "", err
	}

	// account to create tx.
	recent, err := s.RpcClient.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return "", err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				req.Amount,
				fromAddress.PublicKey(),
				toPub,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(feePayer.PublicKey()),
	)
	if err != nil {
		return "", err
	}

	signers := []solana.PrivateKey{
		fromAddress,
		feePayer,
	}
	sig, err := s.SignAndSendTx(tx, signers)
	if err != nil {
		return "", err
	}

	return sig.String(), nil
}

func NewSolanaChain(isMainNet bool) Strategy {
	env := rpc.DevNet_RPC
	if isMainNet {
		env = rpc.MainNetBeta_RPC
	}
	return &solonaImpl{
		RpcClient: rpc.New(env),
	}
}
