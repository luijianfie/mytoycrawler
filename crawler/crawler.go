package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func CrawlPicFromUrl(url, selector, path string, sel interface{}, sleep int) {

	html, err := GetHttpHtmlContent(url, selector, sel, sleep)
	if err != nil {
		fmt.Println(err)
	}
	GetSpecialData(html, "img[src]", "src", "./remote")

}

//get html content from target url
func GetHttpHtmlContent(url string, selector string, sel interface{}, sleep int) (string, error) {

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 20*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Duration(sleep)*time.Second),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		log.Printf("Run err : %v\n", err)
		return "", err
	}

	return htmlContent, nil
}

func GetSpecialData(htmlContent, selector, attr, path string) (string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Errorf("%v\n", err)
		return "", err
	}

	var str string
	dom.Find(selector).Each(func(i int, selection *goquery.Selection) {

		url, exist := selection.Attr(attr)

		if exist {
			savecontent(path, url)
		}
	})
	return str, nil
}

//dir 暂时不设置
func savecontent(dir string, url string) {

	slices := strings.Split(url, "//")
	if len(slices) == 0 {
		fmt.Printf("can't get pic name from url:%s.\n", url)
		return
	}
	urlSub := slices[len(slices)-1]
	resp, err := http.Get("https://" + urlSub)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)

	names := strings.Split(urlSub, "/")
	if len(names) == 0 {
		fmt.Printf("invalid file name %s.\n", urlSub)
	}
	name := names[len(names)-1]

	path := filepath.Join(dir, name)

	out, _ := os.Create(path)
	defer out.Close()
	_, ioerr := io.Copy(out, bytes.NewReader(body))
	if ioerr != nil {
		fmt.Printf("failed to write file %s,url:%s.\n", name, url)
	}

	fmt.Printf("save file %s to disk.\n", name)
	return
}

// 截图，以图片的格式，或者以pdf的方式
func sample1(url1 string, url2 string) {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// run
	var b1, b2 []byte
	if err := chromedp.Run(ctx,
		// emulate iPhone 7 landscape
		chromedp.Emulate(device.NokiaN9landscape),
		chromedp.Navigate(url1),
		// chromedp.CaptureScreenshot(&b1),
		chromedp.FullScreenshot(&b1, 100),
		// reset
		chromedp.Emulate(device.Reset),

		// set really large viewport
		chromedp.EmulateViewport(1920, 2000),
		chromedp.Navigate(url2),
		chromedp.CaptureScreenshot(&b2),
		chromedp.FullScreenshot(&b2, 100),
	); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("baidu_IPhone8Plus.png", b1, 0777); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("baidu_PC.png", b2, 0777); err != nil {
		log.Fatal(err)
	}

	var b3 []byte
	if err := chromedp.Run(ctx, printToPDF(`https://www.ctrip.com/?sid=155952&allianceid=4897&ouid=index`, &b3)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("./remote/new.pdf", b3, 0777); err != nil {
		log.Fatal(err)
	}
}

//生成打印pdf的task
func printToPDF(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return nil
			}
			*res = buf
			return nil
		}),
	}
}
