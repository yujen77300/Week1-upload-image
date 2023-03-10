package main

import (
	// "bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	// "os"

	// "image"
	// "io/ioutil"
	"log"
	"net/http"

	"path/filepath"
	"strings"

	// "github.com/joho/godotenv"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const portNumber = ":8080"

var pictureID int32 = 0

func homePageHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	parsedTemplate, err := template.ParseFiles("./templates/home.tmpl")
	if err != nil {
		fmt.Println("templage parsefile failed,err:", err)
		return
	}
	// if want to share data with tmpl, adding the second argument
	parsedTemplate.Execute(w, nil)
}

func imageUploadHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	region, bucketName, client := ConnectToAWS()
	fmt.Println("到post方法裡面了")
	fmt.Println(client)
	file, header, err := r.FormFile("form")
	if err != nil {
		fmt.Println("測試接收檔案錯誤")
		fmt.Println(err)
		return
	}
	fmt.Println("先先來測試上傳圖片")
	fmt.Println(file)
	// fmt.Printf("Datatype of file : %T\n", file)
	// file資料型態 : multipart.sectionReadCloser
	// hearder資料型態 : *multipart.FileHeader

	//亂數處理用math/rand
	var alphabet []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alphabetSize := len(alphabet)
	var sb strings.Builder
	// 20碼的隨機字串
	for i := 0; i < 20; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	randomFileName := sb.String()
	fmt.Printf("Datatype of file : %s\n", randomFileName)

	fileExt := filepath.Ext(header.Filename)
	fileName := randomFileName + fileExt

	_, error := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    "public-read",
	})
	if error != nil {
		fmt.Printf("Couldn't upload file, Here's why: %v\n", error)
		return
	}

	pictureID += 1
	fmt.Println("測試圖片id")
	fmt.Println(pictureID)

	// 取得url
	url := "https://" + bucketName + ".s3." + region + ".amazonaws.com/" + fileName
	fmt.Println("測試url名稱")
	fmt.Println(url)
	cloudFrontUrl := "https://d1uumvm880lnxp.cloudfront.net/" + fileName

	// 測試文字
	textValue := r.PostFormValue("text")
	fmt.Println("接受文字")
	fmt.Println(textValue)

	// 生成jason檔案

	type UploadInfo struct {
		InfoId        int32
		ImageUrl      string
		Text          string
		CloudFrontUrl string
	}

	uploadInfo := &UploadInfo{
		InfoId:        pictureID,
		ImageUrl:      url,
		Text:          textValue,
		CloudFrontUrl: cloudFrontUrl,
	}

	data, dataError := json.Marshal(uploadInfo)
	if dataError != nil {
		fmt.Printf("json.Marchal failed : %v\n", dataError)
	}
	// fmt.Printf("測試json結果第一次")
	// fmt.Println(string(data))
	// 上傳到rds
	db, _ := ConnectToMYSQL()
	InsertUser(db, url, textValue)
	// 寫入 w
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

func allFileHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, _ := ConnectToMYSQL()
	// fmt.Printf("在allFileHandle裡面的結果")
	res := QueryAllFile(db)
	// fmt.Println(string(res))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(res)))
}

