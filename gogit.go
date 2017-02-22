package main

import (
	//"fmt"
	"flag"
	"github.com/abocd/gogit/glog"
	"log"
	"strconv"
	"net/http"
	"path"
)

var gitrepo string;

func www(w http.ResponseWriter, r *http.Request) {

}

func main(){
	//前台访问端口
	port := flag.Int("p",7878,"Port");
	//git地址
	gitrepo = flag.String("r","","Git path");
	flag.Parse();
	//fmt.Printf("%s %d",*path,*port);
	glog.Info("Start","Git目录:%s,浏览器访问 IP:%d",*gitrepo,*port)

	if(path.)

	http.HandleFunc("/", www)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
