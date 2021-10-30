package main

func main() {

    url := `http://ly88.yntxlx.com/tuxiang003/ynpc/?source=baidu&plan=%E4%BA%91%E5%8D%97PC&unit=%E4%BA%91%E5%8D%972&keyword=%E5%8E%BB%E4%BA%91%E5%8D%97%E7%91%9E%E4%B8%BD%E6%97%85%E6%B8%B8%E6%94%BB%E7%95%A5&e_creative=53855897706&e_keywordid=332031491385&e_keywordid2=282607238015&bd_vid=10371793418743331979&hdfshare=18988154971`
    selector := ``
    sel := `document.querySelector("body")`

    crwaler.CrawlPicFromUrl(url, selector, "./remote", sel, 3)
}
