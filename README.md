# markdownToconfluence

Push markdown files to Confluence Cloud

## 版权声明  如有侵权请告知我会马上删除

本项目的原始作者是 [justmiles](https://github.com/justmiles) 。出于对需求的使用，原功能不能满足我的需求，所以本人在此项目上做了更改，添加的功能如下：

- 修复了深层级文件夹上传时，无法生成对应的文件树。
- 添加了对git变更的支持，只上传发生变化的文件。

原始项目：[go-markdownToconfluence](https://github.com/justmiles/go-markdownToconfluence)

## Installation

- Windows

  Download [the latest release](https://github.com/justmiles/go-markdownToconfluence/releases/download/v3.1.2/go-markdown2confluence_3.1.2_windows_x86_64.tar.gz) and add to your system `PATH`

## Environment Variables

For best practice we recommend you [authenticate using an API token](https://id.atlassian.com/manage/api-tokens).

- CONFLUENCE_USERNAME - username for Confluence Cloud. When using API tokens set this to your full email.
- CONFLUENCE_PASSWORD - API token or password for Confluence Cloud
- CONFLUENCE_ENDPOINT - endpoint for Confluence Cloud, eg `https://mycompanyname.atlassian.net/wiki`

## Usage

```txt
Push markdown files to Confluence Cloud

Usage:
  markdownToconfluence [flags]

Flags:
  -d, --debug                Enable debug logging
  -e, --endpoint string      Confluence endpoint. (Alternatively set CONFLUENCE_ENDPOINT environment variable) (default "https://mydomain.atlassian.net/wiki")
  -x, --exclude strings      list of exclude file patterns (regex) that will be applied on markdown file paths
  -w, --hardwraps            Render newlines as <br />
  -h, --help                 help for markdownToconfluence
  -m, --modified-since int   Only upload files that have modifed in the past n minutes
      --parent string        Optional parent page to next content under
  -p, --password string      Confluence password. (Alternatively set CONFLUENCE_PASSWORD environment variable)
  -s, --space string         Space in which page should be created
  -c, --comment string       Add a comment to the page (optional)
  -t, --title string         Set the page title on upload (defaults to filename without extension)
      --use-document-title   Will use the Markdown document title (# Title) if available
  -u, --username string      Confluence username. (Alternatively set CONFLUENCE_USERNAME environment variable)
      --version              version for markdownToconfluence
```

## Examples

Upload a local directory of markdown files called `markdown-files` to Confluence.

```shell
markdownToconfluence \
  --space 'MyTeamSpace' \
  markdown-files
```

Upload the same directory, but only those modified in the last 30 minutes. This is particurlarly useful for cron jobs/recurring one-way syncs.

```shell
markdownToconfluence \
  --space 'MyTeamSpace' \
  --modified-since 30 \
  markdown-files
```

Upload a single file

```shell
markdownToconfluence \
  --space 'MyTeamSpace' \
  markdown-files/test.md
```

Upload a directory of markdown files in space `MyTeamSpace` under the parent page `API Docs`

```shell
markdownToconfluence \
  --space 'MyTeamSpace' \
  --parent 'API Docs' \
  markdown-files
```

Upload a directory of markdown files in space `MyTeamSpace` under a _nested_ parent page `Docs/API` and _exclude_ mardown files/directories that match `.*generated.*` or `.*temp.md`

```shell
markdownToconfluence \
  --space 'MyTeamSpace' \
  --parent 'API/Docs' \
  --exclude '.*generated.*' \
  --exclude '.*temp.md' \
   markdown-files
```

Upload a directory of markdown files in space `MyTeamSpace` under the parent page  `API Docs` and use the markdown _document-title_ instead of the filname as document title (if available) in Confluence.

```shell
markdownToconfluence \
  --space 'MyTeamSpace' \
  --parent 'API Docs' \
  --use-document-title \
   markdown-files
```

## Enhancements

It is possible to insert Confluence macros using fenced code blocks.
The "language" for this is `CONFLUENCE-MACRO`, exactly like that in all-caps.
Here is an example for a ToC macro using all headlines starting at Level 2:

```markdown
    # Title

    ```CONFLUENCE-MACRO
    name:toc
      minLevel:2
    ```

    ## Section 1
```

In general almost all macros should be possible.
The general syntax is:

```markdown
    ```CONFLUENCE-MACRO
    name:Name of Macro
    attribute:Value of Attribute
      parameter-name:Value of Parameter
      next-parameter:Value of Parameter
    ```
```

So a fully fledged macro could look like:

```markdown
    ```CONFLUENCE-MACRO
    name:toc
    schema-version:1
      maxLevel:5
      minLevel:2
      exclude:Beispiel.*
      style:none
      type:flat
      separator:pipe
    ```
```

Which will translate to:

```XML
<ac:structured-macro ac:name="toc" ac:schema-version="1" >
  <ac:parameter ac:name="maxLevel">5</ac:parameter>
  <ac:parameter ac:name="minLevel">2</ac:parameter>
  <ac:parameter ac:name="exclude">Beispiel.*</ac:parameter>
  <ac:parameter ac:name="style">none</ac:parameter>
  <ac:parameter ac:name="type">flat</ac:parameter>
  <ac:parameter ac:name="separator">pipe</ac:parameter>
</ac:structured-macro>
```
