package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type FrontMatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Description string   `yaml:"description"`
	Slug        string   `yaml:"slug"`
	Image       string   `yaml:"image,omitempty"`
	Categories  []string `yaml:"categories,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

func main() {
	// 解析命令行参数
	filePath := flag.String("file", "", "The path of the markdown file to process")
	flag.Parse()

	// 检查文件路径是否为空
	if *filePath == "" {
		fmt.Println("Error: file path is required")
		os.Exit(1)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// 解析文件内容
	lines := strings.Split(string(content), "\n")
	title := strings.TrimPrefix(filepath.Base(*filePath), "index")
	date := time.Now().Format("2006-01-02 15:04:05-0700")
	description := ""
	image := ""
	categories := []string{}
	tags := []string{}
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			lines = append([]string{line}, lines...)
			break
		}
		if strings.HasPrefix(line, "description: ") {
			description = strings.TrimSpace(strings.TrimPrefix(line, "description: "))
			lines[i] = ""
		}
		if strings.HasPrefix(line, "image: ") {
			image = strings.TrimSpace(strings.TrimPrefix(line, "image: "))
			lines[i] = ""
		}
		if strings.HasPrefix(line, "categories: ") {
			categories = strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "categories: ")), ",")
			for i, cat := range categories {
				categories[i] = strings.TrimSpace(cat)
			}
			lines[i] = ""
		}
		if strings.HasPrefix(line, "tags: ") {
			tags = strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "tags: ")), ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
			lines[i] = ""
		}
	}
	slugStr := slug.Make(title)

	// 生成新的 YAML 头信息
	frontMatter := FrontMatter{
		Title:       title,
		Date:        date,
		Description: description,
		Slug:        slugStr,
		Image:       image,
		Categories:  categories,
		Tags:        tags,
	}
	fmStr := fmt.Sprintf("---\n%v\n---\n", frontMatter)

	// 将新的 YAML 头信息和原始内容写入文件
	newContent := fmStr + strings.Join(lines, "\n")
	err = ioutil.WriteFile(*filePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
}
