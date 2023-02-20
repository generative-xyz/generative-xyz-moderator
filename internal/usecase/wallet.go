package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) GetBTCWalletInfo(address string) (*structure.WalletInfo, error) {
	var result structure.WalletInfo
	apiToken := u.Config.BlockcypherAPI
	walletBasicInfo, err := getWalletInfo(address, apiToken)
	if err != nil {
		return nil, err
	}

	result.BlockCypherWalletInfo = *walletBasicInfo
	outcoins := []string{}
	for _, outcoin := range result.Txrefs {
		o := fmt.Sprintf("%s:%v", outcoin.TxHash, outcoin.TxOutputN)
		outcoins = append(outcoins, o)
	}
	inscriptions, err := u.InscriptionsByOutputs(outcoins)
	if err != nil {
		return nil, err
	}
	result.Inscriptions = inscriptions
	return &result, nil
}

func (u Usecase) InscriptionsByOutputs(outputs []string) (map[string]string, error) {
	result := make(map[string]string)
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	for _, output := range outputs {
		inscription, err := getInscriptionByOutput(ordServer, output)
		if err != nil {
			return nil, err
		}
		if len(inscription.Inscriptions) > 0 {
			result[output] = inscription.Inscriptions[0]
		}
	}
	return result, nil
}

func getInscriptionByOutput(ordServer, output string) (*structure.InscriptionOrdInfo, error) {
	url := fmt.Sprintf("%s/api/output/%s", ordServer, output)
	fmt.Println("url", url)
	var result structure.InscriptionOrdInfo
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getInscriptionByOutput Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func getWalletInfo(address string, apiToken string) (*structure.BlockCypherWalletInfo, error) {
	url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/addrs/%s?unspentOnly=true&includeScript=false&token=%s", address, apiToken)
	fmt.Println("url", url)
	var result structure.BlockCypherWalletInfo
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getInscriptionByOutput Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}
