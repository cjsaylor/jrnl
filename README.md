# jrnl [![](https://drone.chris-saylor.com/api/badges/cjsaylor/jrnl/status.svg)](https://drone.chris-saylor.com/cjsaylor/jrnl)

A quick and easy CLI journaling tool that uses Github Wiki repos for organizing journal entries.

![](jrnl.gif)
---

[View an example wiki journal](https://github.com/cjsaylor/jrnl/wiki)

---

[![GoDoc](https://godoc.org/github.com/cjsaylor/jrnl?status.svg)](https://godoc.org/github.com/cjsaylor/jrnl)

* [Requirements](#requirements)
* [Installation](#installation)
* [Quick Start](#quick-start)
* [Options](#options)
* [Commands](#commands)
	* [Tag Journal Entries](#tag)
	* [Generate Index](#index)
	* [Append Images](#append-an-image)
* [Tips & Tricks](#tips--tricks)
	* [Use `find` command to create a book](#use-find-command-to-create-a-book)
	* [Use `find` and `tag` commands to add a common tag](#use-find-and-tag-commands-to-add-a-common-tag)
* [Development](#development)

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

### Compiling from source

```bash
go get -u github.com/cjsaylor/jrnl/cmd/jrnl
```

## Quick start

[Install `jrnl`](#installation)

```bash
go get -u github.com/cjsaylor/jrnl/cmd/jrnl
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

## Options

You can configure `jrnl` with the following environment variables:

* `JRNL_EDITOR` (`vim`) - Editor to use
* `JRNL_EDITOR_OPTIONS` (`""`) - Additional CLI flags for your editor. IE, for VS Code: `-n $HOME/journal.wiki/`
* `JOURNAL_PATH` (`~/journal.wiki`) - Path to your cloned Github wiki repo.

## Commands

### Tag

`jrnl` has the ability to tag a journal entry so that it can be easily referenced and found.

```bash
jrnl tag -t sometag
```

This would append "sometag" to the journal entry in `frontmatter` format:

```yaml
---
tags:
- sometag
---
```

### Index

`jrnl` has the ability to generate an `Index.md` that allows you to easily reference any journal entry by a tag.

You can use the `tag` command specified above to specify tags.

```bash
jrnl tag -t cooltag -t another cool tag -d "2017-12-01"
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

```bash
jrnl tag -t cooltag -t -d "2017-12-04"
```

Then `jrnl index` would generate:

> ~/journal.wiki/Index.md

> * *cooltag* [2017-12-01](), [2017-12-04]()
> * *another cool tag* [2017-12-01]()

### Append an image

Many developers use hand-written notes (or a whiteboard) and want to store it in a common journal.

You can use the `jrnl image /path/to/image` command to quickly add an image to the journal repo and append it to the current journal entry.

## Tips & Tricks

### Use Find Command to Create a Book

You can use the `find` command (as of `v0.3.0`) to create a single document of all journal entries related to a tag (or tags):

```bash
jrnl find --tag sometag --tag somerelatedtag | xargs cat > sometag.md
```

### Use `find` command to add a common tag

You can use the `find` (as of `v0.3.0`) command and the `tag` command (as of `v0.4.0`) to add a common tag or tags:

```bash
jrnl tag -t newtag $(jrnl find -tag existingtag | xargs -I {} echo "-f {}")
```

## Development

To compile the CLI tool:

```bash
go build -o jrnl cmd/jrnl/main.go
```

To run the unit tests:

```bash
go test $(go list ./... | grep -v /vendor/)
```

To add a dependency use the [`dep` tool](https://github.com/golang/dep)

For distribution, the [`goreleaser` tool](https://goreleaser.com/) is used. Simply run `goreleaser` to tag and distribute.
