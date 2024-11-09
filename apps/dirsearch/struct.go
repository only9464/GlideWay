package dirsearch

import "net/http"

type PathInfo struct {
	URL           string      `json:"url"`
	Path          string      `json:"path"`
	StatusCode    int         `json:"statusCode"`
	ContentType   string      `json:"contentType"`
	ContentLength int64       `json:"contentLength"`
	Header        http.Header `json:"header"`
}

type PathCallback func(PathInfo)
type ProgressCallback func(current, total int)
