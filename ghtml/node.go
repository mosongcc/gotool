package ghtml

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
)

// HTML 节点操作
// 做HTML解析的时使用，比用正则更稳

// Parse 解析HTML字符串
func Parse(htmlStr string) (n *html.Node, err error) {
	return html.Parse(strings.NewReader(htmlStr))
}

// GetNodeByTagKV 根据HTML标签与属性获取节点
// 如果没有找到匹配的标签，rNode==nil ，使用时注意判断
func GetNodeByTagKV(htmlNote *html.Node, tag string, kvMap map[string]string) (rNode *html.Node) {
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			if kvMap == nil { //没有属性，返回当前节点
				rNode = n
				return
			}
			var (
				l = len(kvMap)
				i = 0
			)
			for k, v := range kvMap {
				for _, attr := range n.Attr {
					if attr.Key == k && attr.Val == v {
						i += 1
						break
					}
				}
			}
			// 成功找到匹配的节点
			if l == i {
				rNode = n
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(htmlNote)
	return
}

// NodeToHTML 渲染为HTML字符串
func NodeToHTML(n *html.Node) (htmlStr string, err error) {
	b, err := NodeToBytes(n)
	if err != nil {
		return
	}
	htmlStr = string(b)
	return
}

func NodeToBytes(n *html.Node) (htmlBytes []byte, err error) {
	var buf bytes.Buffer
	err = html.Render(&buf, n)
	if err != nil {
		return
	}
	htmlBytes = buf.Bytes()
	return
}

// GetTagAttributeValue 获取标签属性值，不存在返回空字符串
func GetTagAttributeValue(n *html.Node, key string) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == key {
			return strings.TrimSpace(attr.Val)
		}
	}
	return ""
}
