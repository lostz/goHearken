package hearken

import (
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
)

func ProjectInit(projectName string,projectRoot string) {
    projectDir:=projectRoot+"/"+projectName
    _, err := os.Stat(projectDir)
    if err==nil || !os.IsNotExist(err){
        fmt.Printf("project dir: %s Exists",projectDir)
        os.Exit(0)
    }
    err =os.MkdirAll(projectDir,os.ModePerm)
    if err!=nil {
        fmt.Println("innit project errror")
        fmt.Println(err)
        os.Exit(0)
    }
    fmt.Printf("project dir %s init finish",projectDir)
    projectPostDir:=projectDir+"/posts"
    //projectStaticDir:=projectName+"/static"
    //projectTemplateDir:=projectName+"/templates"
    //projectDir :=[...]string{projectPostDir,projectStaticDir,projectTemplateDir}
    //for i:= range projectDir {
    //    err = os.MkdirAll(projectDir[i],os.ModePerm)
    //    if err!=nil {
    //        fmt.Println("make error")
    //    }
   // }
    err = os.MkdirAll(projectPostDir,os.ModePerm)
    if err!=nil {
        fmt.Println("mkdir post error")
    }
    projectConfig :=projectDir+"/configure"
    _, err = os.Create(projectConfig)
    if err!=nil {
        fmt.Println("config error")
    }


}



type PostInfoObject struct {
     PerPage int
     WebRoot string
     RootUrl string
}



func ReadConfig() PostInfoObject {
    configFile :="configure"
    var  postInfoObject PostInfoObject
    fl,err :=ioutil.ReadFile(configFile)
    if err!=nil {
    fmt.Println("file error")
    }
    json.Unmarshal(fl,&postInfoObject)
    return postInfoObject
}

