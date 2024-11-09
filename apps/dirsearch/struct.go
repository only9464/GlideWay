package dirsearch

import (
	"context"
	"net/http"
)

type PathInfo struct {
	URL           string      `json:"url"`
	Path          string      `json:"path"`
	StatusCode    int         `json:"statusCode"`
	ContentType   string      `json:"contentType"`
	ContentLength int64       `json:"contentLength"`
	Header        http.Header `json:"header"`
}

// 目录扫描相关结构体和变量
type DirsearchProgress struct {
	Current int     `json:"current"`
	Total   int     `json:"total"`
	Speed   float64 `json:"speed"` // 添加速度字段
}

type PathResult struct {
	Path          string `json:"path"`
	FullUrl       string `json:"fullUrl"`
	StatusCode    int    `json:"statusCode"`
	ContentType   string `json:"contentType"`
	ContentLength int64  `json:"contentLength"`
}

type DirsearchControl struct {
	cancel     context.CancelFunc
	scanned    int32
	totalPaths int32
}

type PathCallback func(PathInfo)
type ProgressCallback func(current, total int)
