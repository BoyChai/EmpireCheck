package control

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

type LordModel struct {
	CreatedAt string `json:"created_at"`
	Lord
}

func (LordModel) Index() string {
	return "lord-index-" + time.Now().Format("2006-01-02")
}

func (LordModel) Mapping() string {
	return `
  {
    "mappings": {
      "properties": {
        "created_at": {
          "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time"
        },
          "name": {
          "type": "text"
   		},
          "career": {
          "type": "text"
    	},
        "position": {
          "type": "text"
         },
        "prosperous": {
           "type": "text"
         },
        "week_military_exploit": {
          "type": "text"
        },
        "week_contribute": {
        "type": "text"
      }
      }
    }
  }
  `
}

var ESClient *elastic.Client

func InitEsClient() {
	var err error
	ESClient, err = elastic.NewClient(
		elastic.SetURL(viper.GetString("es.URL")),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("", ""),
	)
	if err != nil {
		log.Fatalln("es连接错误:", err)
	}
}
func ESControl(data []Lord) {
	// 检查索引是否存在
	if exists, _ := ESClient.IndexExists(LordModel{}.Index()).Do(context.Background()); !exists {
		// 创建索引
		_, err := ESClient.CreateIndex(LordModel{}.Index()).BodyString(LordModel{}.Mapping()).Do(context.Background())
		if err != nil {
			log.Fatalln("创建索引失败:", err)
		}
	}

	for _, v := range data {
		// 添加文档
		indexResponse, err := ESClient.Index().Index(LordModel{}.Index()).BodyJson(LordModel{
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			Lord:      v,
		}).Do(context.Background())

		if err != nil {
			log.Fatalln("添加文档失败:", err)
		}
		fmt.Printf("%#v\n", indexResponse)
	}

}
