package main

import (
	"fmt"
	"hello/handler"
	"net/http"
)

func main() {
	//添加路由

	http.HandleFunc("/file/upload", handler.UploadHandler)        //establish routing rules by HandleFunc function
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler) //显示文件上传成功的路由
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)     //显示文件元信息的路由
	http.HandleFunc("/file/download", handler.DownloadHandler)    //下载文件的路由
	http.HandleFunc("/file/update", handler.FileUpdateMetaHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	//address monitoring
	err := http.ListenAndServe(":8080", nil) //address monitoring
	if err != nil {
		fmt.Println("Failed do start server,err:%s", err.Error()) //if address monitor failure
	}
}
