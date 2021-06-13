package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/wujiyu115/yuqueg"
)

type FrontMatter struct {
	Title         string      `json:"title"`
	Slug          string      `json:"slug"`
	Description   interface{} `json:"description"`
	Date          time.Time   `json:"date"`
	Tags          []string    `json:"tags"`
	IsCJKLanguage bool
}

type config struct {
	NameSpace   []string `toml:"namespace"`
	ImageMirror string   `toml:"image_mirror"`
	Dir         string   `toml:"dir"`
}

type Yuque2mdConfig struct {
	YuQue config `toml:"yuque2md"`
}

func sinkToFile(namespace string, slug string, doc yuqueg.DocDetail, config Yuque2mdConfig) error {
	filePath := config.YuQue.Dir + "/" + namespace + "/" + slug + ".md"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	body := doc.Data.Body

	// body = strings.ReplaceAll(body, "<br />", "\n")
	// re3, _ := regexp.Compile(`<a.+><\/a>`)
	// body = re3.ReplaceAllString(body, "")

	frontMatter := FrontMatter{
		Title:         doc.Data.Title,
		Slug:          doc.Data.Slug,
		Description:   doc.Data.CustomDescription,
		Date:          doc.Data.FirstPublishedAt,
		Tags:          []string{doc.Data.Book.Name},
		IsCJKLanguage: true,
	}
	jsonData, err := json.Marshal(frontMatter)
	if err != nil {
		fmt.Println("json转换失败", err)
		return err
	}
	body = string(jsonData) + "\n\n" + body
	body = strings.ReplaceAll(body, "https://cdn.nlark.com", config.YuQue.ImageMirror)

	write.WriteString(body)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return nil
}

func downloadNameSpace(namespace string, token string, config Yuque2mdConfig) error {
	l := yuqueg.L
	yu := yuqueg.NewService(token)

	d, err := yu.Repo.GetToc(namespace)

	if err != nil {
		l.Info(err)
		return err
	}

	os.MkdirAll(config.YuQue.Dir+"/"+namespace, os.ModePerm)
	for i := 0; i < len(d.Data); i++ {
		current := d.Data[i]
		doc, err := yu.Doc.Get(namespace, current.Slug, &yuqueg.DocGet{Raw: 1})
		if err != nil {
			panic(err)
		}
		sinkToFile(namespace, current.Slug, doc, config)
		l.Info("空间：", namespace, " 文章：", current.Title, " 已完成")
	}
	return nil
}

func main() {
	var (
		yq Yuque2mdConfig
	)
	token := os.Getenv("YUQUE_TOKEN")
	_, err := toml.DecodeFile("config.toml", &yq)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(yq.YuQue.NameSpace); i++ {
		downloadNameSpace(yq.YuQue.NameSpace[i], token, yq)
	}
}
