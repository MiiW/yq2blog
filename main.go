package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wujiyu115/yuqueg"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"yq2blog/blog"
)

var (
	ActionTypePublish = "publish"
	ActionTypeUpdate  = "update"
	ActionTypeDelete  = "delete"
)

var (
	repoNamespace string
	yq            *yuqueg.Service
	actor         blog.BlogActor
	slugTitleMap  map[string]string = make(map[string]string)
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if e := viper.ReadInConfig(); e != nil {
		fmt.Println("read in config error:", e.Error())
		os.Exit(-1)
	}
	repoNamespace = viper.GetString("yuque.repo")
	switch viper.GetString("blog.type") {
	case "hugo":
		actor = &blog.HugoActor{
			Deloyment: viper.GetString("blog.deployment"),
			Path:      viper.GetString("blog.path"),
		}
	default:
		fmt.Println("unsupported blog type:", viper.GetString("yuque.repo"))
		os.Exit(-1)
	}

	yq = yuqueg.NewService(viper.GetString("yuque.token"))
	list, _ := yq.Doc.List(repoNamespace)
	for _, v := range list.Data {
		if v.Status == 1 {
			slugTitleMap[v.Slug] = v.Title
		}
	}
}

func main() {
	r := gin.Default()

	r.Handle(http.MethodGet, "/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "_pong",
		})
	})

	r.Handle(http.MethodPost, "/webhook", yqWebhook)

	r.Run(":9091")
}

func yqWebhook(c *gin.Context) {

	//读取并解析webhook
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errMsg := fmt.Sprintf("read body fail: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, errMsg)
		return
	}
	action := HookAction{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		errMsg := fmt.Sprintf("body unmarshal fail: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	//判定仓库是否匹配
	if !strings.HasPrefix(action.Data.Path, repoNamespace) {
		errMsg := fmt.Sprintf("config repo and webhook not match: %s -> %s\n", repoNamespace, action.Data.Path)
		c.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	//根据仓库和slug获取文章内容
	if action.Data.ActionType == ActionTypeDelete {
		if ok, out := actor.DeleteBlog(action.Data.Title); !ok {
			fmt.Println(out)
			return
		}
		for k, v := range slugTitleMap {
			if v == action.Data.Title {
				delete(slugTitleMap, k)
			}
		}
		fmt.Println("doc delete success")
		return
	}

	doc, err := yq.Doc.Get(repoNamespace, action.Data.Slug, &yuqueg.DocGet{Raw: 1})
	if err != nil {
		errMsg := fmt.Sprintf("get doc fail: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	//判定hook类型
	switch action.Data.ActionType {
	case ActionTypePublish:
		if ok, out := actor.PublishBlog(action.Data.FirstPublishedTime, action.Data.Title, doc.Data.Body); !ok {
			fmt.Println(out)
			return
		}
		slugTitleMap[action.Data.Slug] = action.Data.Title
		fmt.Println("doc add success")
	case ActionTypeUpdate:
		if ok, out := actor.UpdateBlog(action.Data.FirstPublishedTime, slugTitleMap[action.Data.Slug], action.Data.Title, doc.Data.Body); !ok {
			fmt.Println(out)
			return
		}
		slugTitleMap[action.Data.Slug] = action.Data.Title
		fmt.Println("doc update success")
	}

}

type HookAction struct {
	Data struct {
		ID                 int    `json:"id"`
		Slug               string `json:"slug"`
		Title              string `json:"title"`
		ActionType         string `json:"action_type"`
		Publish            bool   `json:"publish"`
		Path               string `json:"path"`
		FirstPublishedTime string `json:"first_published_at"`
	} `json:"data"`
}
