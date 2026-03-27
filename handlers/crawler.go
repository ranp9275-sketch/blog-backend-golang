package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/google/uuid"
)

// FetchArticleRequest 请求参数
type FetchArticleRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// FetchArticleResponse 响应数据
type FetchArticleResponse struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CoverImage string `json:"cover_image"`
	CategoryID string `json:"category_id"`
}

// FetchArticleByURL 通过URL抓取外部文章并转换为Markdown
func (h *Handlers) FetchArticleByURL(c *gin.Context) {
	var req FetchArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL provided"})
		return
	}

	// 1. 设置超时，防止请求挂起
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 2. 发起 HTTP GET 请求
	httpReq, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 添加基础请求头，模拟浏览器防止被拦截
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch URL: %v", err)})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to fetch URL, status code: %d", resp.StatusCode)})
		return
	}

	// 3. 使用 goquery 提前处理 HTML 中的图片（解决懒加载和相对路径）
	parsedUrl, err := resp.Request.URL.Parse("")
	if err != nil {
		parsedUrl = resp.Request.URL
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse HTML document: %v", err)})
		return
	}

	var fallbackCover string
	
	type imgProcess struct {
		selection *goquery.Selection
		src       string
	}
	var imgsToDownload []imgProcess

	lazyAttrs := []string{"data-src", "data-original", "data-lazy-src", "data-actualsrc", "origin-src", "data-src-retina", "v-lazy", "src-large"}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		for _, attr := range lazyAttrs {
			if val, exists := s.Attr(attr); exists && val != "" {
				s.SetAttr("src", val)
				break 
			}
		}

		if src, exists := s.Attr("src"); exists && strings.TrimSpace(src) != "" {
			// 使用 Go 标准库的 ResolveReference 完美处理一切相对路径（包括无前缀的目录名、./、../、// 等）
			imgUrl, err := url.Parse(src)
			if err == nil {
				absUrl := parsedUrl.ResolveReference(imgUrl)
				src = absUrl.String()
				s.SetAttr("src", src)
			}

			if strings.HasPrefix(src, "http") {
				imgsToDownload = append(imgsToDownload, imgProcess{selection: s, src: src})
			}
		}
	})

	// 并发下载图片替换资源防盗链
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, item := range imgsToDownload {
		wg.Add(1)
		go func(img imgProcess) {
			defer wg.Done()
			localPath, err := downloadImage(client, img.src, req.URL)
			if err == nil {
				mu.Lock()
				img.selection.SetAttr("src", localPath)
				// 获取第一张常规成功下载图片作为封面备选
				if fallbackCover == "" && !strings.Contains(strings.ToLower(img.src), "icon") && !strings.Contains(strings.ToLower(img.src), "avatar") && !strings.Contains(strings.ToLower(img.src), "logo") {
					fallbackCover = localPath
				}
				mu.Unlock()
			}
		}(item)
	}
	wg.Wait()

	htmlStr, err := doc.Html()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate optimized HTML"})
		return
	}

	// 提取正文内容和标题
	article, err := readability.FromReader(strings.NewReader(htmlStr), parsedUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse article content: %v", err)})
		return
	}

	// 如果没有提取到任何内容，报错
	if article.TextContent == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "No readable content found on this page"})
		return
	}

	// 4. 将 HTML 转换为 Markdown
	converter := md.NewConverter("", true, nil)
	
	// 为了克隆的博客文章更美观，可以在此添加一些原链接的引用
	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		// 如果转换失败，降级使用纯文本提取
		markdown = article.TextContent
	}

	// 追加原文来源信息，尊重版权
	sourceNotice := fmt.Sprintf("\n\n---\n> **说明**：本文通过 URL 克隆抓取。\n> **原文链接**：[%s](%s)\n> **抓取时间**：%s", article.Title, req.URL, time.Now().Format("2006-01-02 15:04:05"))
	markdownContent := markdown + sourceNotice

	// 5. 智能匹配分类
	matchedCategoryID := ""
	categories, err := h.repo.GetCategories()
	if err == nil && len(categories) > 0 {
		// 将文章标题和部分正文转为小写用于匹配
		searchTarget := strings.ToLower(article.Title + " " + article.Excerpt)
		
		for _, cat := range categories {
			catName := strings.ToLower(cat.Name)
			// 如果正文或标题包含分类名称（比如 "Go", "Vue", "前端"），则命中
			if strings.Contains(searchTarget, catName) {
				matchedCategoryID = cat.ID
				break // 找到一个匹配的就退出
			}
		}
	}

	// 6. 提取封面图片
	coverImage := article.Image
	if coverImage == "" {
		coverImage = fallbackCover
	}
	
	// 如果提取出来的封面是相对路径或双斜杠的缺省协议，做一层兜底修复
	if coverImage != "" {
		if strings.HasPrefix(coverImage, "//") {
			coverImage = "https:" + coverImage
		} else if strings.HasPrefix(coverImage, "/") {
			coverImage = fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, coverImage)
		}
	}

	// 7. 返回给前端
	c.JSON(http.StatusOK, FetchArticleResponse{
		Title:      article.Title,
		Content:    markdownContent,
		CoverImage: coverImage,
		CategoryID: matchedCategoryID,
	})
}

// downloadImage 抓取远程图片到本地 uploads 目录，并返回相对路径URL
func downloadImage(client *http.Client, urlStr string, referer string) (string, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	// 模拟正常浏览器请求和原始 Referer 绕过防盗链
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", referer)
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/*,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	ext := ".jpg"
	if strings.Contains(contentType, "png") {
		ext = ".png"
	} else if strings.Contains(contentType, "gif") {
		ext = ".gif"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	} else if strings.Contains(contentType, "svg") {
		ext = ".svg"
	}

	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, 0755)

	filename := fmt.Sprintf("fetch_%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
	dst := filepath.Join(uploadDir, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}
