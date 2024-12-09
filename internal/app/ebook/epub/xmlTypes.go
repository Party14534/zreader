package epub

type EpubContainerXML struct {
    RootFile RootFileXML `xml:"rootfiles>rootfile"`
}

type RootFileXML struct {
    FullPath string `xml:"full-path,attr"`
}

type EpubContentXML struct {
    Meta Metadata `xml:"metadata"`
    Links []ManifestLink `xml:"manifest>item"`
}

type ManifestLink struct {
    ID string `xml:"id,attr"`
    Link string `xml:"href,attr"`
}

type Metadata struct {
    Title string `xml:"http://purl.org/dc/elements/1.1/ title"`
    Creator string `xml:"http://purl.org/dc/elements/1.1/ creator"`
    Language string `xml:"http://purl.org/dc/elements/1.1/ language"`
}

