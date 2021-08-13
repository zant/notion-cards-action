package main

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "regexp"

  "github.com/dstotijn/go-notion"
  "github.com/go-playground/webhooks/v6/github"
  "github.com/joho/godotenv"
)

type CardStatus string
type GithubEnvironmentVariable string
type InputDefaults string
type Page struct {
  Name string
}

const (
  NotionKey         GithubEnvironmentVariable = "NOTION_KEY"
  GitHubEventPath   GithubEnvironmentVariable = "GITHUB_EVENT_PATH"
  InputPageProperty GithubEnvironmentVariable = "INPUT_PAGE_PROPERTY"
  InputOnPR         GithubEnvironmentVariable = "INPUT_ON_PR"
  InputOnMerge      GithubEnvironmentVariable = "INPUT_ON_MERGE"
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

// Extracts last 32 digits
func getIdFromUrl(page string) string {
  return page[len(page)-32:]
}

func extractNotionLink(body string) string {
  markdownRegex := regexp.MustCompile(`(https?://)?(www\.notion\.so|notion\.so)/?[^(\s)]+`)
  results := markdownRegex.FindAllStringSubmatch(body, -1)
  if len(results) < 1 {
    log.Fatalf("No Notion URL was found")
  } else if len(results) > 1 {
    fmt.Printf("First URL matched was:", results[0][1])
  }

  return results[0][1]
}

func check(err error) {
  if err != nil {
    log.Fatalf("Error: %s", err)
  }
}

func inputFromEnv(input GithubEnvironmentVariable) string {
  return os.Getenv(string(input))
}

func updateCard(pageId string, key string, value string) {
  notionClient := notion.NewClient(os.Getenv(string(NotionKey)))

  valueToUpdate := notion.DatabasePageProperty{Select: &notion.SelectOptions{Name: value}}

  databasePageProperties := &notion.DatabasePageProperties{key: valueToUpdate}

  params := notion.UpdatePageParams{DatabasePageProperties: databasePageProperties}
  page, err := notionClient.UpdatePageProps(context.Background(), pageId, params)
  check(err)

  properties := page.Properties.(notion.DatabasePageProperties)
  status := properties[key].Select.Name
  title := properties["Name"].Title[0].Text.Content

  log.Println("\""+title+"\"", "successfully updated to:", status)
}

func valueFromEvent(merged bool, closed bool) (string, error) {
  if !merged && !closed {
    return inputFromEnv(InputOnPR), nil
  } else if merged && closed {
    return inputFromEnv(InputOnMerge), nil
  } else {
    return "", errors.New("not supported")
  }
}

func main() {
  godotenv.Load()

  payload := github.PullRequestPayload{}

  path := os.Getenv(string(GitHubEventPath))
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println(path, "noes not exists")
  }

  data, err := ioutil.ReadFile(path)
  check(err)

  json.Unmarshal(data, &payload)

  // Values from PR payload
  body := payload.PullRequest.Body
  merged := payload.PullRequest.Merged
  closed := payload.Action == "closed"

  // What to update based on payload
  key := inputFromEnv(InputPageProperty)
  value, err := valueFromEvent(merged, closed)
  check(err)

  url := extractNotionLink(body)
  pageId := getIdFromUrl(url)
  updateCard(pageId, key, value)
}
