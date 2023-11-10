package usecase

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"go.uber.org/zap"
	"log"
	"os"
	"rederinghub.io/utils/logger"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
)

func (u Usecase) CaptureContent(id, url string, duration int) (string, error) {

	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		return "", err
	}

	var contextOpts = []chromedp.ContextOption{}
	contextOpts = []chromedp.ContextOption{
		chromedp.WithErrorf(log.Printf),
		chromedp.WithLogf(log.Printf),
		chromedp.WithBrowserOption(),
	}

	opts := []chromedp.ExecAllocatorOption{}
	if os.Getenv("ENV") == "mainnet" {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("google-chrome"),
			chromedp.Flag("headless", eCH),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("no-first-run", true),
		)
	} else {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("no-first-run", true),
		)
	}

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	cctx, cancel := chromedp.NewContext(allocCtx, contextOpts...)

	//avoid overlap html
	ackCtx, cancel := context.WithTimeout(cctx, time.Duration(duration)*5*time.Second)
	defer cancel()

	var buf []byte
	traits := make(map[string]interface{})
	err = chromedp.Run(ackCtx,
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second*time.Duration(duration)),
		chromedp.CaptureScreenshot(&buf),
		chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
	)
	if err != nil {
		return "", err
	}

	image := helpers.Base64Encode(buf)
	image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)
	if image != "" {
		base64Image := image
		i := strings.Index(base64Image, ",")
		if i >= 0 {
			now := time.Now().UTC().Unix()
			name := fmt.Sprintf("capture/%s-%d.png", id, now)
			base64Image = base64Image[i+1:]
			uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
			if err != nil {
				return "", err
			}
			imageURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
			return imageURL, nil
		}
	}
	return "", fmt.Errorf("capture error")
}

func (u Usecase) PublishImageData(req request.CaptureRequest) error {
	return u.PubSub.Producer(utils.PUBSUB_CAPTURE_THUMBNAIL, redis.PubSubPayload{Data: req})
}

func (u Usecase) ParseSvg(req request.ParseSvgRequest) (*string, error) {
	id := req.ID
	url := req.Url

	_, imageType := utils.IsImageURL(url)
	if !strings.Contains(imageType, "svg") && !strings.Contains(imageType, "html") {
		return &url, nil
	}

	duration := req.DelayTime
	if duration == 0 {
		duration = 1
	}

	image, err := u.CaptureContent(id, req.Url, duration)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

type ParsedHtml struct {
	Image  string            `json:"image"`
	Traits map[string]string `json:"traits"`
}

type ParsedRedirectHtml struct {
	Html string `json:"html"`
}

func (u Usecase) CaptureHtmlContent(req request.ParseSvgRequest) (*ParsedHtml, error) {
	id := req.ID
	url := req.Url
	duration := req.DelayTime

	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		logger.AtLog.Logger.Error("CaptureHtmlContent", zap.Any("req", req), zap.Error(err))
		return nil, err
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.ExecPath("google-chrome"),
		chromedp.Flag("headless", eCH),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-first-run", true),
	)

	ctx := context.Background()
	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	cctx, cancel := chromedp.NewContext(allocCtx)

	//avoid overlap html
	ackCtx, cancel := context.WithTimeout(cctx, time.Duration(duration)*5*time.Second)
	defer cancel()

	var buf []byte
	traits := make(map[string]interface{})
	err = chromedp.Run(ackCtx,
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second*time.Duration(duration)),
		chromedp.CaptureScreenshot(&buf),
		chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
	)
	if err != nil {
		logger.AtLog.Logger.Error("CaptureHtmlContent - chromedp.Run", zap.Any("req", req), zap.Error(err))
		return nil, err
	}

	image := helpers.Base64Encode(buf)
	image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)
	if image != "" {
		base64Image := image
		i := strings.Index(base64Image, ",")
		if i >= 0 {
			now := time.Now().UTC().Unix()
			name := fmt.Sprintf("capture/%s-%d.png", id, now)
			base64Image = base64Image[i+1:]
			uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
			if err != nil {
				logger.AtLog.Logger.Error("CaptureHtmlContent - UploadBaseToBucket", zap.Any("req", req), zap.Error(err))
				return nil, err
			}

			traitsResp := make(map[string]string)
			for key, item := range traits {
				traitsResp[key] = fmt.Sprintf("%v", item)
			}

			imageURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
			return &ParsedHtml{
				Image:  imageURL,
				Traits: traitsResp,
			}, nil
		}
	}

	err = fmt.Errorf("capture error")
	logger.AtLog.Logger.Error("CaptureHtmlContent - UploadBaseToBucket", zap.Any("req", req), zap.Error(err))
	return nil, err
}

