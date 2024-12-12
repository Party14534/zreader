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
    Br
    Title
    Xml
    Html
    Head
    Style
    Meta
    Link
    Body
    Abbr
    Bdi
    Bdo
    Cite
    Code
    Data
    Dfn
    Em
    I
    Kbd
    Mark
    Q
    Rp
    Rt
    Ruby
    S
    Samp
    Small
    Span
    Strong
    Sub
    Sup
    Time
    U
    Var
    Wbr
    Undefined
    Single
)

var HtmlTagMap map[string]int
var HtmlStructureTagMap map[string]int
var HtmlInlineTagMap map[string]int
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
        "h1" : H1,
        "h2" : H2,
        "h3" : H3,
        "h4" : H4,
        "h5" : H5,
        "h6" : H6,
        "li" : Li,
        "ol" : Ol,
    }
    HtmlStructureTagMap = map[string]int {
        "xml" : Xml,
        "?xml" : Xml,
        "html" : Html,
        "head" : Head,
        "style" : Style,
        "title": Title,
        "meta" : Meta,
        "link" : Link,
        "body" : Body,
    }
    HtmlInlineTagMap = map[string]int {
        "a" : A,
        "abbr" : Abbr,
        "b": B,
        "bdi" : Bdi,
        "bdo" : Bdo,
        "br" : Br,
        "cite" : Cite,
        "code" : Code,
        "data" : Data,
        "dfn" : Dfn,
        "em" : Em,
        "i" : I,
        "kbd" : Kbd,
        "mark" : Mark,
        "q" : Q,
        "rp" : Rp,
        "rt" : Rt,
        "ruby" : Ruby,
        "s" : S,
        "samp" : Samp,
        "small" : Small,
        "span" : Span,
        "strong" : Strong,
        "sub" : Sub,
        "sup" : Sup,
        "time" : Time,
        "u" : U,
        "var" : Var,
        "wbr" : Wbr,
    }

    mapMade = true
}

func ParseHTML(html string) (elements []HTMLElement) {
    initMap()

    var allElements []HTMLElement
    for i := 0; i < len(html); i++ {
        if html[i] == '<' {
            i++
            var element HTMLElement
            parseHTMLElement(&i, &html, &element, &allElements)
        }
    }

    // Remove structure defining elements
    for _, element := range allElements {
        // If it is a structure node do not include its content
        _, ok := HtmlStructureTagMap[element.Tag];

        // Remove all line breaks from content
        element.Content = removeNewLineFromContent(element)
        if !ok { 
            elements = append(elements, element)
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
                tagCode, ok = HtmlStructureTagMap[element.Tag] 
            }
            if !ok {
                tagCode, ok = HtmlInlineTagMap[element.Tag] 
            }
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
                    code, ok := HtmlTagMap[element.Tag]
                    if ok {
                        element.TagCode = code
                        element.Content = singleLineTagCodeToContent(element.TagCode)
                        // Append here to not fuck up the content parser
                        (*elements) = append((*elements), *element)
                    }

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
                    if (*html)[j+1] == '\n' { 
                        j++ 
                    }
                    break
                }
            }
            
            // Set i equal to j to split up the parent element's content
            *i = j + 1
        } else if (*html)[j] == '<' {
            element.Content = (*html)[*i:j]

            // Move to end of tag
            for ; *i < len(*html); *i++ {
                if (*html)[*i] == '>' {
                    if *i < len(*html) - 1 &&
                        (*html)[*i+1] == '\n' {
                        element.Content += " " 
                    }
                    break
                }   
            }

            // Add extra runes for inline variables
            switch element.TagCode {
                case Sup:
                    element.Content = "^" + element.Content 
                case Sub:
                    element.Content = "_" + element.Content
            }

            (*elements) = append((*elements), *element)
            *i = j + 1

            //fmt.Printf("End: %v\n", element.Content)
            return
        }
    }
}

func singleLineTagCodeToContent(code int) string {
    switch code {
        case Br:
            return "\n"
        default:
            return ""
    }
}

func removeNewLineFromContent(element HTMLElement) (newContent string) {
    // If it is a line break tag dont remove the new line
    if element.TagCode == Br { 
        return element.Content
    }

    for _, r := range element.Content {
        if r != '\n' { 
            newContent += string(r)
        } else {
            newContent += " "
        }
    }

    return newContent
}
