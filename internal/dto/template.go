package dto

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"rederinghub.io/api"
)

type TemplateDTO struct {
	TokenID            string                          `json:"tokenId" bson:"tokenId"`
	FeeStr             string                          `json:"fee" bson:"fee"`
	FeeToken           common.Address                  `json:"feeToken" bson:"feeToken"`
	MintMaxSupplyStr   string                          `json:"mintMaxSupply" bson:"mintMaxSupply"`
	MintTotalSupplyStr string                          `json:"mintTotalSupply" bson:"mintTotalSupply"`
	Script             string                          `json:"script" bson:"script"`
	ScriptType         uint32                          `json:"scriptType" bson:"scriptType"`
	Creator            common.Address                  `json:"creator" bson:"creator"`
	CustomUri          string                          `json:"customUri" bson:"customUri"`
	ProjectName        string                          `json:"projectName" bson:"projectName"`
	ClientSeed         bool                            `json:"clientSeed" bson:"clientSeed"`
	ParamsTemplate     BoilerplateParamParamsOfProject `json:"paramsTemplate" bson:"paramsTemplate"`
}

func (d *TemplateDTO) ToProto() *api.GetTemplateDetailResponse {
	return &api.GetTemplateDetailResponse{
		TokenId:         d.TokenID,
		Fee:             d.FeeStr,
		FeeToken:        d.FeeToken.String(),
		MintMaxSupply:   d.MintMaxSupplyStr,
		MintTotalSupply: d.MintTotalSupplyStr,
		Script:          d.Script,
		ScriptType:      d.ScriptType,
		Creator:         d.Creator.String(),
		CustomUri:       d.CustomUri,
		ProjectName:     d.ProjectName,
		ClientSeed:      d.ClientSeed,
		ParamsTemplate: func(d *TemplateDTO) *api.ParamsTemplate {
			params := make([]*api.Param, 0, len(d.ParamsTemplate.Params))
			for _, item := range d.ParamsTemplate.Params {
				params = append(params, &api.Param{
					TypeValue:       uint32(item.TypeValue),
					Max:             item.MaxStr,
					Min:             item.MinStr,
					Decimal:         uint32(item.Decimal),
					AvailableValues: item.AvailableValues,
					Value:           item.ValueStr,
					Editable:        item.Editable,
				})
			}
			return &api.ParamsTemplate{
				Seed:   d.ParamsTemplate.SeedStr,
				Params: params,
			}
		}(d),
	}
}

func (d *TemplateDTO) Fee(v *big.Int) {
	d.FeeStr = v.String()
}

func (d *TemplateDTO) MintMaxSupply(v *big.Int) {
	d.MintMaxSupplyStr = v.String()
}

func (d *TemplateDTO) MintTotalSupply(v *big.Int) {
	d.MintTotalSupplyStr = v.String()
}

type BoilerplateParamParamsOfProject struct {
	SeedStr string                          `json:"seed" bson:"seed"`
	Params  []BoilerplateParamParamTemplate `json:"params" bson:"params"`
}

func (d *BoilerplateParamParamsOfProject) Seed(seed [32]byte) {
	h := common.BytesToHash(seed[0:])
	d.SeedStr = h.String()
}

type BoilerplateParamParamTemplate struct {
	TypeValue       uint8    `json:"typeValue"`
	MaxStr          string   `json:"max" bson:"max"`
	MinStr          string   `json:"min" bson:"min"`
	Decimal         uint8    `json:"decimal"`
	AvailableValues []string `json:"availableValues" bson:"availableValues"`
	ValueStr        string   `json:"value" bson:"value"`
	Editable        bool     `json:"editable" bson:"editable"`
}

func (d *BoilerplateParamParamTemplate) Value(v *big.Int) {
	d.ValueStr = v.String()
}

func (d *BoilerplateParamParamTemplate) Min(v *big.Int) {
	d.MinStr = v.String()
}

func (d *BoilerplateParamParamTemplate) Max(v *big.Int) {
	d.MaxStr = v.String()
}
