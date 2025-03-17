package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings" // 补充缺失的strings包
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// 创建带超时设置的Collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.Async(true), // 必须启用异步模式
	)

	// 配置超时客户端（关键修复点1）
	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 15 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
	})

	// 创建图片目录
	_ = os.MkdirAll("downloaded_images", 0755)

	// 图片处理逻辑（关键修复点2：保持异步访问）
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		absoluteURL := e.Request.AbsoluteURL(imgSrc) // 使用内置方法解析URL

		// 过滤非图片内容（关键修复点3：使用MIME类型判断更可靠）
		if !strings.HasPrefix(http.DetectContentType(e.Response.Body), "image/") {
			return
		}

		// 异步下载图片
		_ = e.Request.Visit(absoluteURL)
	})

	// 图片下载处理器
	c.OnResponse(func(r *colly.Response) {
		// 类型检测
		if !strings.HasPrefix(r.Headers.Get("Content-Type"), "image/") {
			return
		}

		// 生成唯一文件名
		filename := filepath.Join("downloaded_images",
			fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(r.Request.URL.Path)))

		// 保存文件
		if err := os.WriteFile(filename, r.Body, 0644); err != nil {
			fmt.Printf("保存失败: %v\n", err)
			return
		}
		fmt.Printf("成功下载: %s\n", filename)
	})

	// 错误处理
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("请求失败: %s → %v\n", r.Request.URL, err)
	})

	// 启动爬取
	if err := c.Visit("https://picsum.photos/images"); err != nil {
		fmt.Printf("初始请求错误: %v\n", err)
	}

	c.Wait() // 必须等待异步任务
}
