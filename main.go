package main

import (
	// "bytes"
	"context"
	"fmt"
	"html/template"

	// "image"
	// "io/ioutil"
	"log"
	"net/http"

	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

const portNumber = ":8080"

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
	_, bucketName, client := ConnectToAWS()
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
	fmt.Printf("Datatype of file : %T\n", file)
	// file資料型態 : multipart.sectionReadCloser
	// hearder資料型態 : *multipart.FileHeader
	fmt.Printf("Datatype of header : %T\n", header)
	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	fileName := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + fileExt
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

}

func ConnectToAWS() (string, string, *s3.Client) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// read the fiele
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	AWS_REGION := viper.GetString("AWS_REGION")
	AWS_ACCESS_KEY := viper.GetString("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := viper.GetString("AWS_SECRET_ACCESS_KEY")
	AWS_BUCKET_NAME := viper.GetString("AWS_BUCKET_NAME")

	staticProvider := credentials.NewStaticCredentialsProvider(
		AWS_ACCESS_KEY,
		AWS_SECRET_ACCESS_KEY,
		"",
	)

	// Load the Shared AWS Configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_REGION), config.WithCredentialsProvider(staticProvider))
	fmt.Println(AWS_BUCKET_NAME)
	if err != nil {
		fmt.Println("測試一下")
		log.Fatalln(err)
	}

	// Create an Amazon S3 service client、do operation in s3
	// 一個新的 s3.Client 的指針 (client)
	client := s3.NewFromConfig(cfg)
	fmt.Println(client)
	fmt.Printf("Datatype of client : %T\n", client)

	return AWS_ACCESS_KEY, AWS_BUCKET_NAME, client

}

func main() {
	ConnectToAWS()

	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.GET("/", homePageHandle)
	router.POST("/api/upload/image", imageUploadHandle)
	http.ListenAndServe(portNumber, router)
}

// serverMUX方法
// 靜態檔案的設定
// http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
// http.HandleFunc("/", homePageHandle)
// fmt.Printf("Starting application on port %s ", portNumber)
// http.ListenAndServe(portNumber, nil)
