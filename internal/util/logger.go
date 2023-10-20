package util

import (
	"github.com/highlight/highlight/sdk/highlight-go"
	hlog "github.com/highlight/highlight/sdk/highlight-go/log"
)

func SetupLogger() func() {
	// setup the highlight SDK
	highlight.SetProjectID(GetEnvWithDefault("HIGHLIGHT_ID", "dev"))
	highlight.Start(
		highlight.WithServiceName("top-sites-ranking-api"),
		highlight.WithServiceVersion(GetEnvWithDefault("VERSION", "dev")),
	)

	// setup highlight logrus hook
	hlog.Init()

	return highlight.Stop
}
