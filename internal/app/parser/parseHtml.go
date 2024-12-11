package parser

const (
    P = iota + 1
    Div
    A
    B
    H1
    H2
    H3
    H4
    H5
    H6
    Li
    Ol
    Title
    Undefined
    Single
)

var HtmlTagMap map[string]int
var mapMade = false

type HTMLElement struct {
    Tag string
    TagCode int
    Content string
}

func initMap() {
    if mapMade { return }
    HtmlTagMap = map[string]int {
        "p" : P,
        "div" : Div,
        "a" : A,
        "b": B,
        "h1" : H1,
        "h2" : H2,
        "h3" : H3,
        "h4" : H4,
        "h5" : H5,
        "h6" : H6,
        "li" : Li,
        "ol" : Ol,
        "title": Title,
    }

    mapMade = true
}

func ParseHTML(html string) (elements []HTMLElement) {
    initMap()

    for i := 0; i < len(html); i++ {
        if html[i] == '<' {
            i++
            var element HTMLElement
            parseHTMLElement(&i, &html, &element, &elements)
        }
    }
    
    return elements
}

func parseHTMLElement(i *int, html *string, element *HTMLElement, elements *[]HTMLElement) {
    for j := *i; j < len(*html); j++ {
        if (*html)[j] == '/' {
            return
        } else {
            *i = j
            parseHTMLTag(i, html, element, elements)
            return
        }
    }
}

func parseHTMLTag(i *int, html *string, element *HTMLElement, elements *[]HTMLElement) {
    inQuotes := false
    for j := *i; j < len(*html); j++ {
        if (*html)[j] == ' ' || (*html)[j] == '>' {
            element.Tag = (*html)[*i:j]
            tagCode, ok := HtmlTagMap[element.Tag]
            if !ok {
                tagCode = Undefined
            }
            element.TagCode = tagCode

            for ; j < len(*html); j++ {
                // Don't read tags from single line elements
                if (*html)[j] == '"' { 
                    inQuotes = !inQuotes 
                } else if (*html)[j] == '/' && !inQuotes {
                    *i = j + 2
                    element.TagCode = Single
                    return     
                } else if (*html)[j] == '>' {
                    break
                }
            }

            *i = j + 1
            parseHTMLElementContent(i, html, element, elements)
            return
        }
    }
}

func parseHTMLElementContent(i *int, html *string, element *HTMLElement, elements *[]HTMLElement) {
    for j := *i; j < len(*html); j++ {
        // If there is an element inside the element create a new element
        if (*html)[j] == '<' && (*html)[j+1] != '/' {
            // Create a copy of the current element before sub-element
            var copyElement HTMLElement
            copyElement.Tag = element.Tag
            copyElement.TagCode = element.TagCode
            copyElement.Content = (*html)[*i:j]
            *elements = append((*elements), copyElement)
            //fmt.Printf("Copy: %v\n", copyElement.Content)

            // Parse and add sub-element to elements
            j++
            var subElement HTMLElement
            parseHTMLElement(&j, html, &subElement, elements)
            for ; j < len(*html) && subElement.TagCode != Single; j++ {
                if (*html)[j] == '>' {
                    break
                }
            }
            
            // Set i equal to j to split up the parent element's content
            *i = j + 1
        } else if (*html)[j] == '<' {
            element.Content = (*html)[*i:j]
            (*elements) = append((*elements), *element)
            *i = j + 1
            //fmt.Printf("End: %v\n", element.Content)
            return
        }
    }
}

