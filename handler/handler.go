package handler

/* used to handle some interfaces */
import (
	"encoding/json"
	"fmt"
	"hello/meta"
	"hello/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//upload file interface
//ResponseWriter is the object used to return data to the user
//*http.Request is the object used to receive user requests

//如果是get请求，就显示一个页面。查询，下载用的是get
//如果是post请求，就持续接收文件流，将这个文件保存到本地。上传，删除，修改用的是post
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //back a html page
		data, err := ioutil.ReadFile("./static/view/index.html") //read html file through ReadFile function
		if err != nil {
			io.WriteString(w, "error") //if this html file fails to load, print error information
			return
		}
		io.WriteString(w, string(data)) //if the html file loaded successfully, return the contents of the file
	} else if r.Method == "POST" {
		//receive the file stream uploaded by the user and store it in the local directory
		//接收用户上传的文件流且存储到本地目录

		file, head, err := r.FormFile("file")
		//客户端用的form表单去提交文件的,所以我们这里用FormFile函数去接收
		//FormFile 函数返回文件句柄，文件头，错误信息
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}

		//在这个文件退出之前，要将文件句柄给关掉
		defer file.Close()

		//实例化一个fileMeta用来保存即将上传上来的文件的元信息
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: head.Filename,
			UploadAt: time.Now().Format("2006-01-2 15:04:05"),
		}

		//新建一个文件句柄来准备接收文件流
		newFile, err := os.Create(fileMeta.Location) //新建文件句柄
		if err != nil {
			fmt.Printf("failure to create file,err:%s\n", err.Error())
			return
		}

		//同样的，在这个函数退出之前，要记得先将这个文件句柄关闭
		defer newFile.Close()

		//将内存中的内容，复制到文件buffer中（可以简单理解为，将内存中的信息写入硬盘保存）
		fileMeta.FileSize, err = io.Copy(newFile, file) //如果接收文件流成功，就将这个file拷贝给newFile
		//io.Copy()函数会返回接收的文件的字节长度，还有错误信息
		if err != nil {
			fmt.Printf("Failure to save data into file,err:%s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)                         //将输入位置移动到newFile这个文件的开头位置处
		fileMeta.FileSha1 = util.FileSha1(newFile) //pass the newFile handle to the Filesha1 function, and calculate the hash value of this file

		meta.UpdateFileMeta(fileMeta) //add fileMeta to fileMetas

		//如果文件上传成功了，则显示upload successfully （通过重定向来完成）
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound) //重定向

	}
}

//显示上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload successfully!")
}

//通过文件hash获取文件的元（meta）信息
//get the meta information of the file
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                       //解析客户端请求的操作是什么 解析请求操作
	filehash := r.Form["filehash"][0]   //通过filehash找到客户端想要的 获取到文件的hash值
	fMeta := meta.GetFileMeta(filehash) //通过filehash找到对应的fileMeta，此时还是FileMeta格式

	data, err := json.Marshal(fMeta) //用json.marshal函数将FileMeta格式转换为json格式
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data) //返回转换之后的数据

}

//download interface
//通过文件hash值下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1) //得到对应文件的元信息

	f, err := os.Open(fm.Location) //先通过fm.Location找到这个文件的所在地，
	//然后用os.Open函数，将这个文件在运行内存中打开
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f) //将文件内容写入运行内存中,当然也可以放在其他地方中

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//设置http的响应头，这样就可以让浏览器识别为是文件的下载
	//这里有问题
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment:filename=\""+fm.FileName+"\"")
	w.Write(data) //当文件内容返回给客户端

}

//update meta file
//更新元信息接口（重命名）
func FileUpdateMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                         //解析请求的参数列表,操作，hash，文件名
	opType := r.Form.Get("op")            //操作类型
	fileSha1 := r.Form.Get("filehash")    //要操作的文件hash值
	newFileName := r.Form.Get("filename") //文件名

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" { //如果不是post请求的话
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1) //获取文件的元信息
	curFileMeta.FileName = newFileName        //修改文件名
	meta.UpdateFileMeta(curFileMeta)          //更新

	data, err := json.Marshal(curFileMeta) //解析为json格式
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data) //返回data

}

//delete filemeta interface
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1) //通过hash取出这个文件元信息
	os.Remove(fMeta.Location)           //通过文件的location，达到物理上删除这个文件

	meta.RemoveFileMeta(fileSha1) //这个函数只是删除了文件的元信息，（这句话不是用来真正的删除文件的）

	w.WriteHeader(http.StatusOK)
}
