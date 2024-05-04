package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/longz666/news/lib"
	"log"
	"strings"
)

type Zhihu struct{}

// ZhihuResp 请求结果
type ZhihuResp struct {
	Data []struct {
		Content string `json:"content"`
	} `json:"data"`
}

func (c *Zhihu) Get(index int) (lib.Response, error) {
	url := fmt.Sprintf("https://v2.alapi.cn/api/zaobao?token=nIxqBcVR8gR1bnFt&format=txt", index)
	var resp lib.Response
	body, err := lib.Fetch(url)
	if err != nil {
		log.Printf("fetch zhihu url error:%s\n", err)
		return resp, err
	}
	var zhihuResp ZhihuResp
	if err = json.Unmarshal(body, &zhihuResp); err != nil {
		log.Printf("parse zhihuResp  error:%s\n", err)
		return resp, err
	}

	if len(zhihuResp.Data) > 0 {
		buffer := bytes.NewBuffer([]byte(zhihuResp.Data[0].Content))
		doc, err := goquery.NewDocumentFromReader(buffer)
		if err != nil {
			log.Printf("zhihu go-goquery document error:%s\n", err)
			return resp, err
		}
		doc.Find("p").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			resp.AllData = append(resp.AllData, text)
			if strings.Contains(text, "、") {
				resp.Data.News = append(resp.Data.News, strings.Join(strings.Split(text, "、")[1:], "、"))
			}
		})

		if len(resp.AllData) > 0 {
			resp.Data.Title = resp.AllData[0]
			resp.Data.Date = resp.AllData[1]
			resp.Data.Weiyu = resp.AllData[len(resp.AllData)-1]
			resp.AllData[0], resp.AllData[1] = resp.AllData[1], resp.AllData[0]
		}
	}
	return resp, nil
}
