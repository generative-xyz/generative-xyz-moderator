package model

import (
	"fmt"
	"rederinghub.io/api"
	"rederinghub.io/internal/dto"
)

type OpenSeaAttribute struct {
	TraitType string
	Value     string
}

type RenderedNft struct {
	BaseModel
	ChainID         string               `json:"chainId,omitempty"  bson:"chainId,omitempty"`
	ContractAddress string               `json:"contractAddress,omitempty"  bson:"contractAddress,omitempty"`
	ProjectID       string               `json:"projectId,omitempty"  bson:"projectId,omitempty"`
	TokenID         string               `json:"tokenId,omitempty"  bson:"tokenId,omitempty"`
	Image           *string              `json:"image,omitempty"  bson:"image,omitempty"`
	Glb             *string              `json:"glb,omitempty"  bson:"glb,omitempty"`
	Video           *string              `json:"video,omitempty"  bson:"video,omitempty"`
	Name            string               `json:"name,omitempty"  bson:"name,omitempty"`
	Description     *string              `json:"description,omitempty"  bson:"description,omitempty"`
	ExternalLink    *string              `json:"externalLink,omitempty"  bson:"externalLink,omitempty"`
	Attributes      []*OpenSeaAttribute  `json:"attributes,omitempty"  bson:"attributes,omitempty"`
	EmotionTime     string               `json:"emotionTime,omitempty" bson:"emotionTime,omitempty"`
	Metadata        *RenderedNftMetadata `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Gif             *string              `json:"gif,omitempty"  bson:"gif,omitempty"`
}

type RenderedNftMetadata struct {
	BackgroundColor *string `json:"backgroundColor,omitempty" bson:"backgroundColor,omitempty"`
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

func (r *RenderedNft) ToGenerativeProto(template dto.TemplateDTO) *api.GetGenerativeNFTMetadataResponse {
	resp := &api.GetGenerativeNFTMetadataResponse{
		Name:         r.Name,
		Description:  r.Description,
		Image:        *r.Image,
		AnimationUrl: *r.Glb,
		GlbUrl:       *r.Glb,
		ExternalLink: *r.ExternalLink,
		Attributes:   OpenSeaAttributeSlice(r.Attributes).ToProto(),
	}

	switch template.BlenderType {
	case "confetti", "horn":
		resp.AnimationUrl = ""
	}
	return resp
}

func (r *RenderedNft) ToCandyResponse() *api.GetCandyMetadataResponse {
	return &api.GetCandyMetadataResponse{
		Name:         r.Name,
		Description:  r.Description,
		Image:        *r.Image,
		AnimationUrl: *r.Glb,
		ExternalLink: *r.ExternalLink,
		Attributes:   OpenSeaAttributeSlice(r.Attributes).ToProto(),
	}
}

func (r *RenderedNft) ToAvatarResponse() *api.GetAvatarMetadataResponse {
	gif := ""
	if r.Video != nil {
		gif = *r.Gif
	}
	resp := &api.GetAvatarMetadataResponse{
		Name:         r.Name,
		Description:  r.Description,
		Image:        *r.Image,
		AnimationUrl: gif,
		ExternalLink: *r.ExternalLink,
		Attributes:   OpenSeaAttributeSlice(r.Attributes).ToProto(),
		GlbUrl:       *r.Glb,
	}

	if r.Metadata != nil {
		metadata := r.Metadata
		if metadata.BackgroundColor != nil {
			resp.BackgroundColor = metadata.BackgroundColor
		}
	}
	return resp
}

func (m RenderedNft) CollectionName() string {
	return "rendered_nfts"
}

func (r *RenderedNft) ToRenderingRepsonse() *api.GetGenerativeNFTMetadataResponse {
	return &api.GetGenerativeNFTMetadataResponse{
		Name:        fmt.Sprintf("Rendering on #%s", r.TokenID),
		Image:       "https://cdn.rove.to/metaverse/rove/Rove_logo.png",
		Description: r.Description,
	}
}
