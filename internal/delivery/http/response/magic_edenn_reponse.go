package response

import "rederinghub.io/internal/entity"

type MagicEdenOrdinalMeta struct {
	Name          string                   `json:"name"`
	HighResImgUrl string                   `json:"high_res_img_url"`
	Attributes    []entity.TokenUriAttrStr `json:"attributes"`
}
type MagicEdenOrdinalResponse struct {
	ID   string                `json:"id"`
	Meta *MagicEdenOrdinalMeta `json:"meta"`
}
