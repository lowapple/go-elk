package downloader

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/lowapple/elk/src/common/config"
	"github.com/lowapple/elk/src/common/request"
	"github.com/lowapple/elk/src/static"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type Struct struct {
	// Http Client
	client *http.Client
	// Output File path
	filePath string
	// Output File path (temp)
	tmpFilePath string
	// ProgressBar
	progressBar *progressbar.ProgressBar
	// ProgressBar Enable
	progressBarEnable bool
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func New(progressBarEnable bool) *Struct {
	return &Struct{
		client:            request.DefaultClient(),
		filePath:          config.OutputPath,
		progressBarEnable: progressBarEnable,
	}
}

func (downloader *Struct) CacheFile(URL string, outputPath string, fileType string) (*string, error) {
	// id := strings.Replace(uuid.New().String(), "-", "", -1)
	id := []byte(URL)
	h := sha1.New()
	h.Write(id)
	idStr := hex.EncodeToString(h.Sum(nil))

	fileId := fmt.Sprintf(`%s.%s`, idStr, fileType)
	filePath := fmt.Sprintf(`%s/%s`, outputPath, fileId)
	err := downloader.save(static.URL{URL: URL}, filePath)
	if err != nil {
		return nil, err
	}
	return &fileId, nil
}

func (downloader *Struct) Download(URL string) error {
	// create folder or skipped

	fileNames := strings.Split(URL, "/")
	fileName := fmt.Sprintf(`%s/%s`, config.OutputPath, fileNames[len(fileNames)-1])
	return downloader.save(static.URL{URL: URL}, fileName)
}

func (downloader *Struct) save(url static.URL, fileURI string) error {
	openOpts := os.O_RDWR | os.O_CREATE
	file, err := os.OpenFile(fileURI, openOpts, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		log.Printf(`file "%s" already exists and will be skipped`, fileURI)
		return nil
	}

	return downloader.writeFile(url.URL, file)
}

func (downloader *Struct) writeFile(URL string, file *os.File) error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	if ref := req.Header.Get("Referer"); ref == "" {
		req.Header.Set("Referer", URL)
	}
	res, err := downloader.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		time.Sleep(1 * time.Second)
		res, _ = downloader.client.Get(URL)
	}
	defer res.Body.Close()

	var ws io.Writer
	ws = file
	downloader.initProgressBar(res.ContentLength, fmt.Sprintf("Downloading %s ...", file.Name()), true)
	if downloader.progressBarEnable {
		ws = io.MultiWriter(file, downloader.progressBar)
	}
	_, copyErr := io.Copy(ws, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return fmt.Errorf("file copy error: %s", copyErr)
	}
	return nil
}

func (downloader *Struct) initProgressBar(len int64, description string, asBytes bool) {
	if !downloader.progressBarEnable {
		return
	}
	if asBytes {
		downloader.progressBar = progressbar.NewOptions(
			int(len),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetDescription(description),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
		return
	}
	downloader.progressBar = progressbar.NewOptions(
		int(len),
		progressbar.OptionShowIts(),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
	)
}
