package extract

import (
	urlParsed "net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ghamster0/page-extraction/src/utils"
	"github.com/PuerkitoBio/goquery"
)

type Extractor struct {
	baseRe *regexp.Regexp
}

func NewExtractor() *Extractor {
	return &Extractor{
		baseRe: regexp.MustCompile(""),
	}
}

type Map map[string]interface{}
type Array []interface{}

func Extract(url string, html string, template *Template, args ...interface{}) interface{} {
	var baseUrl *urlParsed.URL
	if len(args) >= 1 {
		baseUrl = args[0].(*urlParsed.URL)
	} else {
		baseUrl, _ = utils.GetBaseUrl(url, html)
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return extract(doc.Find("body"), baseUrl, template)
}

func extract(elems *goquery.Selection, baseUrl *urlParsed.URL, template *Template) interface{} {
	if template.Selector != "" {
		elems = elems.Find(template.Selector)
	}
	switch template.Method {
	case "table":
		return ExtractTable(elems, baseUrl, template.ListData)
	case "object":
		return ExtractObject(elems.First(), baseUrl, template.ListData)
	case "number":
		return ExtractNumber(elems)
	case "html":
		return ExtractHtml(elems)
	case "img":
		return ExtractImg(elems, baseUrl)
	case "imgtxt":
		return ExtractImgText(elems, baseUrl)
	case "text", "nodetext":
		return ExtractTextDirty(elems, baseUrl, template.Method, template.Regexp, template.Attribute)
	case "link":
		return ExtractLinkDirty(elems, baseUrl)
	default:
		return "Unsupport extract type: " + template.Method
	}
}

func ExtractTable(elems *goquery.Selection, baseUrl *urlParsed.URL, templates []Template) []Map {
	res := make([]Map, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		item := ExtractObject(elem, baseUrl, templates)
		if len(item) > 0 {
			res = append(res, item)
		}
	})
	return res
}

func ExtractObject(elems *goquery.Selection, baseUrl *urlParsed.URL, templates []Template) Map {
	res := make(Map)
	for _, template := range templates {
		listName := template.ListName
		extracted := extract(elems, baseUrl, &template)
		if isEmpty(extracted) {
			continue
		}
		extractedType := reflect.TypeOf(extracted).Kind()

		if origin, exist := res[listName]; exist && extractedType == reflect.TypeOf(origin).Kind() {
			switch extractedType {
			case reflect.Map:
				for k, v := range extracted.(Map) {
					origin.(Map)[k] = v
				}
			case reflect.Array, reflect.Slice:
				origin = append(origin.(Array), extracted.(Array)...)
			case reflect.String:
				origin = origin.(string) + "\n" + extracted.(string)
			default:
				origin = extracted
			}
			res[listName] = origin
		} else {
			res[listName] = extracted
		}
	}
	return res
}

// 有attribute先抽，否则若是img抽src，否则抽innerText
// img的src不过正则
// nodetext只抽子元素里的文本
func ExtractTextDirty(elems *goquery.Selection, baseUrl *urlParsed.URL, textType string, params ...string) string {
	var reg *regexp.Regexp = nil
	var attr string = ""
	if len(params) >= 1 && params[0] != "" {
		if r, err := regexp.Compile(params[0]); err == nil {
			reg = r
		}
	}
	if len(params) >= 2 {
		attr = params[1]
	}
	res := make([]string, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		r := ""
		doMatchReg := reg != nil
		if attr != "" {
			r = ExtractAttr(elem, attr)
		} else if goquery.NodeName(elem) == "img" {
			r, _ = getImgLink(elem, baseUrl)
			doMatchReg = false
		} else {
			if textType == "text" {
				r = strings.TrimSpace(elem.Text())
			} else {
				// nodetext
				r = strings.Join(elem.Children().Map(func(i int, s *goquery.Selection) string {
					return strings.TrimSpace(s.Text())
				}), "")
			}
		}
		if !doMatchReg || reg.MatchString(r) {
			res = append(res, r)
		}
	})
	return strings.Join(res, "\n")
}

func ExtractAttr(elems *goquery.Selection, attr string) string {
	res := elems.Map(func(i int, elem *goquery.Selection) string {
		return strings.TrimSpace(elem.AttrOr(attr, ""))
	})
	res = utils.FilterStr(res, func(s string) bool { return s != "" })
	return strings.Join(res, "\n")
}

