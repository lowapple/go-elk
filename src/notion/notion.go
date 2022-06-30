package notion

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/lowapple/go-elk/src/common/config"
	"github.com/lowapple/go-elk/src/common/downloader"
	"github.com/lowapple/go-elk/src/utils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"log"
	"strings"
)

type Notion struct {
}

func New() *Notion {
	return &Notion{}
}

func (e *Notion) Extract(target string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("window-size", "1920,200000"),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	log.Printf("Parsing page %s", target)
	log.Print("Waiting page loaded...")
	err := chromedp.Run(ctx, chromedp.Navigate(target))
	if err != nil {
		return nil
	}
	var htmlContent string
	// 페이지 로딩
	if err := processNotionPageLoaded(ctx, target); err != nil {
		return nil
	}
	// HTML 파일 가져오기
	if err := processNotionPageHtmlContent(ctx, &htmlContent); err != nil {
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	notionDownloader := downloader.New(true)
	// sheets 추출
	if err := processSheets(notionDownloader, doc); err != nil {
		return err
	}
	if err := processImages(notionDownloader, doc); err != nil {
		return err
	}
	// inject sheets
	injectSheets := `
.notion-frame{
	width: 100% !important
}`
	filePath := fmt.Sprintf(`%s/elk.css`, config.OutputPath)
	file := files.CreateFile(filePath)
	err = files.WriteFile(&injectSheets, file)
	if err != nil {
		return err
	}
	// Add Header
	var nodeAttrs []html.Attribute
	nodeAttrs = append(nodeAttrs, html.Attribute{Key: `rel`, Val: `stylesheet`})
	nodeAttrs = append(nodeAttrs, html.Attribute{Key: `href`, Val: `/elk.css`})
	node := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Link,
		Data:     "link",
		Attr:     nodeAttrs,
	}
	doc.Find(`head`).AppendNodes(node)
	// Output index.html
	filePath = fmt.Sprintf(`%s/index.html`, config.OutputPath)
	file = files.CreateFile(filePath)
	htmlContent, err = doc.Html()
	if err != nil {
		return err
	}
	err = files.WriteFile(&htmlContent, file)
	if err != nil {
		return err
	}
	return nil
}

func processImages(downloader *downloader.Struct, doc *goquery.Document) error {
	links := doc.Find(`img`)
	links.Each(func(i int, link *goquery.Selection) {
		src, ok := link.Attr("src")
		if strings.Contains(src, "data:image") {
			return
		}
		if ok {
			// 이미지 주소 자체가 파일로 저장하기 적합하지 않기 때문에 이름을 변경해서 저장할 수 있도록 한다
			path, err := downloader.CacheFile(fmt.Sprintf(`https://www.notion.so%s`, src), config.OutputPath, "jpg")
			if err != nil {
				return
			}
			link.SetAttr("src", *path)
		}
	})
	return nil
}

func processSheets(downloader *downloader.Struct, doc *goquery.Document) error {
	links := doc.Find(`link[rel='stylesheet']`)
	links.Each(func(_ int, link *goquery.Selection) {
		href, ok := link.Attr("href")
		if ok {
			if strings.HasPrefix(href, "/") {
				if strings.Contains(href, "vendors~") {
					return
				}
			}
			path, err := downloader.CacheFile(fmt.Sprintf(`https://www.notion.so%s`, href), config.OutputPath, "css")
			if err != nil {
				return
			}
			link.SetAttr("href", *path)
		}
	})
	return nil
}

// 노션 페이지 로딩
//
// Return:
//	- error: 에러 발생
func processNotionPageLoaded(ctx context.Context, URL string) error {
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
				//log.Printf("Waiting for page content to load "+
				//	"pending blocks: %d ,"+
				//	"loading blocks: %d ,"+
				//	"loading scrollers: %d/%d", pendingLen, loadingSpinnersLen, loadedScrollersLen, scrollersLen)
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

func processNotionPageHtmlContent(ctx context.Context, htmlContent *string) error {
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
