package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/github"
)

var (
	description bool
	fork        bool
	language    bool
	star        bool
	watch       bool
	sortBy      string
	style       string

	user string

	langList    []string
	langRepoMap map[string][]github.Repository
)

const shieldsPrefix = "https://img.shields.io/github/"

type shieldsType struct {
	name, label string
}

var shieldsTypeFork = shieldsType{"forks", "Fork"}
var shieldsTypeStar = shieldsType{"stars", "Stars"}
var shieldsTypeWatch = shieldsType{"watchers", "Watch"}

func init() {
	flag.BoolVar(&description, "d", false, "hide description")
	flag.BoolVar(&fork, "f", false, "hide forks")
	flag.BoolVar(&language, "l", false, "hide language")
	flag.BoolVar(&star, "s", false, "hide stars")
	flag.BoolVar(&watch, "w", false, "hide watches")
	flag.StringVar(&sortBy, "sort", "lang", "how to sort the repository list: created | updated | pushed | full_name | lang. default: lang")
	flag.StringVar(&style, "stl", "flat", "badget styles: flat | flat-square | for-the-badge | plastic | social. default: flat")

	langList = []string{}
	langRepoMap = map[string][]github.Repository{}
}

func main() {
	checkFlag()
	fetchData()
	if sortByByLang() {
		printByLang()
	}
}

func checkFlag() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: username undefined")
		os.Exit(1)
	}
	user = args[0]
}

func fetchData() {
	client := github.NewClient(nil)

	sortByTemp := sortBy
	if sortByByLang() {
		sortByTemp = "created"
	}
	options := &github.ActivityListStarredOptions{Sort: sortByTemp}

	fmt.Println(markdownTag(user+"'s Star List", "##"))
	fmt.Println()
	for page := 1; ; page++ {
		options.Page = page

		reps, res, err := client.Activity.ListStarred(context.Background(), user, options)
		if err != nil {
			fmt.Printf("ListStarred: %s\n", err)
		}

		for _, rep := range reps {
			r := *rep.Repository

			lang := "None"
			if r.Language != nil {
				lang = *r.Language
			}

			list, ok := langRepoMap[lang]
			if !ok {
				langList = append(langList, lang)
				list = []github.Repository{}
			}
			list = append(list, r)
			langRepoMap[lang] = list

			if !sortByByLang() {
				printRep(r)
			}
		}

		if page >= res.LastPage {
			break
		}
	}
}

func printRep(rep github.Repository) {
	fmt.Println(markDownUrl(*rep.FullName, *rep.HTMLURL))

	if !star {
		fmt.Printf(markDownImg(shieldsBadget(shieldsTypeStar, rep), shieldsTypeStar.name) + "\t")
	}
	if !fork {
		fmt.Printf(markDownImg(shieldsBadget(shieldsTypeFork, rep), shieldsTypeFork.name) + "\t")
	}
	if !watch {
		fmt.Printf(markDownImg(shieldsBadget(shieldsTypeWatch, rep), shieldsTypeWatch.name) + "\t")
	}
	if !language {
		fmt.Printf(markDownImg(sheildsBadgetLanguage(rep), "lang"))
	}
	fmt.Println()

	fmt.Println()
	if !description {
		fmt.Println(*rep.Description)
	}
	fmt.Println()
}

func printByLang() {
	printLangList()
	printRepsByLang()
}

func printLangList() {
	sort.Strings(langList)
	fmt.Println(markdownTag("Content", "##"))
	fmt.Println()
	for _, lang := range langList {
		fmt.Println(markdownTag(markDownUrl(lang, "#"+strings.ToLower(lang)), "-"))
	}
	fmt.Println()
}

func printRepsByLang() {
	for _, lang := range langList {
		reps, _ := langRepoMap[lang]
		fmt.Println(markdownTag(lang, "###"))
		for _, rep := range reps {
			printRep(rep)
		}
	}
}

func sortByByLang() bool {
	return sortBy == "lang"
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
