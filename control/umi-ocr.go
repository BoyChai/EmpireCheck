package control

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type requestStruct struct {
	Base64  string `json:"base64"`
	Options struct {
		OcrLanguage     string `json:"ocr.language"`
		OcrCls          bool   `json:"ocr.cls"`
		OcrLimitSideLen int    `json:"ocr.limit_side_len"`
		TbpuParser      string `json:"tbpu.parser"`
		DataFormat      string `json:"data.format"`
	} `json:"options"`
}

type responseStruct struct {
	Code int `json:"code"`
	Data []struct {
		Box      [][]int `json:"box"`
		ClsLabel int     `json:"cls_label"`
		ClsScore float64 `json:"cls_score"`
		Score    float64 `json:"score"`
		Text     string  `json:"text"`
		End      string  `json:"end"`
	} `json:"data"`
	Score     float64 `json:"score"`
	Time      float64 `json:"time"`
	Timestamp float64 `json:"timestamp"`
}

type Lord struct {
	// 主公名字
	Name string `json:"name"`
	// 主公职业
	Career string `json:career`
	// 主公职位
	Position string `json:"position"`
	// 总繁荣
	// Prosperous int64
	Prosperous string `json:"prosperous"`
	// 周武勋
	// WeekMilitaryExploit int64
	WeekMilitaryExploit string `json:"week_military_exploit"`
	// 周贡献
	// WeekContribute int64
	WeekContribute string `json:"week_contribute"`
}

func UmiOcr(base64 string) ([]Lord, error) {
	// 创建请求数据
	request := requestStruct{
		Base64: base64,
	}
	request.Options.OcrLanguage = "models/config_chinese.txt"
	request.Options.OcrCls = true
	request.Options.OcrLimitSideLen = 999999
	request.Options.TbpuParser = "single_code"
	request.Options.DataFormat = "dict"

	// 将请求数据编码为JSON格式
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, errors.New("Error encoding JSON")
	}

	// 发送POST请求
	resp, err := http.Post(viper.GetString("ocr.URL")+"/api/ocr", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return nil, errors.New("Error sending POST request")
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, errors.New("Error reading response body")
	}

	// 解析响应数据
	var response responseStruct
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, errors.New("Error decoding JSON response")

	}
	var lordList []Lord
	for i, v := range response.Data[3 : len(response.Data)-1] {
		if i%2 == 0 {
			fields := strings.Fields(v.Text)
			lordList = append(lordList, Lord{
				Name:                fields[0],
				Career:              fields[1],
				Position:            fields[2],
				Prosperous:          fields[3],
				WeekMilitaryExploit: fields[4],
				WeekContribute:      fields[5],
			})
		}
	}
	// 输出
	jsonData, err = json.Marshal(lordList)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return nil, errors.New("Error converting to JSON")
	}
	fmt.Println(jsonData)

	// 输出 JSON 字符串
	return lordList, nil
}
