# Notion Cards GH Action

This card updates a property from a page linked in a PR description. Commonly used to update the "Status" property of a card used to keep track of features.

**Variables**
| Key | Description |
|-------|-------|
| NOTION_KEY | The token url retrieved from the token cookie in your browser |

**Inputs**
| Key | Description |
|-------|-------|
| page-property | The name of the property to update. Default is "Status" |
| on-pr | The value of PAGE_PROPERTY to be updated on PR event. Default is "Code Review" |

## Example usage

On PR body:
```
This PR implements [Notion Card](www.notion.so/Card-1234)
```

```
on: [pull_request]

jobs:
  update_card:
    runs-on: ubuntu-latest
    name: Updates Notion Card
    steps:
      - name: Updates to Code Review
        uses: zant/notion-cards-action@main
        env:
          NOTION_KEY: --- 
          NOTION_DATABASE_ID: ---
```