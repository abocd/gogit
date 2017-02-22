package main

import (
	//"fmt"
	"flag"
	"github.com/abocd/gogit/glog"
	//"log"
	"strconv"
	"net/http"
	"path"
	"os"
	//"fmt"
)

var dirpath string;

func www(w http.ResponseWriter, r *http.Request) {

}

func main(){
	//前台访问端口
	port := flag.Int("p",7878,"Port");
	//git地址
	gitrepo := flag.String("r","","Git path");
	flag.Parse();
	dirpath = *gitrepo;
	//fmt.Printf("%s %d",*path,*port);
	fileInfo,err :=os.Stat(dirpath);
	if err != nil{
		glog.Error("Break","%s %s",dirpath,err);
		return;
	}
	//fmt.Println(fileInfo);
	if !fileInfo.IsDir(){
		glog.Error("Break","%s 不是一个目录",dirpath);
		return;
	}
	_,err2 := os.Stat(path.Clean(dirpath+"/.git"));
	if err2 != nil{
		glog.Error("Break","%s 不是一个有效的git版本库",dirpath);
		return;
	}
	http.HandleFunc("/", www)
	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		glog.Error("Bad","ListenAndServe:%s", err);
		return;
	}
	glog.Info("Start","Git目录:%s,浏览器访问 IP:%d",dirpath,*port)
}