func ExtractText(elems *goquery.Selection, params ...string) string {
	res := elems.Map(func(i int, e *goquery.Selection) string {
		return strings.TrimSpace(e.Text())
	})
	return strings.Join(res, "\n")
}

func ExtractNumber(elems *goquery.Selection) float64 {
	var res float64 = 0
	NotNumberPrefix := regexp.MustCompile(`^[^\d]*`)
	elems.EachWithBreak(func(_ int, elem *goquery.Selection) bool {
		t := strings.TrimSpace(elem.Text())
		t = NotNumberPrefix.ReplaceAllString(t, "")
		if n, err := strconv.ParseFloat(t, 64); err == nil {
			res = n
			return false
		}
		return true
	})
	return res
}

func ExtractHtml(elems *goquery.Selection) string {
	res := make([]string, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		if h, err := goquery.OuterHtml(elem); err == nil {
			res = append(res, strings.TrimSpace(h))
		}
	})
	return strings.Join(res, "\n")
}

func ExtractImg(elems *goquery.Selection, baseUrl *urlParsed.URL) string {
	res := make([]string, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		if goquery.NodeName(elem) == "img" {
			if src, exist := elem.Attr("src"); exist {
				res = append(res, utils.AbsoluteURL(baseUrl, strings.TrimSpace(src)))
			}
		} else {
			elem.Find("img").Each(func(_ int, e *goquery.Selection) {
				if src, exist := e.Attr("src"); exist {
					res = append(res, utils.AbsoluteURL(baseUrl, strings.TrimSpace(src)))
				}
			})
		}
	})
	// return utils.UniqueStr(res)
	return strings.Join(utils.UniqueStr(res), "\n")
}

func ExtractImgText(elems *goquery.Selection, baseUrl *urlParsed.URL) string {
	res := make([]string, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		res = append(res, getImgText(elem, baseUrl))
	})
	// return utils.UniqueStr(res)
	return strings.Join(utils.UniqueStr(res), "\n")
}

func ExtractLinkDirty(elems *goquery.Selection, baseUrl *urlParsed.URL) []string {
	res := make([]string, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		r := ExtractLink(elem, baseUrl)
		if len(r) == 0 {
			for true {
				elem = elem.Parent()
				nn := goquery.NodeName(elem)
				if nn == "body" || nn == "" {
					break
				} else if nn == "a" {
					if href, vaild := getLink(elem, baseUrl); vaild {
						r = append(r, href)
					}
					break
				}
			}
		}
		res = append(res, r...)
	})
	return utils.UniqueStr(res)
}

func ExtractLink(elems *goquery.Selection, baseUrl *urlParsed.URL) []string {
	res := make([]string, 0)
	elems.Each(func(_ int, elem *goquery.Selection) {
		if goquery.NodeName(elem) == "a" {
			if href, vaild := getLink(elem, baseUrl); vaild {
				res = append(res, href)
			}
		} else {
			elem.Find("a").Each(func(_ int, e *goquery.Selection) {
				if href, vaild := getLink(e, baseUrl); vaild {
					res = append(res, href)
				}
			})
		}
	})
	return utils.UniqueStr(res)
}

func getImgText(elem *goquery.Selection, baseUrl *urlParsed.URL) string {
	switch goquery.NodeName(elem) {
	case "img":
		src, _ := getImgLink(elem, baseUrl)
		return "#img#" + src + "#img#"
	case "#text":
		return elem.Text()
	default:
		r := ""
		elem.Children().Each(func(_ int, e *goquery.Selection) {
			r += getImgText(e, baseUrl)
		})
		return r
	}
}

func getImgLink(elem *goquery.Selection, baseUrl *urlParsed.URL) (string, bool) {
	return getLink(elem, baseUrl, "src")
}

func getLink(elem *goquery.Selection, baseUrl *urlParsed.URL, params ...string) (string, bool) {
	attr := "href"
	if len(params) > 0 {
		attr = params[0]
	}
	if href, exist := elem.Attr(attr); exist {
		return utils.AbsoluteURL(baseUrl, strings.TrimSpace(href)), true
	}
	return "", false
}

func getSimpleType(val interface{}) string {
	switch reflect.TypeOf(val).Kind() {
	case reflect.String:
		return "string"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Map:
		return "map"
	default:
		return ""
	}
}

func isEmpty(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return reflect.ValueOf(val).Len() == 0
	default:
		return true
	}
}
