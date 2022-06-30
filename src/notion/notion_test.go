package notion

import (
	"github.com/lowapple/go-elk/src/common/config"
	files "github.com/lowapple/go-elk/src/utils"
	"testing"
)

func TestParseURL(t *testing.T) {
	config.OutputPath = "./blog.lowapple.io"
	err := files.IsNotExistMkDir(config.OutputPath)
	if err != nil {
		return
	}
	t.Run("notion parsing data", func(t *testing.T) {
		err = New().Extract("https://lowapple.notion.site/e1db500567ca46bcbb59ed2f575325e4")
		if err != nil {
			return
		}
	})
}
