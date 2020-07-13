package meta

//文件元信息结构体
type FileMeta struct {
	FileSha1 string //文件的唯一标志
	FileName string
	FileSize int64
	Location string
	UploadAt string //time stamp
}

/*
func functionName(parameter_list) (return_value_list) {
   …
}
*/

//定义了一个map类型的fileMetas,用来保存所有文件的元信息，这个map的key是string类型，value是FileMeta类型
var fileMetas map[string]FileMeta //map是go中的集合

func init() {
	fileMetas = make(map[string]FileMeta) //initialize map with make function
}

//update key-value pairs in map
//use FileSha1 to find FileMeta

func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//get file msg with filesha1
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
