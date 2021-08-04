package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"

  "github.com/dstotijn/go-notion"
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

func getJson(url string, target interface{}) error {
  fmt.Println(url)
  httpClient := &http.Client{}
  r, err := httpClient.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()
  return json.NewDecoder(r.Body).Decode(target)
}

func handleError(err error) {
  log.Fatalf("Error %s", err)
}

func main() {
  godotenv.Load()
  client := notion.NewClient(os.Getenv("NOTION_KEY"))

  path := os.Getenv("GITHUB_EVENT_PATH")
  payload := github.PullRequestPayload{}
  err := getJson(path, payload)
  handleError(err)

  pageId := getIdFromUrl(cardLinked)
  databasePageProperties := &notion.DatabasePageProperties{"Status": notion.DatabasePageProperty{Select: &notion.SelectOptions{Name: string(CardStatusCodeReview)}}}
  params := notion.UpdatePageParams{DatabasePageProperties: databasePageProperties}
  page, err := client.UpdatePageProps(context.Background(), pageId, params)
  handleError(err)

  // Create Page
  // databasePageProperties := notion.DatabasePageProperties{"title": notion.DatabasePageProperty{Title: []notion.RichText{{Text: &notion.Text{Content: "New card"}}}}}
  // params := notion.CreatePageParams{ParentID: databaseId, ParentType: notion.ParentTypeDatabase, DatabasePageProperties: &databasePageProperties}
  // page, err := client.CreatePage(context.Background(), params)

  properties := page.Properties.(notion.DatabasePageProperties)
  status := properties["Status"].Select.Name
  title := properties["Name"].Title[0].Text.Content

  log.Println("\""+title+"\"", "successfully updated to:", status)
}
