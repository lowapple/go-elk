package notion

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	New().Extract("https://lowapple.notion.site/e1db500567ca46bcbb59ed2f575325e4")

	//tests := []struct {
	//	Name string
	//	Args test.Args
	//}{
	//	{
	//		Name: "Notion Parse 테스트",
	//		Args: test.Args{
	//			URL: ,
	//		},
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.Name, func(t *testing.T) {
	//		data, err := New().Extract(tt.Args.URL)
	//		test.CheckError(t, err)
	//		test.Check(t, tt.Args, data[0])
	//	})
	//}
}
