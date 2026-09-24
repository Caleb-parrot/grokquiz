# quizgrok

Local quiz in one terminal: question after question, one miss ends the run. A streak counter is kept in memory for the sitting.

## Path

```
category menu
      |
  TUI (Bubble Tea)
      |
  internal/grok
      |
  pkg/grokipedia Search
      |
  grokipedia.com
```

The game does not run `grokipedia-mcp`, does not call Grokipedia from a browser, and does not store responses. `GetPage` in the library still targets `/api/page`, which currently 404s; questions use search snippets instead.

## Round

1. Player picks a category (Space, History, Nature, Music, Sports, Places, Science, Myths).
2. A topic name from that category is chosen at random, skipping topics already answered this run.
3. `grokipedia.Search` runs with an 8s timeout. The hit whose title matches the topic supplies the fact.
4. Three other topic names from the same category are the wrong choices.
5. A correct answer bumps the streak and draws again. A miss shows the topic and stops the run.

429 and 5xx skip that topic and try another. If a draw still fails, enter retries and esc returns to the menu.

## Later

- Save the best streak under `XDG_STATE_HOME` if a sitting should outlive the process.
- Swap the search adapter if Grokipedia's API changes.
- An offline pack only if the live search is too slow or too noisy to play.
