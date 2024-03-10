package es

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

var (
	client   *elastic.Client
	dataChan chan *EsData
)

type EsData struct {
	DataLog map[string]interface{}
	Index   string
}

// Init
func Init(address string, chanSize int) (err error) {

	dataChan = make(chan *EsData, chanSize)

	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		fmt.Println("create client failed, err:", err)
		return
	}
	fmt.Println("create es client success!")
	go sendToES()
	return
}

// SemdToEs 发送数据到Es
func sendToES() {
	for {
		select {
		case data := <-dataChan:
			put1, err := client.Index().Index(data.Index).Type("go").BodyJson(data.DataLog).Do(context.Background())
			if err != nil {
				fmt.Println("send to Es failed, err:", err)
			}
			fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// 开放结构体接口
func SendToChan(data *EsData) {
	dataChan <- data
}
