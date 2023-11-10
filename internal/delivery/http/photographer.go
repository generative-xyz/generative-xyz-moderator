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
			err = h.Usecase.PublishImageData(req)
			if err != nil {
				return nil, err
			}

			return req, nil

		},
	).ServeHTTP(w, r)
}

// Capture
// @Summary captures svg
// @Description captures svg and response png
// @Tags Photographer
// @Accept  json
// @Produce  json
// @Param json body request.ParseSvgRequest true "capture request"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.CaptureResponse}
// @Router /photo/pare-svg [POST]
func (h *httpDelivery) ParseSvg(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var req request.ParseSvgRequest
			err := _req.BindJson(r, &req)
			if err != nil {
				return nil, err
			}
			image, err := h.Usecase.ParseSvg(req)
			if err != nil {
				return nil, err
			}

			return image, nil

		},
	).ServeHTTP(w, r)
}

// Capture
// @Summary parse html
// @Description parse html
// @Tags Photographer
// @Accept  json
// @Produce  json
// @Param json body request.ParseSvgRequest true "capture request"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.CaptureResponse}
// @Router /photo/pare-html [POST]
func (h *httpDelivery) ParseHtml(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var req request.ParseSvgRequest
			err := _req.BindJson(r, &req)
			if err != nil {
				return nil, err
			}
			image, err := h.Usecase.CaptureHtmlContent(req)
			if err != nil {
				return nil, err
			}

			return image, nil

		},
	).ServeHTTP(w, r)
}

// Capture
// @Summary open url
// @Description open url
// @Tags Photographer
// @Accept  json
// @Produce  json
// @Param json body request.ParseSvgRequest true "capture request"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.CaptureResponse}
// @Router /photo/open-url [POST]
func (h *httpDelivery) OpenUrl(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var req request.ParseSvgRequest
			err := _req.BindJson(r, &req)
			if err != nil {
				return nil, err
			}
			image, err := h.Usecase.OpenUrl(req)
			if err != nil {
				return nil, err
			}

			return image, nil

		},
	).ServeHTTP(w, r)
}

// Capture
// @Summary parse html
// @Description parse html
// @Tags Photographer
// @Accept  json
// @Produce  json
// @Param json body request.ParseSvgRequest true "capture request"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.CaptureResponse}
// @Router /photo/pare-html-v2 [POST]
func (h *httpDelivery) ParseHtmlV2(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var req request.ParseSvgRequest
			err := _req.BindJson(r, &req)
			if err != nil {
				return nil, err
			}
			image, err := h.Usecase.CaptureHtmlContentv2(req)
			if err != nil {
				return nil, err
			}

			return image, nil

		},
	).ServeHTTP(w, r)
}
