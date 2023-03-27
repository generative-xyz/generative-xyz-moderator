package structure

type TokenTraits struct {
	ID string `json:"id"`
	Atrributes []TraitAttribute `json:"attributes"`
}

type TraitAttribute struct {
	TraitType string `json:"trait_type"`
	Value string `json:"value"`
}