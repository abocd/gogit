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
	"bufio"
	//"errors"
	"io"
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"io/ioutil"
	"github.com/kataras/iris/core/errors"
	"time"
)

var dirpath string
var viewLimit int

func www(w http.ResponseWriter, r *http.Request) {
	//f,err := exec.Command("cd",dirpath).Output()
	//fmt.Println(r.RequestURI)
	var html interface{};

		tpl := template.New("index.html")
		tpl.ParseFiles("themes/index.html")
		err := tpl.Execute(w,html)
		if err != nil{
			glog.Error("Html","%s",err)
		}

}


type logData struct{
	Commit string `json:"Commit"`
	Author string `json:"Author"`
	Date string `json:"Date"`
	Memo string `json:"Memo"`
}

func log(w http.ResponseWriter, r *http.Request){
	//f,err := exec.Command("git","log").Output()
	r.ParseForm()
	//page,_ := strconv.Atoi(r.FormValue("page"))
	//fmt.Println(page,"page")
	cmd := exec.Command("git","log")
	cmd.Dir = dirpath //指定command的目录
	f,err :=cmd.Output()
	if err != nil{
		glog.Error("Cmd Error",err.Error())
		return
	}
	data := string(f)
	//fmt.Println(data)
	glog.Info("Cmd","%d (%s)",len(data),data[0])
	//var logRegexp = regexp.MustCompilePOSIX("^commit (.*?)Author: (.*?)Date: (.*?)$")
	var logRegexp = regexp.MustCompile(`commit (\w+)\nAuthor: (.*?)\nDate:   (\w{3} \w{3} \d{2} \d{2}:\d{2}:\d{2} \d{4} [+|-]\d{4})\n{1,}([\s\S]*?)\n`)
	result := logRegexp.FindAllStringSubmatch(data,-1)
	//fmt.Println(result)
	var logList []logData
	//logList = make([]logData)
	timeLayout := "Mon Jan 02 15:04:05 2006 -0700"
	for _,val := range result {
		//fmt.Println("....", val,len(val))
		//Wed Mar 20 10:30:01 2019 +0800
		//Mon Jan 02 15:04:05 2006 -0700
		timeDate,_ := time.Parse(timeLayout,val[3])
		dlog := logData{Commit:val[1],Author:val[2],Date:timeDate.Format("2006-01-02 15:04:05"),Memo:strings.Trim(val[4]," ")};
		logList = append(logList,dlog)
	}
	//fmt.Println(logList)
	w.Write([]byte(_json(logList)))
	//fmt.Fprintf(w,_json(logList))
}

func _json(a interface{}) string{
	jsonList,err := json.Marshal(a)
	if err != nil{
		fmt.Println(err)
		return "{}"
	}
	return string(jsonList)
}

type fileChange struct{
	FileName string `json:"Filename"`
	FileMd5 string `json:"FileMd5"`
	Lines    []string `json:"Lines"`
}

type ViewJson struct{
	Start    int
	StartEnd int
	Files    []fileChange
}

/**

 */
func getCacheContent(from,to string,start,startEnd int) (error,[]byte,string) {
	cacheFile := fmt.Sprintf("%s/%s_%s_%d_%d.cache",cacheDir,from,to,start,startEnd)
	_,err := os.Stat(cacheFile)
	if err == nil{
		if content,err := ioutil.ReadFile(cacheFile);err == nil{
			return nil,content,cacheFile
		}
	}
	return errors.New("内容不存在"),nil,cacheFile
}

