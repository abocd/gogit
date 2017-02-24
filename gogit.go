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
	"regexp"
	"html/template"
	"encoding/json"
)

var dirpath string;

func www(w http.ResponseWriter, r *http.Request) {
	//f,err := exec.Command("cd",dirpath).Output();
	result,_ := exec.LookPath(dirpath);
	glog.Info("Log",result);
	//fmt.Println(r.RequestURI);
	var html interface{};
	var Path = strings.SplitN(r.RequestURI,"?",2);
	var viewRegexp = regexp.MustCompile("^/view");
	var logRegexp = regexp.MustCompile("^/log");
	fmt.Println(Path[0],"Path");
	if viewRegexp.MatchString(Path[0]){
		view();
	} else if logRegexp.MatchString(Path[0]) {
		log(w,r);
	}else{
		tpl := template.New("index.html");
		tpl.ParseFiles("themes/index.html");
		err := tpl.Execute(w,html);
		if err != nil{
			glog.Error("Html","%s",err);
		}
	}

}


type logData struct{
	Commit string `json:commit`
	Author string `json:author`
	Date string `json:date`
	Memo string `json:meo`
}

func log(w http.ResponseWriter, r *http.Request){
	f,err := exec.Command("git","log").Output();
	if err != nil{
		glog.Error("Cmd Error",err.Error());
		return;
	}
	data := string(f);
	fmt.Println(data);
	glog.Info("Cmd","%d (%s)",len(data),data[0]);
	//var logRegexp = regexp.MustCompilePOSIX("^commit (.*?)Author: (.*?)Date: (.*?)$");
	var logRegexp = regexp.MustCompile(`commit (\w+)\nAuthor: (.*?)\nDate:   (\w{3} \w{3} \d{2} \d{2}:\d{2}:\d{2} \d{4} [+|-]\d{4})\n{1,}([\s\S]*?)\n`);
	result := logRegexp.FindAllStringSubmatch(data,-1);
	//fmt.Println(result);
	var logList []logData;
	//logList = make([]logData);
	for _,val := range result {
		fmt.Println("....", val,len(val));
		dlog := logData{Commit:val[1],Author:val[2],Date:val[3],Memo:strings.Trim(val[4]," ")};
		logList = append(logList,dlog);
	}
	//fmt.Println(logList);
	jsonList,_ := json.Marshal(logList);
	fmt.Fprintf(w,string(jsonList));
}

func view(){
	fmt.Println("view...");
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
