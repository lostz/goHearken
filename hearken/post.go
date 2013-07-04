package hearken

import (
    "github.com/knieriem/markdown"
    "bytes"
    "io/ioutil"
    "time"
    "fmt"
    "os"
    "io"
    "strconv"
    "text/template"
    "sort"
)

type Post struct {
    Link string
    Title string
    Date  int
    DateFormat string
    Html  string
}

type Index struct {
    Pre  string
    Next string
    Posts []Post

}

type PostSort []Post
func (s PostSort ) Len() int {return len(s) }
func (s PostSort) Less(i,j int) bool {return s[j].Date<s[i].Date}
func (s PostSort) Swap(i,j int) {s[i],s[j] = s[j],s[i]}
func generateIndexes(webroot string,posts []Post,perPage int ){
    countPosts :=len(posts)
    sort.Sort(PostSort(posts))
    var indexHtml *os.File
    pre :=""
    next :=""
    postsPage := []Post{}
    for i:=0;i<countPosts;i+=perPage {
        page:=(i/perPage)
        if page==0 {
            indexHtml,_= os.Create(webroot+"index.html")
        } else {
            indexHtml,_= os.Create(webroot+"index_"+strconv.Itoa(page)+".html")
            if page==1 {
                pre="./index.html"
            } else {
                pre="./index_"+strconv.Itoa(page-1)+".html"
            }
        }
        defer indexHtml.Close()
        if i<countPosts {
            if (i+1)*perPage>=countPosts {
                postsPage=posts[i:countPosts]
                next=""
            } else {
                next="./index_"+strconv.Itoa(page+1)+".html"
                postsPage=posts[i:(i+1)*perPage]

            }
            indexParser,err :=template.ParseFiles("./templates/header.html","./templates/index.html","./templates/footer.html")
            if err!=nil {
                fmt.Println(err)
            }
            indexInfo :=Index{Pre:pre,Next:next,Posts:postsPage}
            indexParser.ExecuteTemplate(indexHtml,"index",indexInfo)
            pre =""
            next =""
       }
    }

}



func ConvertDate(date string) string {
     timeLayout := "200601021504"
     ztime,err :=time.Parse(timeLayout,date)
     if err!=nil{
         fmt.Println(err)
     } 
     return ztime.Format("2006年1月2日")
}


func GeneratePost(webroot string,perPage int ) {
    mds,_ :=ioutil.ReadDir("./posts")
    posts := []Post{}
    err :=os.RemoveAll(webroot)
    if err!=nil {
    fmt.Println(err)
    }
    err = os.MkdirAll(webroot,os.ModePerm)
    if err!=nil {
    fmt.Println(err)
    }
    copyStatic(webroot)
    for _,md :=range mds {
        contents :=""
        lineNum :=0
        titleLineNum :=0
        contentLineNum :=0
        title :=""
        date  :=""
        if md.Name()[0]=='.'{
            continue
        }
        mdData,err :=ioutil.ReadFile("./posts/"+md.Name())
        if (err!=nil) {
            fmt.Printf("Got error %s",err)
        }

        for i := range mdData {
            if string(mdData[i])=="\n"{
                lineNum +=1
            }
            if lineNum==1&&titleLineNum==0 {
                title =string(mdData[:i])
                titleLineNum=i
            }
            if lineNum==2 {
                date = string(mdData[titleLineNum+1:i])
                contentLineNum =i
                break
            }
        }
        contents =string(mdData[contentLineNum:])
        htmlContents := convert2Html(contents)
        dateFormat :=ConvertDate(date)
        dateInt,_ := strconv.Atoi(date)
        post :=Post{Link:"./"+date+".html",Title:title,Date:dateInt,DateFormat:dateFormat,Html:htmlContents}
        htmlFile,err := os.Create(webroot+date+".html")
        if err!=nil {
            fmt.Println("html create error")
        }
        defer htmlFile.Close()
        postParser,err :=template.ParseFiles("./templates/header.html","./templates/post.html","./templates/footer.html")
        if err!=nil {
        fmt.Println(err)
        }
        postParser.ExecuteTemplate(htmlFile,"post",post)
        posts =append(posts,post)
    }
    generateIndexes(webroot,posts,perPage)
}

func copyStatic(webroot string) {
    statics ,_ :=ioutil.ReadDir("./static")
    for _,static :=range statics {
        srcStatic,err := os.Open("./static/"+static.Name())
        if err!=nil {
        fmt.Println(err)
        }
        defer srcStatic.Close()
        dstStatic,err := os.Create(webroot+static.Name())
        if err!=nil {
            fmt.Println(err)
        }
        defer dstStatic.Close()
        io.Copy(dstStatic,srcStatic)
    }
}

func convert2Html(contents string) string {
    mdParser := markdown.NewParser(&markdown.Extensions{Smart: true})
    buf := bytes.NewBuffer(nil)
    mdParser.Markdown(bytes.NewBufferString(contents), markdown.ToHTML(buf))
    contentHtml := buf.String()
    return contentHtml
}




