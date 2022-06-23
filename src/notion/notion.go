package notion

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/lowapple/elk/src/static"
	"log"
	"strings"
)

type extractor struct {
}

func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	log.Printf("Parsing page %s", URL)
	log.Print("Parsing page config ...")
	err := chromedp.Run(ctx, chromedp.Navigate(URL))
	if err != nil {
		return nil, nil
	}
	var htmlContent string
	// 페이지 로딩
	if err := waitForPageLoaded(ctx, URL); err != nil {
		return nil, nil
	}
	// HTML 파일 가져오기
	if err := getHtmlContent(ctx, &htmlContent); err != nil {
		return nil, nil
	}
	log.Print(htmlContent)
	return nil, nil
}

// 노션 페이지 로딩
//
// Return:
//	- error: 에러 발생
func waitForPageLoaded(ctx context.Context, URL string) error {
	c := chromedp.FromContext(ctx)
	err := chromedp.Tasks([]chromedp.Action{
		chromedp.Navigate(URL),
		chromedp.WaitVisible(`.notion-presence-container`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				node, err := dom.GetDocument().Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					return fmt.Errorf("page not loaded")
				}
				documentHTML, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					return fmt.Errorf("page not loaded")
				}
				documentReader, err := goquery.NewDocumentFromReader(strings.NewReader(documentHTML))
				if err != nil {
					return fmt.Errorf("page not loaded")
				}
				pending := documentReader.Find(`.notion-unknown-block`).Nodes
				pendingLen := len(pending)
				loadingSpinners := documentReader.Find(`.loading-spinner`).Nodes
				loadingSpinnersLen := len(loadingSpinners)
				scrollers := documentReader.Find(`.notion-scroller`)
				scrollersLen := len(scrollers.Nodes)
				var loadedScrollers []*goquery.Selection
				idx := 0
				for idx < scrollersLen {
					scrollerWithChildren := scrollers.Find(`div`)
					scrollerWithChildrenLen := len(scrollerWithChildren.Nodes)
					if scrollerWithChildrenLen > 0 {
						loadedScrollers = append(loadedScrollers, scrollerWithChildren)
					}
					idx += 1
				}
				loadedScrollersLen := len(loadedScrollers)
				log.Printf("Waiting for page content to load "+
					"pending blocks: %d ,"+
					"loading blocks: %d ,"+
					"loading scrollers: %d/%d", pendingLen, loadingSpinnersLen, loadedScrollersLen, scrollersLen)
				if pendingLen == 0 && loadingSpinnersLen == 0 && scrollersLen == loadedScrollersLen {
					return nil
				}
			}
		}),
	}).Do(cdp.WithExecutor(ctx, c.Target))
	if err != nil {
		return err
	}
	return nil
}

func getHtmlContent(ctx context.Context, htmlContent *string) error {
	c := chromedp.FromContext(ctx)
	node, err := dom.GetDocument().Do(cdp.WithExecutor(ctx, c.Target))
	if err != nil {
		htmlContent = nil
		return fmt.Errorf("페이지에 연결할 수 없습니다")
	}
	*htmlContent, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(cdp.WithExecutor(ctx, c.Target))
	if err != nil {
		return fmt.Errorf("페이지에 연결할 수 없습니다")
	}
	return nil
}

func parsePage(URL string) {
	log.Printf("Parsing page %s", URL)
	//logger.Printf("Using page config: %s")
	load(URL)
}

// Notion page loaded
func load(URL string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(URL),
		// .notion-presence-container loaded
		chromedp.WaitVisible(`.notion-presence-container`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf(">>>> .notion-presence-container is visible")
			return nil
		}),
		// if unknown blocks

		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return nil
			}
			do, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			log.Print(do)
			return nil
		}),
	}
}

func parseURL(URL string) []string {
	//htmlString, err := request.Request()
	//if err != nil {
	//	return nil
	//}
	//fmt.Println(htmlString)
	return nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.FullScreenshot(res, quality),
	}
}
