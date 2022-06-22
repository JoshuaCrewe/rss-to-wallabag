package main

// https://www.reddit.com/r/golang/comments/cs55ul/comment/excm8ts/?utm_source=share&utm_medium=web2x&context=3
import (
    "fmt"
    "github.com/integrii/flaggy"
    "joshuaCrewe/go-rss-wallabag/src"
)

var version = "0.1"

var commandInit   *flaggy.Subcommand
var commandAdd    *flaggy.Subcommand
var commandRemove *flaggy.Subcommand

var tags string

func init() {
  /*
  * GENERAL CONFIG
  */
  flaggy.SetName("Go RSS to Wallabag")
  flaggy.SetDescription("Send articles from RSS feeds to your wallabag instance")

  /*
  * SUB COMMANDS
  */
  commandInit = flaggy.NewSubcommand("init")
  commandInit.Description = "Configure your secrets"
  flaggy.AttachSubcommand(commandInit, 1)

  commandAdd = flaggy.NewSubcommand("add")
  commandAdd.Description = "Add a new RSS feed"
  commandAdd.String(&tags, "t", "tags", "Define comma separated list of tags")
  flaggy.AttachSubcommand(commandAdd, 1)

  commandRemove = flaggy.NewSubcommand("remove")
  commandRemove.Description = "Remove an RSS feed"
  flaggy.AttachSubcommand(commandRemove, 1)
  
  /*
  * BOOTSTRAP
  */
  flaggy.SetVersion(version)
  flaggy.Parse()
}

func main() {
     if commandInit.Used {
        fmt.Println("🐛 RUN INIT")
        gobag.Init()
    } else if commandAdd.Used {
        fmt.Println("🐛 ADD RSS Feed ", tags)
        gobag.Add()
    } else if commandRemove.Used {
        fmt.Println("🐛 Remove RSS Feed")
        gobag.Remove()
    } else {
        gobag.Run()
    }
}
