# jrnl

A quick and easy CLI journaling tool that uses Github Wiki repos for organizing journal entries.

![](jrnl.gif)
---

[View an example wiki journal](https://github.com/cjsaylor/jrnl/wiki)

---

[![GoDoc](https://godoc.org/github.com/cjsaylor/jrnl?status.svg)](https://godoc.org/github.com/cjsaylor/jrnl)

* [Quick Start](#quick-start)
* [Requirements](#requirements)
* [Installation](#installation)
* [Options](#options)
* [Generate Index](#index)
* [Append Images](#append-an-image)
* [Operate on a different date](#operate-on-a-different-date)
* [Use `find` command to create a book](#use-find-command-to-create-a-book)
* [Development](#development)

## Quick start

Install `jrnl`:

```bash
go get github.com/cjsaylor/jrnl
```

Clone a github wiki you want to act as the store of your journal:

```bash
git clone https://github.com/<yourname>/<journel_repo>.wiki.git ~/journal.wiki
```

Quickly drop into the editor of your choice (default `vim` but configurable via `$JRNL_EDITOR`):

```bash
jrnl
```

Write your entry and then "memorize":

```bash
jrnl memorize
```

## Requirements

* Git
* Github account (and access to a repo for the wiki)

## Installation

### MacOS

```bash
brew install cjsaylor/tap/jrnl
```

### Windows & Linux

Download the binary from the [latest release](https://github.com/cjsaylor/jrnl/releases/latest)

## Options

You can configure `jrnl` with the following environment variables:

* `JRNL_EDITOR` (`vim`) - Editor to use
* `JRNL_EDITOR_OPTIONS` (`""`) - Additional CLI flags for your editor. IE, for VS Code: `-n $HOME/journal.wiki/`
* `JOURNAL_PATH` (`~/journal.wiki`) - Path to your cloned Github wiki repo.

## Index

`jrnl` has the ability to generate an `Index.md` that allows you to easily reference any journal entry by a tag.

In your journal entry, you can specify tags in `frontmatter` format:

> ~/journal.wiki/entries/2017-12-01.md
```
---
tags: [cooltag, another cool tag]
---

My journal content here
```

Then run the index command:

```bash
jrnl index
```

Which generates:

> ~/journal.wiki/Index.md

> * *cooltag* [2017-12-01]()
> * *another cool tag* [2017-12-01]()

If there is more than one journal entry that uses the same tag:

> ~/journal.wiki/Index.md

> * *cooltag* [2017-12-01](), [2017-12-04]()
> * *another cool tag* [2017-12-01]()

## Append an image

Many developers use hand-written notes (or a whiteboard) and want to store it in a common journal.

You can use the `jrnl image /path/to/image` command to quickly add an image to the journal repo and append it to the current journal entry.

## Operate on a different date

If you need to open or add an image to a different day's entry, you can use the `-date` flag:

```bash
jrnl -date 2017-12-05 open
```

Additionally, certain commands like `open` and `image` support a subject `-s` flag which allows you to name the file anything you want.

## Use Find Command to Create a Book

You can use the `find` command (as of `v0.3.0`) to create a single document of all journal entries related to a tag (or tags):

```bash
jrnl find --tag sometag --tag somerelatedtag | xargs cat > sometag.md
```

## Development

To compile the CLI tool:

```bash
go build -o jrnl cmd/cli/main.go
```

To add a dependency use the [`dep` tool](https://github.com/golang/dep)

For distribution, the [`goreleaser` tool](https://goreleaser.com/) is used. Simply run `goreleaser` to tag and distribute.