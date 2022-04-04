# markdownToconfluence

Push markdown files to Confluence Cloud

## 版权声明  如有侵权请告知我会马上删除

本项目的原始作者是 [justmiles](https://github.com/justmiles) 。出于对需求的使用，原功能不能满足我的需求，所以本人在此项目上做了更改，添加的功能如下：

- 修复了深层级文件夹上传时，无法生成对应的文件树。
- 添加了对git变更的支持，只上传发生变化的文件。

原始项目：[go-markdown2confluence](https://github.com/justmiles/go-markdown2confluence)

## 安装

- Windows

  Download [the latest release](https://github.com/wyw-s/markdownToConfluence/releases/download/v4.1.0/markdownToConfluence_4.1.0_windows_x86_64.tar.gz) and add to your system `PATH`

## 环境变量

For best practice we recommend you [authenticate using an API token](https://id.atlassian.com/manage/api-tokens).

| 变量名称                | 默认值                                     | 可选值 | 说明                                                         |
| ----------------------- | ------------------------------------------ | ------ | ------------------------------------------------------------ |
| CONFLUENCE_USERNAME     |                                            |        | Confluence Cloud的用户名。 当使用API令牌时，将此设置为您的完整电子邮件。 |
| CONFLUENCE_PASSWORD     |                                            |        | Confluence Cloud的API令牌或密码                              |
| CONFLUENCE_ENDPOINT     | `https://mycompanyname.atlassian.net/wiki` |        | 你的confluence地址                                           |
| CONFLUENCE_SPACE        | ""                                         |        | 你的团队空间名                                               |
| CONFLUENCE_PARENT       | ""                                         |        | 你需要上传到的confluence父页面                               |
| CONFLUENCE_GIT_SYNC_DIR | ""                                         | docs   | 你本地需要同步的文件夹                                       |
| CONFLUENCE_MODEL        | ""                                         | Git    | 是否基于GIT                                                  |

> 注意：`CONFLUENCE_SPACE`、``CONFLUENCE_PARENT`、`CONFLUENCE_GIT_SYNC_DIR`的使用需要启动 `CONFLUENCE_MODEL`

## Usage

```txt
Push markdown files to Confluence Cloud

Usage:                                                                                                                                                         
  markdown2confluence [flags]                                                                                                                                  
                                                                                                                                                               
Flags:                                                                                                                                                         
  -c, --comment string        (Optional) Add comment to page                                                                                                   
  -d, --debug                 Enable debug logging                                                                                                             
  -e, --endpoint string       Confluence endpoint. (Alternatively set CONFLUENCE_ENDPOINT environment variable) (default "https://mydomain.atlassian.net/wiki")
  -x, --exclude strings       list of exclude file patterns (regex) for that will be applied on markdown file paths                                            
  -g, --git-sync-dir string   Example Set the local synchronization directory                                                                                  
  -w, --hardwraps             Render newlines as <br />                                                                                                        
  -h, --help                  help for markdown2confluence                                                                                                     
      --model string          Is it based on git                                                                                                               
  -m, --modified-since int    Only upload files that have modifed in the past n minutes
      --parent string         Optional parent page to next content under
  -p, --password string       Confluence password. (Alternatively set CONFLUENCE_PASSWORD environment variable)
  -s, --space string          Space in which page should be created
  -t, --title string          Set the page title on upload (defaults to filename without extension)
      --use-document-title    Will use the Markdown document title (# Title) if available
  -u, --username string       Confluence username. (Alternatively set CONFLUENCE_USERNAME environment variable)
      --version               version for markdown2confluence
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
