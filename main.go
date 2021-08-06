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

type CardStatus string
type GithubEnvironmentVariable string
type InputDefaults string

const (
  NotionKey         GithubEnvironmentVariable = "NOTION_KEY"
  GitHubEventPath   GithubEnvironmentVariable = "GITHUB_EVENT_PATH"
  InputPageProperty GithubEnvironmentVariable = "INPUT_PAGE_PROPERTY"
  InputOnPR         GithubEnvironmentVariable = "INPUT_ON_PR"
)

const (
  InputPagePropertyDefault InputDefaults = "Status"
  InputOnPRDefault         InputDefaults = "Code Review"
)

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

func inputFromEnv(input GithubEnvironmentVariable) string {
  env := os.Getenv(string(input))
  fmt.Println(env)
  switch input {
  case InputPageProperty:
    if env == "" {
      return string(InputPagePropertyDefault)
    }
    return env
  case InputOnPR:
    if env == "" {
      return string(InputOnPRDefault)
    }
    return env
  default:
    return ""
  }
}

func main() {
  godotenv.Load()
  notionClient := notion.NewClient(os.Getenv(string(NotionKey)))

  payload := github.PullRequestPayload{}

  path := os.Getenv(string(GitHubEventPath))
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println(path, "Does not exists")
  }

  data, err := ioutil.ReadFile(path)
  check(err)

  json.Unmarshal(data, &payload)

  body := payload.PullRequest.Body
  url := extractNotionLink(body)

  pageId := getIdFromUrl(url)

  inputOnPr := inputFromEnv(InputOnPR)
  propertyToUpdate := notion.DatabasePageProperty{Select: &notion.SelectOptions{Name: inputOnPr}}

  inputPageProperty := inputFromEnv(InputPageProperty)
  databasePageProperties := &notion.DatabasePageProperties{inputPageProperty: propertyToUpdate}

  params := notion.UpdatePageParams{DatabasePageProperties: databasePageProperties}
  page, err := notionClient.UpdatePageProps(context.Background(), pageId, params)
  check(err)

  properties := page.Properties.(notion.DatabasePageProperties)
  status := properties[inputPageProperty].Select.Name
  title := properties["Name"].Title[0].Text.Content

  log.Println("\""+title+"\"", "successfully updated to:", status)
}