func ConnectToAWS() (string, string, *s3.Client) {
	// 透過viper讀取環境變數
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// 透過dotenv讀取環境變數
	// godotenv.Load()
	AWS_REGION := viper.GetString("AWS_REGION")
	AWS_ACCESS_KEY := viper.GetString("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := viper.GetString("AWS_SECRET_ACCESS_KEY")
	AWS_BUCKET_NAME := viper.GetString("AWS_BUCKET_NAME")
	// AWS_REGION := os.Getenv("AWS_REGION")
	// AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY_ID")
	// AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	// AWS_BUCKET_NAME := os.Getenv("AWS_BUCKET_NAME")
	// fmt.Println("測試測試")
	// fmt.Println(testAWS_REGION)
	// fmt.Println(testAWS_ACCESS_KEY)
	// fmt.Println(testAWS_SECRET_ACCESS_KEY)
	// fmt.Println(testAWS_BUCKET_NAME)

	staticProvider := credentials.NewStaticCredentialsProvider(
		AWS_ACCESS_KEY,
		AWS_SECRET_ACCESS_KEY,
		"",
	)

	// Load the Shared AWS Configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_REGION), config.WithCredentialsProvider(staticProvider))
	fmt.Println(AWS_BUCKET_NAME)
	if err != nil {
		// fmt.Println("測試一下")
		log.Fatalln(err)
	}

	// Create an Amazon S3 service client、do operation in s3
	// 一個新的 s3.Client 的指針 (client)
	client := s3.NewFromConfig(cfg)
	fmt.Println(client)
	fmt.Printf("Datatype of client : %T\n", client)

	return AWS_REGION, AWS_BUCKET_NAME, client

}

func ConnectToMYSQL() (*sql.DB, error) {
	// godotenv.Load()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	const (
		NETWORK = "tcp"
		PORT    = 3306
	)

	// TEXTUSERNAME := os.Getenv("USERNAME")
	// PASSWORD := os.Getenv("PASSWORD")
	// DATABASE := os.Getenv("DATABASE")
	// SERVER := os.Getenv("SERVER")
	USERNAME := viper.GetString("USERNAME")
	PASSWORD := viper.GetString("PASSWORD")
	DATABASE := viper.GetString("DATABASE")
	SERVER := viper.GetString("SERVER")
	// fmt.Println("測試測試")
	// fmt.Println(TEXTUSERNAME)
	// fmt.Println(testPASSWORD)
	// fmt.Println(testDATABASE)
	// fmt.Println(testSERVER)

	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

	db, err := sql.Open("mysql", conn)
	fmt.Printf("db的資料類型")
	fmt.Printf("Datatype of file : %T\n", db)
	if err != nil {
		return nil, fmt.Errorf("開啟 MySQL 連線發生錯誤，原因為： %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("資料庫連線錯誤，原因為： %v", err)
	}
	return db, nil
}

func main() {
	ConnectToAWS()
	// ConnectToMYSQL()

	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.GET("/", homePageHandle)
	router.GET("/api/allfile", allFileHandle)
	router.POST("/api/upload/image", imageUploadHandle)
	http.ListenAndServe(portNumber, router)
}

// serverMUX方法
// 靜態檔案的設定
// http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
// http.HandleFunc("/", homePageHandle)
// fmt.Printf("Starting application on port %s ", portNumber)
// http.ListenAndServe(portNumber, nil)

func InsertUser(DB *sql.DB, image, text string) error {
	_, err := DB.Exec("INSERT INTO information(imageUrl,textInfo) VALUES(?,?)", image, text)
	if err != nil {
		fmt.Printf("建立檔案失敗，原因是：%v", err)
		return err
	}
	fmt.Println("建立檔案成功！")
	return nil

}

func QueryAllFile(db *sql.DB) []byte {
	allData, err := db.Query("SELECT * FROM information;")
	if err != nil {
		fmt.Printf("查詢資料庫失敗，原因為：%v\n", err)
		return nil
	}
	type MysqlData struct {
		Id       int16
		ImageUrl string
		Text     string
	}
	// 建立一個slice來儲存資料
	var files []MysqlData
	for allData.Next() {
		var file MysqlData
		err := allData.Scan(&file.Id, &file.ImageUrl, &file.Text)
		if err != nil {
			fmt.Printf("映射失敗，原因為：%v\n", err)
			return nil
		}
		file.ImageUrl = strings.Replace(file.ImageUrl, "https://uploadimage-dylan.s3.ap-northeast-1.amazonaws.com/", "https://d1uumvm880lnxp.cloudfront.net/", 1)

		files = append(files, file)
	}
	res, err := json.Marshal(files)
	if err != nil {
		fmt.Printf("轉換JSON失敗，原因為：%v\n", err)
		return nil
	}

	return res
}
