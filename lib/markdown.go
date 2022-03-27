package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/justmiles/go-confluence"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	e "markdownToConfluence/lib/extension"
)

const (
	// DefaultEndpoint provides an example endpoint for users
	DefaultEndpoint = "https://mydomain.atlassian.net/wiki"

	// Parallelism determines how many files to convert and upload at a time
	Parallelism = 5
)

// Markdown2Confluence stores the settings for each run
type Markdown2Confluence struct {
	Space               string
	Comment             string
	Title               string
	File                string
	Ancestor            string
	Debug               bool
	UseDocumentTitle    bool
	WithHardWraps       bool
	Since               int
	Username            string
	Password            string
	Endpoint            string
	Parent              string
	SourceMarkdown      []string
	DeleteMarkdown      []string
	ExcludeFilePatterns []string
	client              *confluence.Client
	IsGitEnv            bool
	GitSyncDir          string
}

// CreateClient returns a new markdown clietn
func (m *Markdown2Confluence) CreateClient() {
	m.client = new(confluence.Client)
	m.client.Username = m.Username
	m.client.Password = m.Password
	m.client.Endpoint = m.Endpoint
	m.client.Debug = m.Debug
}

// SourceEnvironmentVariables overrides Markdown2Confluence with any environment variables that are set
//  - CONFLUENCE_USERNAME
//  - CONFLUENCE_PASSWORD
//  - CONFLUENCE_ENDPOINT
func (m *Markdown2Confluence) SourceEnvironmentVariables() {
	var s string
	s = os.Getenv("CONFLUENCE_USERNAME")
	if s != "" {
		m.Username = s
	}

	s = os.Getenv("CONFLUENCE_PASSWORD")
	if s != "" {
		m.Password = s
	}

	s = os.Getenv("CONFLUENCE_ENDPOINT")
	if s != "" {
		m.Endpoint = s
	}

	s = os.Getenv("CONFLUENCE_SPACE")
	if s != "" {
		m.Space = s
	}

	s = os.Getenv("CONFLUENCE_PARENT")
	if s != "" {
		m.Parent = s
	}

	s = os.Getenv("CONFLUENCE_GIT_SYNC_DIR")
	if s != "" {
		m.GitSyncDir = s
	}

	//slice := []string{os.Getenv("CONFLUENCE_FOLDER_NAME")}
	//if slice[0] != "" {
	//	m.SourceMarkdown = slice
	//}
}

// Validate required configs are set
func (m Markdown2Confluence) Validate() error {
	if m.Space == "" {
		return fmt.Errorf("--space is not defined")
	}
	if m.Username == "" {
		return fmt.Errorf("--username is not defined")
	}
	if m.Password == "" {
		return fmt.Errorf("--password is not defined")
	}
	if m.Endpoint == "" {
		return fmt.Errorf("--endpoint is not defined")
	}
	if m.Endpoint == DefaultEndpoint {
		return fmt.Errorf("--endpoint is not defined")
	}
	if len(m.SourceMarkdown) == 0 {
		return fmt.Errorf("please pass a markdown file or directory of markdown files")
	}
	if len(m.SourceMarkdown) > 1 && m.Title != "" {
		return fmt.Errorf("You can not set the title for multiple files")
	}
	return nil
}

func (m *Markdown2Confluence) IsExcluded(p string) bool {
	for _, pattern := range m.ExcludeFilePatterns {
		r := regexp.MustCompile(pattern)
		if r.MatchString(p) {
			fmt.Printf("excluding markdown file '%s': exclude pattern '%s'\n", p, pattern)
			return true
		}
	}

	return false
}

