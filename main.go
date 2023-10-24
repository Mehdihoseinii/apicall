package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
)

type Article struct {
	ID      uint
	Title   string `json:""`
	Desc    string
	Content string
}

func main() {
	handleRequests()
}

func returnAllarticles(ctx echo.Context) error {
	articles, err := GetAllFromCSV()
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, articles)
}

// func returnSingleArticle(ctx echo.Context) error{

// }

func createNewArticle(ctx echo.Context) error {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our articles array.
	article := new(Article)
	if err := ctx.Bind(article); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	// update our global articles array to include
	// our new Article

	addToCSV(*article) // TODO: Handle error-- اینجا بگو که چطور اررو رو تعریف کنم اول
	return ctx.JSON(http.StatusCreated, "hi")
}

func handleRequests() {
	router := echo.New()
	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"message": "hi bitch",
		})
	})
	router.GET("/articles", returnAllarticles)
	router.POST("/articles", createNewArticle)
	// router.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	// router.HandleFunc("/article/{id}", returnSingleArticle)
	fmt.Println("Running webserver on port 80")
	log.Fatal(http.ListenAndServe("localhost:80", router))
}

func (article Article) ToString() []string {
	return []string{
		fmt.Sprint(article.ID),
		article.Title,
		article.Content,
		article.Desc,
	}
}

func addToCSV(article Article) error {
	file, err := os.OpenFile("data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	csvWiter := csv.NewWriter(file)
	err = csvWiter.Write(article.ToString())
	if err != nil {
		return err
	}
	csvWiter.Flush()
	return nil
}

// تمام مقاله هارو میگیری
// آیدی رو از توش پاک میکنی
// اطلاعات قبلی رو پاک میکنی
// اطلاعات جدید رو مینویسی

func deleteFromCSV(ID uint) {
	articles, err := GetAllFromCSV()
	if err != nil {
		fmt.Println(err)
	}
	articles = removeSpecificArticle(ID, articles)
	file, err := os.OpenFile("data.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	if err := csvWriter.WriteAll(ArticlesToRecords(articles)); err != nil {
		fmt.Println(err)
	}
	csvWriter.Flush()
}
func updateFromCSV(ID uint) {
	articles, err := GetAllFromCSV()
	if err != nil {
		fmt.Println(err)
	}

	articles = updateSpificArticle(ID, articles)
	file, err := os.OpenFile("data.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write(ArticlesToRecords(articles)); err != nil {
		fmt.Println(err)
	}
	csvWriter.Flush()

}

func GetAllFromCSV() ([]Article, error) {
	// os.O_APPEND|os.O_CREATE|os.O_WRONLY,
	//os.O_WRONLY|os.O_APPEND, perm
	//باز کردن فایل
	file, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	//بستن فایل
	defer file.Close()
	//خواندن فایل
	CsvReader := csv.NewReader(file)
	records, err := CsvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	articles := []Article{}
	for _, record := range records {
		id, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		article := Article{
			ID:      uint(id),
			Title:   record[1],
			Desc:    record[2],
			Content: record[3],
		}
		articles = append(articles, article)
	}
	return articles, nil
	// []Article
}

func ArticlesToRecords(articles []Article) [][]string {
	records := make([][]string, 0)
	for _, article := range articles {
		records = append(records, article.ToString())
	}
	return records
}

func removeSpecificArticle(ID uint, articles []Article) []Article {
	newArticles := make([]Article, len(articles)-1)
	for _, article := range articles {
		if article.ID != ID {
			newArticles = append(newArticles, article)
		}
	}
	return newArticles
}
func updateSpificArticle(ID uint, articles []Article) []Article {
	newArticles := make([]Article, len(articles)-1)
	for _, article := range articles {
		if article.ID != ID {
			newArticles = append(newArticles, article)
		}
	}
	return newArticles
}

var test string
