package usecase

import (
	"context"
	"fmt"
	"os"
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
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.ExecPath("google-chrome"),
		chromedp.Flag("headless", eCH),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-first-run", true),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	cctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte
	//traits := make(map[string]interface{})
	err = chromedp.Run(cctx,
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second*time.Duration(duration)),
		chromedp.CaptureScreenshot(&buf),
		//chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
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
