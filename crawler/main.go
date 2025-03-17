package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// 创建带缓存的Collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.CacheDir("./.cache"), // 启用缓存避免重复下载
		colly.Async(true),
	)

	// 配置限制规则
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*picsum.photos*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	// 创建图片保存目录
	_ = os.Mkdir("downloaded_images", 0755)

	// 图片处理回调
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		imgUrl, err := url.Parse(imgSrc)
		if err != nil {
			fmt.Printf("URL解析错误: %v\n", err)
			return
		}

		// 生成绝对URL
		absoluteURL := e.Request.URL.ResolveReference(imgUrl)
		fmt.Printf("[发现图片] %s\n", absoluteURL)

		// 异步下载图片
		_ = c.Visit(absoluteURL.String())
	})

	// 处理图片下载响应
	c.OnResponse(func(r *colly.Response) {
		// 过滤非图片内容
		contentType := http.DetectContentType(r.Body)
		if contentType[:5] != "image" {
			fmt.Printf("跳过非图片内容: %s\n", r.Request.URL)
			return
		}

		// 生成唯一文件名
		filename := filepath.Join("downloaded_images",
			fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(r.Request.URL.Path)))

		// 保存文件
		if err := os.WriteFile(filename, r.Body, 0644); err != nil {
			fmt.Printf("保存失败 %s: %v\n", filename, err)
			return
		}
		fmt.Printf("成功保存: %s (%s)\n", filename, humanSize(len(r.Body)))
	})

	// 错误处理
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("请求失败: %s → %v\n", r.Request.URL, err)
	})

	// 起始URL
	if err := c.Visit("https://picsum.photos/images"); err != nil {
		fmt.Printf("初始请求错误: %v\n", err)
	}

	c.Wait()
	fmt.Println("抓取完成！")
}

// 辅助函数：转换字节数为易读格式
func humanSize(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
