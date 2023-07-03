package token_explorer

import "rederinghub.io/utils/helpers"

type Token struct {
	Address         string `json:"address"`
	Symbol          string `json:"symbol"`
	Decimal         int    `json:"decimal"`
	Name            string `json:"name"`
	TotalSupply     string `json:"total_supply"`
	Owner           string `json:"owner"`
	DeployedAtBlock int    `json:"deployed_at_block"`
}

type WalletAddressToken struct {
	Symbol           string `json:"symbol"`
	Decimal          int    `json:"decimal"`
	Name             string `json:"name"`
	Contract         string `json:"contract"`
	Balance          string `json:"balance"`
	UpdatededAtBlock int    `json:"updated_at_block"`
}

type Response struct {
	Code   string      `json:"code"`
	Error  error       `json:"error"`
	Result interface{} `json:"result"`
}

type SearchToken struct {
	Totals int     `json:"totals"`
	Tokens []Token `json:"tokens"`
}

func (r *Response) ToSearchTokens() (*SearchToken, error) {
	var resp SearchToken
	err := helpers.JsonTransform(r.Result, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *Response) ToTokens() ([]Token, error) {
	var resp []Token
	err := helpers.JsonTransform(r.Result, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Response) ToWalletAddressTokens() ([]WalletAddressToken, error) {
	var resp []WalletAddressToken
	err := helpers.JsonTransform(r.Result, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Response) ToToken() (*Token, error) {
	var resp Token
	err := helpers.JsonTransform(r.Result, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
