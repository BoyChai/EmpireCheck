package bot

import (
	"EmpireCheck/control"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/BoyChai/CoralBot/action"
	"github.com/BoyChai/CoralBot/bot"
	"github.com/BoyChai/CoralBot/run"
	"github.com/BoyChai/CoralBot/task"
	"github.com/spf13/viper"
)

func RunDingding() {
	Dhandle, _ := action.NewDingDingHandle(viper.GetString("dingding.AppKey"), viper.GetString("dingding.AppSecret"))
	dEvent := bot.DingDingEvent{
		Content: struct {
			RichText []struct {
				Text                string `json:"text"`
				PictureDownloadCode string `json:"pictureDownloadCode"`
				DownloadCode        string `json:"downloadCode"`
				Type                string `json:"type"`
			} `json:"richText"`
		}{
			RichText: make([]struct {
				Text                string `json:"text"`
				PictureDownloadCode string `json:"pictureDownloadCode"`
				DownloadCode        string `json:"downloadCode"`
				Type                string `json:"type"`
			}, 1), // 初始化大小为1
		},
	}
	task.NewTask(task.Task{
		Condition: []task.Condition{
			{Key: &dEvent.Msgtype,
				Value: "richText",
			},
		},
		Run: func() {
			str := *&dEvent.Content.RichText[0].Text
			// fmt.Println([]byte(str[len(str)-3 : len(str)]))
			// fmt.Println([]byte("/sm"))
			if str[len(str)-3:] == "/sm" {
				fmt.Printf("aaa")
				for i, v := range *&dEvent.Content.RichText {
					if i == 0 {
						continue
					}
					imgUrl, err := Dhandle.GetMessageFilesUrl(Dhandle.AppKey, v.DownloadCode)
					if err != nil {
						fmt.Println("1")
						msg := action.NewTextMsg("图片识别错误")
						Dhandle.SendGroupMessages(Dhandle.AppKey, msg, dEvent.ConversationId)
						return
					}
					str := imgUrlToBase64(imgUrl)
					lordList, err := control.UmiOcr(str)
					if err != nil {
						fmt.Println("2")
						msg := action.NewTextMsg("图片识别错误")
						Dhandle.SendGroupMessages(Dhandle.AppKey, msg, dEvent.ConversationId)
						return
					}
					control.ESControl(lordList)
				}
				msg := action.NewTextMsg("导入数据成功！！！")
				Dhandle.SendGroupMessages(Dhandle.AppKey, msg, dEvent.ConversationId)
			}
		},
	})

	run.Run(&dEvent, ":8080", false)

}

func imgUrlToBase64(url string) string {

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	return base64.StdEncoding.EncodeToString(data)
}
