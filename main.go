package main

import (
  "context"
  "flag"
  "fmt"
  "os"
  "strings"

	"github.com/google/go-github/github"
)

var (
  description bool
  fork        bool
  language    bool
  star        bool
  watch       bool
  style       string

  langRepoMap map[string][]github.Repository
  languageList []string
)

const shieldsPrefix = "https://img.shields.io/github/"

type shieldsType struct {
  name, label string
}

var shieldsTypeFork = shieldsType{"forks", "Fork"}
var shieldsTypeStar = shieldsType{"stars", "Stars"}
var shieldsTypeWatch = shieldsType{"watchers", "Watch"}

func init() {
  flag.BoolVar(&description, "d",   true,   "show description")
  flag.BoolVar(&fork,        "f",   false,  "show forks")
  flag.BoolVar(&language,    "l",   true,   "show language")
  flag.BoolVar(&star,        "s",   true,   "show stars")
  flag.BoolVar(&watch,       "w",   false,  "show watches")
  flag.StringVar(&style,     "stl", "flat", "badget styles: flat | flat-square | for-the-badge | plastic | social")
}

func main() {
  flag.Parse()

  args := flag.Args()
  if len(args) < 1 {
		fmt.Println("Error: username undefined")
    os.Exit(1)
  }
  user := args[0]

  client := github.NewClient(nil)

  options := &github.ActivityListStarredOptions{Sort: "created"}
  fmt.Println(markdownTag(user + "'s Star List", "##"))
  fmt.Println()
  for page := 1; ; page++ {
		options.Page = page

		starredRepos, res, err := client.Activity.ListStarred(context.Background(), user, options)
		if err != nil {
			fmt.Printf("ListStarred: %s\n", err)
		}

		for _, starredRepo := range starredRepos {
			fmt.Println(markDownUrl(*starredRepo.Repository.FullName, *starredRepo.Repository.HTMLURL))

      if star {
        fmt.Printf(markDownImg(shieldsBadget(shieldsTypeStar, *starredRepo.Repository), shieldsTypeStar.name) + "\t")
      }
      if fork {
        fmt.Printf(markDownImg(shieldsBadget(shieldsTypeFork, *starredRepo.Repository), shieldsTypeFork.name) + "\t")
      }
      if watch {
        fmt.Printf(markDownImg(shieldsBadget(shieldsTypeWatch, *starredRepo.Repository), shieldsTypeWatch.name) + "\t")
      }
      if language {
        fmt.Println(markDownImg(sheildsBadgetLanguage(*starredRepo.Repository), "lang"))
      }
      if description {
        fmt.Println(*starredRepo.Repository.Description)
      }
      fmt.Println()
		}

		if page >= res.LastPage {
			break
		}
  }
}

func shieldsBadget(t shieldsType, rep github.Repository) string {
  s := []string{shieldsPrefix, t.name, "/", *rep.Owner.Login, "/", *rep.Name, ".svg?style=", style, "&logo=github&label=", t.label}
  return strings.Join(s, "")
}

func sheildsBadgetLanguage(rep github.Repository) string {
  s := []string{shieldsPrefix, "languages/top/", *rep.Owner.Login, "/", *rep.Name, ".svg?style=", style}
  return strings.Join(s, "")
}

func markdownTag(text string, tag string) string {
  s := []string{tag, " ", text}
  return strings.Join(s, "")
}

func markDownUrl(text string, url string) string {
  s := []string{"[", text, "](", url, ")"}
  return strings.Join(s, "")
}

func markDownImg(src string, tag string) string {
  s := []string{"![", tag, "](", src, ")"}
  return strings.Join(s, "")
}
