package anna

import (
	"fmt"
	"net/url"
	"sync"

	"strings"

	"github.com/PuerkitoBio/goquery"
	colly "github.com/gocolly/colly/v2"
)

const (
	AnnasSearchEndpoint   = "https://annas-archive.org/search?q=%s"
	AnnasDownloadEndpoint = "https://annas-archive.org/dyn/api/fast_download.json?md5=%s&key=%s"
)

func FindBook(query string) ([]*Book, error) {
	c := colly.NewCollector(
		colly.Async(true),
	)

	bookList := make([]*colly.HTMLElement, 0)
	var mu sync.Mutex

	c.OnHTML("a[href*='/md5/'].js-vim-focus", func(e *colly.HTMLElement) {
		mu.Lock()
		bookList = append(bookList, e)
		mu.Unlock()
	})

	c.OnHTML("a[href*='/md5/']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		alreadyFound := false
		mu.Lock()
		for _, existing := range bookList {
			if existing.Attr("href") == href {
				alreadyFound = true
				break
			}
		}
		if !alreadyFound {
			bookList = append(bookList, e)
		}
		mu.Unlock()
	})

	c.OnRequest(func(r *colly.Request) {
	})

	fullURL := fmt.Sprintf(AnnasSearchEndpoint, url.QueryEscape(query))
	c.Visit(fullURL)
	c.Wait()

	bookListParsed := make([]*Book, 0, len(bookList))

	for _, e := range bookList {
		container := e.DOM.Closest("div.flex")
		if container.Length() == 0 {
			continue
		}

		title := strings.TrimSpace(e.Text)

		var authors string
		titleLink := e.DOM
		authorLink := titleLink.NextFiltered("a")
		if authorLink.Length() > 0 {
			authors = strings.TrimSpace(authorLink.Text())
			authors = strings.TrimSpace(strings.Replace(authors, "👤", "", -1))
			authors = strings.TrimSpace(authors)
		}

		var meta string

		container.Find("div").Each(func(j int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())

			if strings.Contains(text, "·") &&
				(strings.Contains(text, "MB") || strings.Contains(text, "KB") || strings.Contains(text, "GB")) &&
				(strings.Contains(text, "[") || strings.Contains(text, "ZIP") || strings.Contains(text, "PDF") || strings.Contains(text, "EPUB")) {

				cleanText := text

				if saveIndex := strings.Index(cleanText, "Save"); saveIndex != -1 {
					cleanText = strings.TrimSpace(cleanText[:saveIndex])
				}

				if funcIndex := strings.Index(cleanText, "(function"); funcIndex != -1 {
					cleanText = strings.TrimSpace(cleanText[:funcIndex])
				}

				if len(cleanText) > 10 && strings.Contains(cleanText, "·") {
					meta = cleanText
				}
			}
		})

		language, format, size := extractMetaInformation(meta)

		trimmedFormat := strings.TrimSpace(format)
		if trimmedFormat != "" {
			trimmedFormat = strings.Trim(trimmedFormat, "[]()·")
		}

		var publisher string

		/*var filePath string
		container.Find("div.text-\\[9px\\].text-gray-500").Each(func(j int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if strings.Contains(text, "/") && (strings.Contains(text, ".pdf") || strings.Contains(text, ".epub") || strings.Contains(text, ".mobi")) {
				filePath = text
			}
		})*/

		var coverURL, coverData string

		coverLink := container.Find("a.custom-a.block").First()
		if coverLink.Length() > 0 {
			coverImg := coverLink.Find("img")
			if coverImg.Length() > 0 {
				coverURL = coverImg.AttrOr("src", "")
				if coverURL != "" {
					coverURL = e.Request.AbsoluteURL(coverURL)
				}
			}

			if coverURL == "" {
				fallbackDiv := coverLink.Find("div.js-aarecord-list-fallback-cover")
				if fallbackDiv.Length() > 0 {
					bgColor, _ := fallbackDiv.Attr("style")
					titleDiv := fallbackDiv.Find("div.font-bold.text-violet-900")
					authorDiv := fallbackDiv.Find("div.font-bold.text-amber-900")

					coverData = fmt.Sprintf("fallback:bg=%s;title=%s;author=%s",
						bgColor,
						titleDiv.Text(),
						authorDiv.Text())
				}
			}
		}

		if coverURL == "" {
			allImgs := container.Find("img")
			allImgs.Each(func(idx int, img *goquery.Selection) {
				src := img.AttrOr("src", "")
				if coverURL == "" && src != "" {
					coverURL = e.Request.AbsoluteURL(src)
				}
			})
		}

		if coverURL == "" {
			allLinks := container.Find("a")
			allLinks.Each(func(idx int, link *goquery.Selection) {
				linkHTML, _ := link.Html()
				if strings.Contains(linkHTML, "<img") {
					img := link.Find("img").First()
					if img.Length() > 0 {
						src := img.AttrOr("src", "")
						if src != "" && coverURL == "" {
							coverURL = e.Request.AbsoluteURL(src)
						}
					}
				}
			})
		}

		if coverURL == "" {
			parent := container.Parent()
			if parent.Length() > 0 {
				parentImgs := parent.Find("img")
				parentImgs.Each(func(idx int, img *goquery.Selection) {
					src := img.AttrOr("src", "")
					if coverURL == "" && src != "" {
						coverURL = e.Request.AbsoluteURL(src)
					}
				})
			}
		}

		link := e.Attr("href")
		hash := strings.TrimPrefix(link, "/md5/")

		book := &Book{
			Language:  strings.TrimSpace(language),
			Format:    trimmedFormat,
			Size:      strings.TrimSpace(size),
			Title:     strings.TrimSpace(title),
			Publisher: strings.TrimSpace(publisher),
			Authors:   strings.TrimSpace(authors),
			URL:       e.Request.AbsoluteURL(link),
			Hash:      hash,
			CoverURL:  coverURL,
			CoverData: coverData,
		}

		bookListParsed = append(bookListParsed, book)
	}

	return bookListParsed, nil
}
