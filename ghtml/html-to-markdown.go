package ghtml

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// ** Html 转 Markdown

// TagFunc HTML标签对应Markdown标记的处理，可按需扩展与替换处理函数
var TagFunc = make(map[string]func(node *html.Node) string)

func init() {
	TagFunc = map[string]func(node *html.Node) string{
		"br": func(node *html.Node) string {
			return "\n\n"
		},
		"hr": func(node *html.Node) string {
			return "\n\n---\n\n"
		},
		"h1": func(node *html.Node) string {
			return fmt.Sprintf("\n# %s\n\n", NodeToMarkdown(node))
		},
		"h2": func(node *html.Node) string {
			return fmt.Sprintf("\n## %s\n\n", NodeToMarkdown(node))
		},
		"h3": func(node *html.Node) string {
			return fmt.Sprintf("\n### %s\n\n", NodeToMarkdown(node))
		},
		"h4": func(node *html.Node) string {
			return fmt.Sprintf("\n#### %s\n\n", NodeToMarkdown(node))
		},
		"h5": func(node *html.Node) string {
			return fmt.Sprintf("\n##### %s\n\n", NodeToMarkdown(node))
		},
		"h6": func(node *html.Node) string {
			return fmt.Sprintf("\n###### %s\n\n", NodeToMarkdown(node))
		},
		"div": func(node *html.Node) string {
			return fmt.Sprintf("%s    \n\n", NodeToMarkdown(node))
		},
		"p": func(node *html.Node) string {
			return fmt.Sprintf("%s    \n\n", NodeToMarkdown(node))
		},
		"strong": func(node *html.Node) string {
			return fmt.Sprintf("**%s**", NodeToMarkdown(node))
		},
		"em": func(node *html.Node) string {
			return fmt.Sprintf("*%s*", NodeToMarkdown(node))
		},
		"a": func(node *html.Node) string {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					return fmt.Sprintf("[%s](%s)", NodeToMarkdown(node), strings.TrimSpace(attr.Val))
				}
			}
			return NodeToMarkdown(node)
		},
		// img 当不存在src属性时，用空字符串
		"img": func(node *html.Node) string {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					return fmt.Sprintf("![%s](%s)", "img", strings.TrimSpace(attr.Val))
				}
			}
			return ""
		},
		"blockquote": func(node *html.Node) string {
			v := NodeToMarkdown(node)
			v = strings.ReplaceAll(v, "\n", "")         // 去掉换行符
			v = strings.ReplaceAll(v, "    ", "    \n") // 4个空格后添加换行符
			return fmt.Sprintf("\n> %s\n", v)
		},
		"code": func(node *html.Node) string {
			return fmt.Sprintf("`%s`", NodeToMarkdown(node))
		},
		"pre": func(node *html.Node) string {
			return fmt.Sprintf("\n\n```\n%s\n```\n\n", strings.TrimSpace(NodeToText(node)))
		},
		// dl 标签不做转换
		"dl": func(node *html.Node) string {
			v, _ := NodeToHTML(node)
			return "\n\n" + v + "\n\n"
		},
		// ul 标签不做转换
		"ul": func(node *html.Node) string {
			v, _ := NodeToHTML(node)
			return "\n\n" + v + "\n\n"
		},
		// ol 标签不做转换
		"ol": func(node *html.Node) string {
			v, _ := NodeToHTML(node)
			return "\n\n" + v + "\n\n"
		},
		// Table 标签，奇奇怪怪的用法太多，比如在表格里有换行、有代码等情况无法处理。
		// Markdown默认的表格功能太弱，保留标签不变
		"table": func(node *html.Node) string {
			v, _ := NodeToHTML(node)
			return "\n\n" + v + "\n\n"
		},
		// 忽略标签
		"script": func(node *html.Node) string {
			return ""
		},
		"style": func(node *html.Node) string {
			return ""
		},
		"iframe": func(node *html.Node) string {
			return ""
		},
	}
}

// 转 Markdown
func toMarkdown(node *html.Node) string {
	if node == nil {
		return ""
	}
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data)
	}
	if node.Type == html.ElementNode {
		f := TagFunc[node.Data]
		if f != nil {
			return f(node)
		}
	}
	return NodeToMarkdown(node)
}

// NodeToMarkdown HTML节点转为Markdown
func NodeToMarkdown(node *html.Node) (md string) {
	if node == nil {
		return
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		md += toMarkdown(c)
	}
	return
}

// NodeToText 读取HTML文本内容
func NodeToText(node *html.Node) (md string) {
	if node == nil {
		return
	}
	if node.Type == html.TextNode {
		md = strings.TrimSpace(node.Data)
		return
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		md += NodeToText(c)
	}
	return
}
