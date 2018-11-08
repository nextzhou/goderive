package utils

import (
	"fmt"
	"runtime"
	"strings"
)

type DeriveComment struct {
	Plugin     string
	OptionsStr string
}

func MatchDeriveComment(cmt string) (*DeriveComment, error) {
	cmt = strings.TrimPrefix(cmt, "//")
	cmt = strings.TrimSpace(cmt)
	if !strings.HasPrefix(cmt, "derive-") {
		return nil, nil
	}
	cmt = strings.TrimPrefix(cmt, "derive-")
	splitIdx := strings.Index(cmt, ":")
	dc := new(DeriveComment)
	if splitIdx == -1 {
		dc.Plugin = cmt
	} else {
		dc.Plugin = cmt[:splitIdx]
		dc.OptionsStr = strings.TrimSpace(cmt[splitIdx+1:])
	}
	if !ValidateIdentName(dc.Plugin) {
		return nil, fmt.Errorf("invalid plugin name %#v", dc.Plugin)
	}
	return dc, nil
}

func MatchPluginComment(cmt string) (*DeriveComment, error) {
	// get plugin name by caller file path
	const pluginDir = "/goderive/plugin/"
	_, file, _, _ := runtime.Caller(1)
	pos := strings.LastIndex(file, pluginDir) + len(pluginDir)
	offset := strings.IndexByte(file[pos:], '/')
	plugin := file[pos : pos+offset]

	dc, err := MatchDeriveComment(cmt)
	if err != nil {
		return nil, err
	}
	if dc != nil && dc.Plugin != plugin {
		return nil, &UnmatchedError{Ident: "plugin", Got: dc.Plugin, Expected: plugin}
	}
	return dc, nil
}