// Run the sync
func (m *Markdown2Confluence) Run() []error {
	var markdownFiles []MarkdownFile
	var deleteMarkdownFiles []MarkdownFile
	var now = time.Now()
	m.CreateClient()

	for _, f := range m.SourceMarkdown {
		file, err := os.Open(f)
		defer file.Close()
		if err != nil {
			return []error{fmt.Errorf("Error opening file %s", err)}
		}

		stat, err := file.Stat()
		if err != nil {
			return []error{fmt.Errorf("Error reading file meta %s", err)}
		}

		var md MarkdownFile

		if stat.IsDir() {

			// prevent someone from accidently uploading everything under the same title
			if m.Title != "" {
				return []error{fmt.Errorf("--title not supported for directories")}
			}

			err := filepath.Walk(f,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if strings.HasSuffix(path, ".md") && !m.IsExcluded(path) {

						// Only include this file if it was modified m.Since minutes ago
						if m.Since != 0 {
							if info.ModTime().Unix() < now.Add(time.Duration(m.Since*-1)*time.Minute).Unix() {
								if m.Debug {
									fmt.Printf("skipping %s: last modified %s\n", info.Name(), info.ModTime())
								}
								return nil
							}
						}

						var tempTitle string
						var tempParents []string

						if strings.HasSuffix(path, "README.md") {
							tempTitle = strings.Split(path, "/")[len(strings.Split(path, "/"))-2]
							tempParents = deleteFromSlice(deleteFromSlice(strings.Split(filepath.ToSlash(filepath.Dir(strings.TrimPrefix(path, f))), "/"), "."), tempTitle)
						} else {
							tempTitle = strings.TrimSuffix(filepath.Base(path), ".md")
							tempParents = deleteFromSlice(strings.Split(filepath.ToSlash(filepath.Dir(strings.TrimPrefix(path, f))), "/"), ".")
						}

						if m.UseDocumentTitle == true {
							doc_title := getDocumentTitle(path)
							if doc_title != "" {
								tempTitle = doc_title
							}
						}

						md = MarkdownFile{
							Path:    path,
							Parents: tempParents,
							Title:   tempTitle,
						}

						if m.Parent != "" {
							parents := strings.Split(m.Parent, "/")
							md.Parents = append(parents, md.Parents...)
							md.Parents = deleteEmpty(md.Parents)
						}

						markdownFiles = append(markdownFiles, md)

					}
					return nil
				})
			if err != nil {
				return []error{fmt.Errorf("Unable to walk path: %s", f)}
			}

		} else {
			if strings.HasSuffix(f, ".md") && !m.IsExcluded(f) {

				// 当前上传的文件属于gitchangelist
				if m.IsGitEnv {
					var tempTitle string
					var tempParents []string

					if strings.HasSuffix(f, "README.md") {
						tempTitle = strings.Split(f, "/")[len(strings.Split(f, "/"))-2]
						if m.GitSyncDir != "" {
							tempParents = deleteFromSlice(deleteFromSlice(strings.Split(filepath.ToSlash(filepath.Dir(strings.TrimPrefix(f, m.GitSyncDir))), "/"), "."), tempTitle)
						} else {
							tempParents = deleteFromSlice(deleteFromSlice(strings.Split(filepath.ToSlash(filepath.Dir(f)), "/"), "."), tempTitle)
						}
					} else {
						tempTitle = strings.TrimSuffix(filepath.Base(f), ".md")
						if m.GitSyncDir != "" {
							tempParents = deleteFromSlice(strings.Split(filepath.ToSlash(filepath.Dir(strings.TrimPrefix(f, m.GitSyncDir))), "/"), ".")
						} else {
							tempParents = deleteFromSlice(strings.Split(filepath.ToSlash(filepath.Dir(f)), "/"), ".")
						}
					}

					if m.UseDocumentTitle == true {
						doc_title := getDocumentTitle(f)
						if doc_title != "" {
							tempTitle = doc_title
						}
					}

					md = MarkdownFile{
						Path:    f,
						Parents: tempParents,
						Title:   tempTitle,
					}
				} else {
					md = MarkdownFile{
						Path:  f,
						Title: m.Title,
					}

					if md.Title == "" {
						if m.UseDocumentTitle == true {
							md.Title = getDocumentTitle(f)
						}
						if md.Title == "" {
							md.Title = strings.TrimSuffix(filepath.Base(f), ".md")
						}
					}
				}

				if m.Parent != "" {
					parents := strings.Split(m.Parent, "/")
					md.Parents = append(parents, md.Parents...)
					md.Parents = deleteEmpty(md.Parents)
				}

				markdownFiles = append(markdownFiles, md)
			}
		}

	}

	for _, f := range m.DeleteMarkdown {
		if strings.HasSuffix(f, ".md") && !m.IsExcluded(f) {
			var md MarkdownFile
			md = MarkdownFile{
				Path:  f,
				Title: m.Title,
			}

			if md.Title == "" {
				if m.UseDocumentTitle == true {
					md.Title = getDocumentTitle(f)
				}
				if md.Title == "" {
					md.Title = strings.TrimSuffix(filepath.Base(f), ".md")
				}
			}

			deleteMarkdownFiles = append(deleteMarkdownFiles, md)
		}
	}
	var errors []error

	for _, markdownFile := range deleteMarkdownFiles {
		var err error
		_, err = markdownFile.DeletePage(m)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		fmt.Printf("删除云端文件：%s \n", markdownFile.Path)
	}

	var (
		wg    = sync.WaitGroup{}
		queue = make(chan MarkdownFile)
	)

	// Process the queue
	for worker := 0; worker < Parallelism; worker++ {
		wg.Add(1)
		go m.queueProcessor(&wg, &queue, &errors)
	}

	for _, markdownFile := range markdownFiles {

		// Create parent pages synchronously
		if len(markdownFile.Parents) > 0 {
			var err error
			markdownFile.Ancestor, err = markdownFile.FindOrCreateAncestors(m)
			if err != nil {
				errors = append(errors, err)
				continue
			}
		}

		queue <- markdownFile
	}

	close(queue)

	wg.Wait()

	return errors
}

func (m *Markdown2Confluence) queueProcessor(wg *sync.WaitGroup, queue *chan MarkdownFile, errors *[]error) {
	defer wg.Done()

	for markdownFile := range *queue {
		url, err := markdownFile.Upload(m)
		if err != nil {
			*errors = append(*errors, fmt.Errorf("Unable to upload markdown file %s: \n\t%s", markdownFile.Path, err))
		}
		fmt.Printf("已上传云端文件：%s %s %s\n", markdownFile.Path, markdownFile.FormattedPath(), url)
	}
}

func validateInput(s string, msg string) {
	if s == "" {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func renderContent(filePath, s string, withHardWraps bool) (content string, images []string, err error) {
	confluenceExtension := e.NewConfluenceExtension(filePath)
	ro := goldmark.WithRendererOptions(
		html.WithXHTML(),
	)
	if withHardWraps {
		ro = goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		)
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.DefinitionList),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		ro,
		goldmark.WithExtensions(
			confluenceExtension,
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(s), &buf); err != nil {
		return "", nil, err
	}

	return buf.String(), confluenceExtension.Images(), nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func deleteFromSlice(s []string, del string) []string {
	for i, v := range s {
		if v == del {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}
	return s
}

func getDocumentTitle(p string) string {
	// Read file to check for the content
	file_content, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	// Convert []byte to string and print to screen
	text := string(file_content)

	// check if there is a
	e := `^#\s+(.+)`
	r := regexp.MustCompile(e)
	result := r.FindStringSubmatch(text)
	if len(result) > 1 {
		// assign the Title to the matching group
		return result[1]
	}

	return ""
}
