package blog

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type HugoActor struct {
	Deloyment string
	Path      string
}

func (h *HugoActor) ActorTest() (bool, string) {
	if out, err := exec.Command("hugo", "version").Output(); err != nil {
		fmt.Println(out)
		return false, err.Error()
	}
	//maybe more test
	return true, ""
}

func (h *HugoActor) PublishBlog(publishTime, title, content string) (bool, string) {
	doc := generateDoc(publishTime, title, content)
	if err := ioutil.WriteFile(h.Path+"\\content\\"+title+".md",
		[]byte(doc), 0755); err != nil {
		return false, fmt.Sprintf("create md file error while publishing: %s", err.Error())
	}

	//cmd := exec.Command("hugo", "server --theme=hyde --buildDrafts")
	//if out,err := cmd.Output(); err != nil {
	//	fmt.Println(string(out))
	//	return false, err.Error()
	//}
	//maybe more test
	return true, ""
}

func (h *HugoActor) UpdateBlog(firstPublishTime, oldTitle, newTitle, content string) (bool, string) {
	if err := os.Remove(h.Path + "\\content\\" + oldTitle + ".md"); err != nil {
		return false, fmt.Sprintf("delete file error while updating: %s", err.Error())
	}

	if err := ioutil.WriteFile(h.Path+"\\content\\"+newTitle+".md",
		[]byte(generateDoc(firstPublishTime, newTitle, content)), 0755); err != nil {
		return false, fmt.Sprintf("create md file error while updating: %s", err.Error())
	}

	return true, ""
}

func (h *HugoActor) DeleteBlog(title string) (bool, string) {
	if err := os.Remove(h.Path + "\\content\\" + title + ".md"); err != nil {
		return false, fmt.Sprintf("delete file error while updating: %s", err.Error())
	}
	return true, ""
}

func generateDoc(date, title, content string) string {
	const gendocFrontmatterTemplate = `---
date: %s
title: "%s"
---
`
	return fmt.Sprintf(gendocFrontmatterTemplate, date, title) + content
}
