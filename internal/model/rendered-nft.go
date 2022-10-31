package model

import "rederinghub.io/api"

type OpenSeaAttribute struct {
	TraitType string
	Value     string
}

type RenderedNft struct {
	BaseModel
	ChainID         string              `json:"chainId,omitempty"  bson:"chainId,omitempty"`
	ContractAddress string              `json:"contractAddress,omitempty"  bson:"contractAddress,omitempty"`
	ProjectID       string              `json:"projectId,omitempty"  bson:"projectId,omitempty"`
	TokenID         string              `json:"tokenId,omitempty"  bson:"tokenId,omitempty"`
	Image           *string             `json:"image,omitempty"  bson:"image,omitempty"`
	Glb             *string             `json:"glb,omitempty"  bson:"glb,omitempty"`
	Video           *string             `json:"video,omitempty"  bson:"video,omitempty"`
	Name            string              `json:"name,omitempty"  bson:"name,omitempty"`
	Description     *string             `json:"description,omitempty"  bson:"description,omitempty"`
	ExternalLink    *string             `json:"externalLink,omitempty"  bson:"externalLink,omitempty"`
	Attributes      []*OpenSeaAttribute `json:"attributes,omitempty"  bson:"attributes,omitempty"`
}

func (o *OpenSeaAttribute) ToProto() *api.OpenSeaAttribute {
	return &api.OpenSeaAttribute{
		TraitType: o.TraitType,
		Value:     o.Value,
	}
}

type OpenSeaAttributeSlice []*OpenSeaAttribute

func (o OpenSeaAttributeSlice) ToProto() []*api.OpenSeaAttribute {
	res := make([]*api.OpenSeaAttribute, len(o))
	for i := 0; i < len(o); i++ {
		res[i] = o[i].ToProto()
	}
	return res
}

func (r *RenderedNft) ToProto() *api.GetRenderedNftResponse {
	return &api.GetRenderedNftResponse{
		Name:         r.Name,
		Description:  r.Description,
		Image:        *r.Image,
		AnimationUrl: *r.Glb,
		ExternalLink: *r.ExternalLink,
		Attributes:   OpenSeaAttributeSlice(r.Attributes).ToProto(),
	}
}

func (m RenderedNft) CollectionName() string {
	return "rendered_nfts"
}
