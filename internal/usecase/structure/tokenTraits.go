package structure

type TokenTraits struct {
	ID string `json:"id"`
	Atrributes []TraitAttribute `json:"atrributes"`
}

type TraitAttribute struct {
	TraitType string `json:"trait_type"`
	Value string `json:"value"`
}