func (u Usecase) elementScreenshot(sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Screenshot(sel, res, chromedp.NodeReady),
	}
}

func (u Usecase) loadHTMLFromStringActionFunc(content string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		ch := make(chan bool, 1)
		defer close(ch)

		go chromedp.ListenTarget(ctx, func(ev interface{}) {
			if _, ok := ev.(*page.EventLoadEventFired); ok {
				ch <- true
			}
		})

		frameTree, err := page.GetFrameTree().Do(ctx)
		if err != nil {
			return err
		}

		if err := page.SetDocumentContent(frameTree.Frame.ID, content).Do(ctx); err != nil {
			return err
		}

		select {
		case <-ch:
			return nil
		case <-ctx.Done():
			return context.DeadlineExceeded
		}
	})
}

func (u Usecase) CaptureHtmlContentv2(req request.ParseSvgRequest) (*ParsedHtml, error) {
	id := req.ID
	url := req.Url
	duration := req.DelayTime

	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		logger.AtLog.Logger.Error("CaptureHtmlContent", zap.Any("req", req), zap.Error(err))
		return nil, err
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.ExecPath("google-chrome"),
		chromedp.Flag("headless", eCH),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-first-run", true),
	)

	ctx := context.Background()
	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	cctx, cancel1 := chromedp.NewContext(allocCtx)
	defer cancel1()

	//avoid overlap html
	ackCtx, cancel2 := context.WithTimeout(cctx, time.Duration(duration)*50*time.Second)
	defer cancel2()

	var buf []byte
	traits := make(map[string]interface{})

	w := 650
	h := 500
	if req.Width != 0 {
		w = req.Width
	}

	if req.Height != 0 {
		h = req.Height
	}

	actions := []chromedp.Action{
		chromedp.EmulateViewport(int64(w), int64(h)),
	}

	actions = append(actions, chromedp.Navigate(url))
	actions = append(actions, chromedp.Sleep(time.Second*time.Duration(duration)))
	actions = append(actions, chromedp.CaptureScreenshot(&buf))

	err = chromedp.Run(ackCtx, actions...)
	if err != nil {
		logger.AtLog.Logger.Error("CaptureHtmlContent - chromedp.Run", zap.Any("req", req), zap.Error(err))
		return nil, err
	}

	image := helpers.Base64Encode(buf)
	image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)
	if image != "" {
		base64Image := image
		i := strings.Index(base64Image, ",")
		if i >= 0 {
			now := time.Now().UTC().Unix()
			name := fmt.Sprintf("capture/%s-%d.png", id, now)
			base64Image = base64Image[i+1:]
			uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
			if err != nil {
				logger.AtLog.Logger.Error("CaptureHtmlContent - UploadBaseToBucket", zap.Any("req", req), zap.Error(err))
				return nil, err
			}

			traitsResp := make(map[string]string)
			for key, item := range traits {
				traitsResp[key] = fmt.Sprintf("%v", item)
			}

			imageURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
			return &ParsedHtml{
				Image:  imageURL,
				Traits: traitsResp,
			}, nil
		}
	}

	err = fmt.Errorf("capture error")
	logger.AtLog.Logger.Error("CaptureHtmlContent - UploadBaseToBucket", zap.Any("req", req), zap.Error(err))
	return nil, err
}

func (u Usecase) OpenUrl(req request.ParseSvgRequest) (*ParsedRedirectHtml, error) {
	id := req.ID
	url := req.Url
	duration := req.DelayTime

	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		logger.AtLog.Logger.Error("OpenUrl",
			zap.Any("req", req),
			zap.Error(err),
			zap.String("app", id),
		)
		return nil, err
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.ExecPath("google-chrome"),
		chromedp.Flag("headless", eCH),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-first-run", true),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	cctx, cancel := chromedp.NewContext(allocCtx)

	//avoid overlap html
	ackCtx, cancel := context.WithTimeout(cctx, time.Duration(duration)*5*time.Second)
	defer cancel()
	var body string

	err = chromedp.Run(ackCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second*time.Duration(duration)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return fmt.Errorf("could not get doc: %w", err)
			}

			body, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return fmt.Errorf("could not get html: %w", err)
			}
			return nil
		}),
	)

	if err != nil {
		logger.AtLog.Logger.Error("OpenUrl",
			zap.Any("req", req),
			zap.Error(err),
			zap.String("app", id),
		)
		return nil, err
	}

	return &ParsedRedirectHtml{
		Html: body,
	}, err
}
