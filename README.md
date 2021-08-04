# Notion Cards GH Action

This action moves the _only_ linked Notion Card URL to a "Code Review" status

## Example usage

```
on: [pull_request]

jobs:
  update_card:
    runs-on: ubuntu-latest
    name: A test
    steps:
      - name: Updates card
        uses: zant/notion-cards-action@main
        env:
          NOTION_KEY: --- 
          NOTION_DATABASE_ID: ---
```