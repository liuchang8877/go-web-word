package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
)

func main() {
	r := gin.Default()

	// 静态文件服务，用于提供下载链接
	r.Static("/documents", "./documents")

	r.POST("/generate", func(c *gin.Context) {
		var req struct {
			Title   string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查模板文件是否存在
		templatePath := "./template.docx"
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Template file does not exist"})
			log.Printf("Error: Template file %s does not exist", templatePath)
			return
		}

		// 读取模板文件
		r, err := docx.ReadDocxFile(templatePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read template file"})
			log.Printf("Error reading template file %s: %v", templatePath, err)
			return
		}
		defer r.Close()

		doc := r.Editable()
		// 替换占位符
		doc.Replace("{{title}}", req.Title, -1)
		doc.Replace("{{content}}", req.Content, -1)

		// 确保documents目录存在
		outputDir := "./documents"
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			os.Mkdir(outputDir, 0755)
		}

		// 保存到文件
		filename := fmt.Sprintf("document-%d.docx", time.Now().Unix())
		filepath := filepath.Join(outputDir, filename)
		err = doc.WriteToFile(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to file"})
			log.Printf("Error writing to file %s: %v", filepath, err)
			return
		}

		// 返回下载链接
		downloadURL := fmt.Sprintf("%s/%s", "/documents", filename)
		c.JSON(http.StatusOK, gin.H{
			"download_url": downloadURL,
		})
	})

	if err := r.Run(":3000"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
