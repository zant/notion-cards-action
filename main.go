package main

import (
  "context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "regexp"

  "github.com/dstotijn/go-notion"
  "github.com/go-playground/webhooks/v6/github"
  "github.com/joho/godotenv"
)

// Extracts last 32 digits
func getIdFromUrl(page string) string {
  return page[len(page)-32:]
}

type Page struct {
  Name string
}

type CardStatus string

func extractNotionLink(body string) string {
  markdownRegex := regexp.MustCompile(`\[[^][]+]\((https?://(www.notion.so|notion.so)[^()]+)\)`)
  results := markdownRegex.FindAllStringSubmatch(body, -1)
  if len(results) < 1 {
    log.Fatalf("No Notion URL was found")
  } else if len(results) > 1 {
    log.Fatalf("Please link only one Notion URL")
  }
  return results[0][1]
}

const (
  CardStatusCodeReview CardStatus = "Code Review"
  CardStatusQATesting  CardStatus = "QA Testing"
  CardStatusReleased   CardStatus = "Released"
)

func check(err error) {
  if err != nil {
    log.Fatalf("Error: %s", err)
  }
}

func main() {
  godotenv.Load()
  notionClient := notion.NewClient(os.Getenv("NOTION_KEY"))

  payload := github.PullRequestPayload{}

  path := os.Getenv("GITHUB_EVENT_PATH")
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println(path, "Does not exists")
  }

  data, err := ioutil.ReadFile(path)
  check(err)

  json.Unmarshal(data, &payload)

  body := payload.PullRequest.Body
  url := extractNotionLink(body)

  pageId := getIdFromUrl(url)
  databasePageProperties := &notion.DatabasePageProperties{"Status": notion.DatabasePageProperty{Select: &notion.SelectOptions{Name: string(CardStatusCodeReview)}}}
  params := notion.UpdatePageParams{DatabasePageProperties: databasePageProperties}
  page, err := notionClient.UpdatePageProps(context.Background(), pageId, params)
  check(err)

  properties := page.Properties.(notion.DatabasePageProperties)
  status := properties["Status"].Select.Name
  title := properties["Name"].Title[0].Text.Content

  log.Println("\""+title+"\"", "successfully updated to:", status)
}
