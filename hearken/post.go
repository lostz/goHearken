package hearken

import (
    . "github.com/russross/blackfriday"
    "io/ioutil"
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
    Date  string
    Html  string
}

type Index struct {
    Pre  string
    Next string
    Posts []Post

}

type PostSort []Post
func (s PostSort) Len() int {return len(s) }
func (s PostSort) Swap(i,j int) {s[i],s[j] = s[j],s[i]}
type ByDate struct { PostSort }

func (s ByDate) Less(i,j int) bool {return s.PostSort[i].Date > s.PostSort[j].Date}

func generateIndexes(webroot string,posts []Post,perPage int ){
    countPosts :=len(posts)
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
            sort.Sort(ByDate{postsPage})
            indexInfo :=Index{Pre:pre,Next:next,Posts:postsPage}
            indexParser.ExecuteTemplate(indexHtml,"index",indexInfo)
            pre =""
            next =""
       }
    }

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
                date =string(mdData[titleLineNum+1:i])
                contentLineNum =i
                break
            }
        }
        contents =string(mdData[contentLineNum:])
        htmlContents := convert2Html(contents)
        post :=Post{Link:"./"+date+".html",Title:title,Date:date,Html:htmlContents}
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



func convert2Html(contents string)   string {
    htmlFlags :=0
    htmlFlags |= HTML_USE_XHTML
	htmlFlags |= HTML_USE_SMARTYPANTS
	htmlFlags |= HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= HTML_SMARTYPANTS_LATEX_DASHES
    renderer := HtmlRenderer(htmlFlags, "", "")
    extensions := 0
	extensions |= EXTENSION_NO_INTRA_EMPHASIS
	extensions |= EXTENSION_TABLES
	extensions |= EXTENSION_FENCED_CODE
	extensions |= EXTENSION_AUTOLINK
	extensions |= EXTENSION_STRIKETHROUGH
	extensions |= EXTENSION_SPACE_HEADERS
    contentHtml := string(Markdown([]byte(contents),renderer,extensions))
    return contentHtml
}