func view(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	from := r.FormValue("from")
	to := r.FormValue("to")
	_start := r.FormValue("start")
	if to == ""{
		to = "."
	}
	start,err := strconv.Atoi(_start)
	if err != nil{
		start = 0
	}
	startEnd := start + viewLimit
	fmt.Println("view...",from,to,viewLimit)

	err,content,cacheFile := getCacheContent(from,to,start,startEnd)
	//fmt.Println(err,content)
	if err == nil{

		w.Write(content)
		return
	}

	cmd := exec.Command("git","diff",from,to)
	cmd.Dir = dirpath
	//此处单行输出比较好
	stdout,_ := cmd.StdoutPipe()
	cmd.Start()
	bio := bufio.NewReader(stdout)
	//var line []byte;
	//var hasMoreInLine bool;
	//var err error;
	//line,_ :=bio.ReadString('\n')
	var lines []string
	var fileChangeList []fileChange
	var fileChangeInfo fileChange
	var viewJson ViewJson
	//var isLine = true;
	fileRegexp := regexp.MustCompile(`^diff \-\-git`)
	/**
	 当前行数量
	 */
	var lineNum int
	for {
		line,_,err := bio.ReadLine()
			if err != nil || err == io.EOF{
				break
			}
		newline := string(line)
		//fmt.Print(newline)
		if fileRegexp.MatchString(newline){
			//一个文件开始了
			fileChangeInfo.Lines = lines
			if len(fileChangeInfo.Lines)> 0 {
				fileChangeList = append(fileChangeList, fileChangeInfo)
			}
			//清空line
			fileChangeInfo.FileName = newline
			fileChangeInfo.FileMd5 = cmd5(newline)
			lines = []string{}
			//isLine = false;

		} else {
			if lineNum >= start && lineNum < startEnd {
				lines = append(lines, newline)
			}
			lineNum ++
		}
	}
	lineCount := len(lines)
	if lineCount >0 {
		fileChangeInfo.Lines = lines
		fileChangeList = append(fileChangeList, fileChangeInfo)
	}
	fmt.Println("lineCount",lineCount)
	//fmt.Print(fileChangeList)
	//fmt.Fprintf(w,_json(fileChangeList))
	viewJson.Start = start
	viewJson.StartEnd  = startEnd
	viewJson.Files = fileChangeList
	fileInfo,_ := os.OpenFile(cacheFile,os.O_WRONLY|os.O_CREATE,0777)
	defer fileInfo.Close()
	fileInfo.Write([]byte(_json(viewJson)))
	glog.Info("cacheFile",cacheFile)

	w.Write([]byte(_json(viewJson)))
}

var Tips = `-r gitpath
-p port  default 7878
example ./gogit -p=7878 -r=/var/www/html/phpecorg
`
var cacheDir string

func main(){
	//前台访问端口
	port := flag.Int("p",7878,"Port")
	//git地址
	gitrepo := flag.String("r","","Git path")
	_viewLimit := flag.Int("limit",200,"View limit num")
	viewLimit = *_viewLimit
	flag.Parse()
	dirpath = *gitrepo
	if dirpath ==""{
		glog.Error("Fail","参数错误：\n%s",Tips)
		return
	}
	//fmt.Printf("%s %d",*path,*port)
	dirpath = path.Clean(dirpath)
	fileInfo,err :=os.Stat(dirpath)
	if err != nil{
		glog.Error("Break","%s %s",dirpath,err)
		return
	}
	//fmt.Println(fileInfo)
	if !fileInfo.IsDir(){
		glog.Error("Break","%s 不是一个目录",dirpath)
		return
	}
	_,err2 := os.Stat(path.Clean(dirpath+"/.git"))
	if err2 != nil{
		glog.Error("Break","%s 不是一个有效的git版本库",dirpath)
		return
	}
	dirpathMd5 := cmd5(dirpath)
	cacheDir = getCurrentCacheDir(dirpathMd5)
	glog.Asset("Dir",cacheDir)
	glog.Info("Start","Git目录：%s,浏览器访问：IP:%d",dirpath,*port)
	//fmt.Println(http.Dir("/static/"))
	http.HandleFunc("/log/",log)
	http.HandleFunc("/view/",view)
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", www)
	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		glog.Error("Bad","ListenAndServe:%s", err)
		return
	}
}

func getCurrentCacheDir(filedir string) string {
	dir,err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil{
		glog.Error("Cache",err.Error())
		return ""
	}
	dir = filepath.Clean(fmt.Sprintf("%s/%s/%s",dir,"cache",filedir))
	_,err = os.Stat(dir)
	if err != nil{
		if os.IsNotExist(err){
			err = os.MkdirAll(dir,0777)
			if err != nil{
				glog.Error("Cache",err.Error())
				return ""
			}
		} else {
			glog.Error("Cache",err.Error())
			return ""
		}
	}
	return dir
}

/**
 md5
 */
func cmd5(str string)string{
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
