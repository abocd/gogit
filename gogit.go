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
	"os/exec"
	//"bytes"
	"strings"
	"fmt"
)

var dirpath string;

func www(w http.ResponseWriter, r *http.Request) {
	//f,err := exec.Command("cd",dirpath).Output();
	result,_ := exec.LookPath(dirpath);
	glog.Info("Log",result);

	f,err := exec.Command("git","log").Output();
	if err != nil{
		glog.Error("Cmd Error",err.Error());
		return;
	}
	data := strings.Split(string(f),"\n");
	fmt.Println(data);
	glog.Info("Cmd","%d (%s)",len(data),data[0]);
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
	glog.Info("Start","Git目录:%s,浏览器访问 IP:%d",dirpath,*port)
	http.HandleFunc("/", www)
	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		glog.Error("Bad","ListenAndServe:%s", err);
		return;
	}
}
