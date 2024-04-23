package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
	"github.com/yuin/goldmark"
)

func main() {
	createDirIfNotExists("posts")
	createDirIfNotExists("public")

	posts := readMarkdownPosts("posts")
	createIndexPage(posts)

	// c, err := readConf("conf.yaml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("conf: %#v\n", c.Conf.Languages)
}

func layout(title string, content elem.Node) string {
	bodyStyle := styles.Props{
		styles.Margin: "0",
	}

	headerStyle := styles.Props{
		styles.BackgroundColor: "lightblue",
		styles.Padding:         "10px",
		styles.TextAlign:       "center",
		attrs.Class:            "grid grid-cols-2",
	}

	footerStyle := styles.Props{
		styles.BackgroundColor: "lightgrey",
		styles.Padding:         "10px",
		styles.TextAlign:       "center",
	}

	headerContent := elem.Div(attrs.Props(headerStyle),
		elem.Div(attrs.Props{attrs.Class: "flex justify-center items-center"},
			elem.A(attrs.Props{attrs.Href: "./index.html"}, elem.Text("Home page")),
			elem.Button(attrs.Props{attrs.Class: "bg-blue-500 rounded-lg w-24 h-12", attrs.ID: "click-here"}, elem.Text("Click here")),
		),
		elem.Div(nil,
			elem.Select(nil,
				elem.Option(attrs.Props{attrs.Class: "flex justify-center items-center px-4 cursor-pointer", attrs.Value: "pt-BR"},
					elem.Text("Português"),
				),

				elem.Option(attrs.Props{attrs.Class: "flex justify-center items-center px-4 cursor-pointer", attrs.Value: "en"},
					elem.Text("English"),
				),
			),
		),
	)

	footerContent := elem.Div(
		attrs.Props{attrs.Style: "", attrs.Class: "grid grid-cols-3"},
		elem.Div(nil,
			elem.Ul(nil,
				elem.Li(nil, elem.Text("Item 1")),
				elem.Li(nil, elem.Text("Item 2")),
				elem.Li(nil, elem.Text("Item 3")),
			),
		),
		elem.Div(nil,
			elem.Ul(nil,
				elem.Li(nil, elem.Text("Item 4")),
				elem.Li(nil, elem.Text("Item 5")),
				elem.Li(nil, elem.Text("Item 6")),
			),
		),
		elem.Div(nil, elem.Text("Footer third content here"),
			elem.Ul(nil,
				elem.Li(nil, elem.Text("Item 7")),
				elem.Li(nil, elem.Text("Item 8")),
				elem.Li(nil, elem.Text("Item 9")),
			),
		),
	)

	mainStyle := styles.Props{
		styles.Padding: "20px",
	}

	htmlPage := elem.Html(nil,
		elem.Head(nil,
			elem.Title(nil, elem.Text(title)),
			elem.Meta(attrs.Props{attrs.Charset: "UTF-8"}),
			// elem.Link(attrs.Props{attrs.Rel: "stylesheet", attrs.Href: "styles.css"}),
			elem.Meta(attrs.Props{attrs.Name: "viewport", attrs.Content: "width=device-width, initial-scale=1.0"}),
			elem.Script(attrs.Props{attrs.Src: "https://cdn.tailwindcss.com/"}),
			elem.Script(attrs.Props{attrs.Src: "../js/main.js"}),
		),
		elem.Body(
			attrs.Props{attrs.Style: bodyStyle.ToInline(), attrs.Class: "text-blue-800"},
			elem.Header(attrs.Props{attrs.Style: headerStyle.ToInline()},
				headerContent,
			),
			elem.Main(attrs.Props{attrs.Style: mainStyle.ToInline()},
				content,
			),
			elem.Footer(attrs.Props{attrs.Style: footerStyle.ToInline()},
				footerContent,
			),
		),
	)

	return htmlPage.Render()
}

// precisamos criar uma função para escrever o HTML gerdo na func layout
// essa função vai escrever com base no 'title'
// Para converter o tipo desejado pela lib, precisamos passar o elm.Raw ( pega HTML cru ) e retorna um elem.Node, que é o tipo
// compatível com a lib
func createHTMLPage(title, content string) string {
	htmlOutput := layout(title, elem.Raw(content))

	postFileName := title + ".html"
	filepath := filepath.Join("public", postFileName)
	os.WriteFile(filepath, []byte(htmlOutput), 0644)

	return postFileName
}

// precisamos de uma função que vai receber o Markdown e usar o goldmark para converter para HTML
func markdownHTML(content string) string {
	var buf bytes.Buffer
	md := goldmark.New()

	if err := md.Convert([]byte(content), &buf); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

// agora precisamos de uma função que vai percorrer todos nossos posts, ler o arquivo e passar o markdown para a func 'markdownToHtml'
// o retorn vai ser um array com todos posts convertidos
func readMarkdownPosts(dir string) []string {
	var posts []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			htmlContent := markdownHTML(string(content))
			title := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			postFilename := createHTMLPage(title, htmlContent)

			posts = append(posts, postFilename)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return posts
}

func createIndexPage(postFilenames []string) {
	listItems := make([]elem.Node, len(postFilenames))

	for i, filename := range postFilenames {
		link := elem.A(attrs.Props{attrs.Href: "./" + filename}, elem.Text(filename))
		listItems[i] = elem.Li(nil, link)
	}

	indexContent := elem.Ul(nil, listItems...)
	htmlOutput := layout("Home", indexContent)

	filepath := filepath.Join("public", "index.html")
	os.WriteFile(filepath, []byte(htmlOutput), 0644)
}

func createDirIfNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755) // ou 0700 se quiser que fique privado
		if err != nil {
			log.Fatal(err)
		}
	}
}
