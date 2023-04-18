package http

import (
	"context"
	"net/http"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	_req "rederinghub.io/utils/request"
)

// Capture
// @Summary captures url content as image
// @Description captures url content as image
// @Tags Photographer
// @Accept  json
// @Produce  json
// @Param json body request.CaptureRequest true "capture request"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.CaptureResponse}
// @Router /photo/capture [POST]
func (h *httpDelivery) Capture(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var req request.CaptureRequest
			err := _req.BindJson(r, &req)
			if err != nil {
				return nil, err
			}
			url, err := h.Usecase.CaptureContent(req.ID, req.Url)
			if err != nil {
				return nil, err
			}
			return response.CaptureResponse{
				ImageUrl: url,
				ID:       req.ID,
			}, nil

		},
	).ServeHTTP(w, r)
}
