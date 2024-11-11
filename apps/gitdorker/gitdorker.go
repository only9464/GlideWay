package gitdorker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GitdorkerCalculate(a, b int) int {
	fmt.Println("Gitdorker")
	return a + b
}

// 创建HTTP客户端
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SearchGithub 从GitHub API查询数据
func SearchGithub(ctx context.Context, mainKeyword string, subKeyword string, token string) *GithubResult {
	const apiUrl = "https://api.github.com/search/code"
	const maxPages = 10 // 最多获取10页数据，每页100条

	searchQuery := fmt.Sprintf("%s %s", mainKeyword, subKeyword)
	allItems := []string{}

	for page := 1; page <= maxPages; page++ {
		uri, _ := url.Parse(apiUrl)
		param := url.Values{}
		param.Set("q", searchQuery)
		param.Set("per_page", "100")
		param.Set("page", fmt.Sprintf("%d", page))
		uri.RawQuery = param.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			break
		}

		req.Header.Set("accept", "application/vnd.github.v3+json")
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
		req.Header.Set("User-Agent", "HelloGitHub")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error performing request:", err)
			break
		}

		// 检查速率限制
		if resp.StatusCode == 403 {
			fmt.Println("搜索完成  或  达到 GitHub API 速率限制")
			break
		}

		// 检查是否还有更多页
		if resp.StatusCode == 404 {
			break
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}

		var githubResp GithubResponse
		if err = json.Unmarshal(body, &githubResp); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			break
		}

		// 如果这页没有数据了，就退出
		if len(githubResp.Items) == 0 {
			break
		}

		// 添加这一页的结果
		allItems = append(allItems, extractUrls(githubResp.Items)...)

		// 输出进度
		fmt.Printf("已获取 %d 页数据，当前共 %d 条结果\n", page, len(allItems))

		// 添加延时以避免触发 GitHub API 限制
		time.Sleep(time.Second)
	}

	return &GithubResult{
		Status: true,
		Total:  float64(len(allItems)), // 更新为实际获取的数量
		Items:  allItems,
		Link:   fmt.Sprintf("https://github.com/search?q=%s&type=Code", url.QueryEscape(searchQuery)),
	}
}

// 辅助函数：提取URL列表
func extractUrls(items []struct {
	Html_url string `json:"html_url"`
}) []string {
	urls := make([]string, len(items))
	for i, item := range items {
		urls[i] = item.Html_url
	}
	return urls
}
