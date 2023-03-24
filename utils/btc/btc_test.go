package btc

import (
	"log"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"gotest.tools/assert"
)

func TestBtc(t *testing.T) {

	bs := NewBlockcypherService("https://api.blockcypher.com/v1/btc/main/addrs", "", "", &chaincfg.MainNetParams)

	addresses := []string{"bc1qn74ftxrvh862jcre972ulnvmve9ek50ewngwyx", "bc1qc7gal7g4snw2lxyeuatrjwurfw0skmy6q2jj9z"}
	objectAddress, err := bs.BTCGetAddrInfoMulti(addresses)

	log.Println(objectAddress, err)

	for _, v := range objectAddress {
		log.Println("address: ", v.Address, "==> balance: ", v.Balance)
	}

	assert.Equal(t, false, true)
}
