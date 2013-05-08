package main

import (
    "flag"
    "os"
    "go_hearken/hearken"
)

func main(){
    flag.Parse()
    args :=flag.Args()
    if len(args)==0 ||len(args)>3 {
    os.Exit(1)
    }
    switch args[0] {
        case "init":
            hearken.ProjectInit(args[1],args[2])
        case "post":
            postInfo:=hearken.ReadConfig()
            hearken.GeneratePost(postInfo.WebRoot,postInfo.PerPage)
        }

}
