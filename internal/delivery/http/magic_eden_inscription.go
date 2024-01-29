package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"strings"
)

func SerializeMagicEdenResponse(arr []entity.TokenUri) []response.MagicEdenOrdinalResponse {
	var magicEdenResponse []response.MagicEdenOrdinalResponse
	for _, a := range arr {
		attrs := []entity.TokenUriAttrStr{}
		for _, a := range a.ParsedAttributesStr {
			if strings.ToLower(a.TraitType) == "hash" {
				continue
			}
			attrs = append(attrs, a)
		}
		r := response.MagicEdenOrdinalResponse{
			ID: a.TokenID,
			Meta: &response.MagicEdenOrdinalMeta{
				Name:          fmt.Sprintf("Modular #%d", a.OrderInscriptionIndex),
				Attributes:    attrs,
				HighResImgUrl: a.Thumbnail,
			},
		}
		magicEdenResponse = append(magicEdenResponse, r)
	}
	return magicEdenResponse
}
func (h *httpDelivery) GetListInscriptionWithMagicEdenFormat(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	inscriptions, err := h.Usecase.Repo.FindTokenByProjectIDWithMagicEdenMetadata(context.Background(), projectID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	res, _ := json.Marshal(SerializeMagicEdenResponse(inscriptions))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
