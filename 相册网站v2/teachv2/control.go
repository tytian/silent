package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//控制器

// 主页面
func IndexView(w http.ResponseWriter, r *http.Request) {
	html := loadHtml("./views/index.html")
	_, err := w.Write(html)
	if err != nil {
		log.Println(err)
	}
}

// 上传页面
func UploadView(w http.ResponseWriter, r *http.Request) {
	html := loadHtml("./views/upload.html")
	_, _ = w.Write(html)
}

// 上传多张页面
func UploadMoreView(w http.ResponseWriter, r *http.Request) {
	html := loadHtml("./views/uploadmore.html")
	_, _ = w.Write(html)
}

// 图片上传
func ApiUpload(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	f, h, err := r.FormFile("file")
	if err != nil {
		_, _ = io.WriteString(w, "上传错误")
		return
	}
	t := h.Header.Get("Content-Type")
	if !strings.Contains(t, "image") {
		_, _ = io.WriteString(w, "文件类型错误")
		return
	}
	_ = os.Mkdir("./static", 0666)
	now := time.Now()
	name := now.Format("2006-01-02-150405") + path.Ext(h.Filename) //获取后缀名
	out, err := os.Create("./static/" + name)
	if err != nil {
		_, _ = io.WriteString(w, "文件创建错误")
		return
	}
	_, _ = io.Copy(out, f)
	out.Close()
	_ = f.Close()
	mod := Info{
		Name: h.Filename,
		Path: "/static/" + name,
		Note: r.Form.Get("note"),
		Unix: now.Unix(),
	}
	//log.Println(InfoAdd(&mod))
	if err = InfoAdd(&mod); err != nil {
		log.Println("error: info add fail!",err)
		return
	}
	http.Redirect(w, r, "/list", 302)
}

//一次上传多张图片
func ApiUploadMore(w http.ResponseWriter, r *http.Request) {
	//设置内存大小
	_ = r.ParseMultipartForm(32 << 20)
	//获取上传的文件组
	files := r.MultipartForm.File["file"]
	len := len(files)
	for i := 0; i < len; i++ {
		//log.Println("more走了一次")
		//打开上传文件
		f, err := files[i].Open()
		if err != nil {
			log.Fatal(err)
		}
		t := files[i].Header.Get("Content-Type")
		if !strings.Contains(t, "image") {
			_, _ = io.WriteString(w, "文件类型错误")
			return
		}
		_ = os.Mkdir("./static", 0666)
		now := time.Now()
		name := now.Format("2006-01-02-150405") + files[i].Filename// +path.Ext(files[i].Filename) //获取后缀名
		out, err := os.Create("./static/" + name)
		if err != nil {
			_, _ = io.WriteString(w, "文件创建错误")
			return
		}
		_, _ = io.Copy(out, f)
		out.Close()
		_ = f.Close()
		mod := Info{
			Name: files[i].Filename,
			Path: "/static/" + name,
			Note: r.Form.Get("note"),
			Unix: now.Unix(),
		}
		if err = InfoAdd(&mod); err != nil {
			log.Println("error: info add fail!",err)
			return
		}
	}
	http.Redirect(w, r, "/list", 302)
}

// 详细页面
func DetailView(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	idStr := r.Form.Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	mod, _ := InfoGet(id)
	html := loadHtml("./views/detail.html")
	date := time.Unix(mod.Unix, 0).Format("2006-01-02 15:04:05")
	html = bytes.Replace(html, []byte("@src"), []byte(mod.Path), 1)
	html = bytes.Replace(html, []byte("@note"), []byte(mod.Note), 1)
	html = bytes.Replace(html, []byte("@unix"), []byte(date), 1)
	_, _ = w.Write(html)
}

// 相册列表
func ApiList(w http.ResponseWriter, r *http.Request) {
	mods, _ := InfoList()
	buf, _ := json.Marshal(mods)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(buf)
}

// 删除
func ApiDrop(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	idStr := r.Form.Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := InfoDrop(id)
	if err != nil {
		_, _ = io.WriteString(w, "数据库删除失败")
		return
	}
	_, _ = io.WriteString(w, "删除成功")
	return
}

// 列表页你
func ListView(w http.ResponseWriter, r *http.Request) {
	html := loadHtml("./views/list.html")
	_, _ = w.Write(html)
}

// 加载html文件
func loadHtml(name string) []byte {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return []byte("")
	}
	return buf
}
