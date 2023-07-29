package main

import (
	"idcoding/utils"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	router.Static("/static", "static")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		message := c.DefaultQuery("message", "")
		c.HTML(http.StatusOK, "index.html", gin.H{"message": message})
	})
	router.GET("/lesson/4", func(c *gin.Context) {
		message := c.DefaultQuery("message", "")
		c.HTML(http.StatusOK, "4.html", gin.H{"message": message})
	})
	router.POST("/set/4", handleRepoURL)
	router.POST("/upload", handleFileUpload)

	return router
}
func handleRepoURL(c *gin.Context) {
	// Получаем имя из формы
	name := c.PostForm("name")

	// Проверка наличия имени
	if name == "" {
		c.Redirect(http.StatusFound, "/?message=Имя не заполнено")
		log.Println("Нет имени")
		return
	}

	repo := c.PostForm("repo")

	if repo == "" {
		c.Redirect(http.StatusFound, "/?message=Ссылка на репо не указана")
		log.Println("Нет имени")
		return
	}
	curr := utils.HomeworkMap[4]

	if curr.HomeWorks == nil {
		curr.HomeWorks = make(map[string]string)
	}

	curr.HomeWorks[name] = repo

	utils.HomeworkMap[4] = curr

	utils.UpdateCache()
	// Все прошло успешно
	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно загружен.", "name": name})
}

func handleFileUpload(c *gin.Context) {
	// Максимальный размер загружаемого файла (измените по своему усмотрению)
	maxSize := int64(10 << 20) // 10 MB

	// Получаем файл из формы
	file, err := c.FormFile("file")
	if err != nil {
		c.Redirect(http.StatusFound, "/?message=Файл не выбран")
		log.Println(err)
		return
	}

	// Проверка размера файла
	if file.Size > maxSize {
		c.Redirect(http.StatusFound, "/?message=Размер файла превышает лимит")
		log.Println(file.Size)

		return
	}

	if !isArchive(file) {
		c.Redirect(http.StatusFound, "/?message=Файл не является архивом")
		log.Println("Не архив")
		return
	}

	// Получаем имя из формы
	name := c.PostForm("name")

	// Проверка наличия имени
	if name == "" {
		c.Redirect(http.StatusFound, "/?message=Имя не заполнено")
		log.Println("Нет имени")
		return
	}
	lesson := c.PostForm("lesson")

	if lesson == "" {
		c.Redirect(http.StatusFound, "/?message=Урок не выбран")
		log.Println("Нет имени")
		return
	}

	// Получаем временное имя файла
	tempFile := filepath.Join("./files/lesson"+lesson+"/"+name, file.Filename)

	// Сохраняем файл на сервере
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.Redirect(http.StatusFound, "/?message=Ошибка во время сохранения")
		log.Println(err)
		return
	}

	// Все прошло успешно
	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно загружен.", "name": name})
}

func main() {
	r := setupRouter()
	utils.PopulateFromCache()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8081")
}

func isArchive(file *multipart.FileHeader) bool {
	allowedExtensions := []string{".zip", ".tar", ".gz", ".rar"} // Укажите допустимые расширения архивов

	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	for _, ext := range allowedExtensions {
		if ext == fileExtension {
			return true
		}
	}

	return false
}
