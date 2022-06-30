package downloader

import (
	"github.com/lowapple/go-elk/src/common/config"
	"testing"
)

func TestDownload(t *testing.T) {
	config.OutputPath = "./blog.lowapple.io"
	downloader := New(false)
	t.Run("Stylesheet 데이터 다운로드", func(t *testing.T) {
		err := downloader.Download("https://lowapple.notion.site/e1db500567ca46bcbb59ed2f575325e4")
		if err != nil {
			t.Errorf("데이터 다운로드 실패")
		}
	})
}
