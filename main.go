package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strings"

  "github.com/go-playground/webhooks/v6/github"
  "github.com/joho/godotenv"
)

func getIdFromUrl(page string) string {
  return strings.Split(page, "Card-")[1]
}

type Page struct {
  Name string
}

type CardStatus string

const cardLinked = "https://www.notion.so/Card-ef2e02bf5f0f4a37a6c7fe48ff5de280"

const (
  CardStatusCodeReview CardStatus = "Code Review"
  CardStatusQATesting  CardStatus = "QA Testing"
  CardStatusReleased   CardStatus = "Released"
)

func check(err error) {
  log.Fatalf("Error: %s", err)
}

func main() {
  godotenv.Load()
  // client := notion.NewClient(os.Getenv("NOTION_KEY"))

  path := os.Getenv("GITHUB_EVENT_PATH")
  payload := github.PullRequestPayload{}
  data, err := ioutil.ReadFile(path)
  check(err)

  json.Unmarshal(data, &payload)

  fmt.Println(payload.PullRequest.Body)

  // pageId := getIdFromUrl(cardLinked)
  // databasePageProperties := &notion.DatabasePageProperties{"Status": notion.DatabasePageProperty{Select: &notion.SelectOptions{Name: string(CardStatusCodeReview)}}}
  // params := notion.UpdatePageParams{DatabasePageProperties: databasePageProperties}
  // page, err := client.UpdatePageProps(context.Background(), pageId, params)
  // check(err)

  // Create Page
  // databasePageProperties := notion.DatabasePageProperties{"title": notion.DatabasePageProperty{Title: []notion.RichText{{Text: &notion.Text{Content: "New card"}}}}}
  // params := notion.CreatePageParams{ParentID: databaseId, ParentType: notion.ParentTypeDatabase, DatabasePageProperties: &databasePageProperties}
  // page, err := client.CreatePage(context.Background(), params)

  // properties := page.Properties.(notion.DatabasePageProperties)
  // status := properties["Status"].Select.Name
  // title := properties["Name"].Title[0].Text.Content

  // log.Println("\""+title+"\"", "successfully updated to:", status)
}
