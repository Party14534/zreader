package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
)

const (
    EPUBMIMETYPE = "application/epub+zip"
    CONTAINERXMLLOCATION = "META-INF/container.xml"
)


func LoadEpubBook(path, dest string, e *ebooktype.EBook) error {
    // UnZip the epub
    epubDest, err := unzipEpub(path, dest)
    if err != nil {
        return err
    }

    // Validate MIME type
    err = validateMIMEType(epubDest)
    if err != nil {
        fmt.Println("Failed validation")
        return err
    }

    // Get container.xml values
    container, err := getContainerXMLFile(epubDest)
    if err != nil {
        return err
    }

    // Use container.xml to find the content.opf file
    content, err := getContentOPFFile(epubDest, container.RootFile.FullPath)
    if err != nil {
        return err
    }

    e.Dest = epubDest
    e.Type = ebooktype.EPUB

    e.Title = content.Meta.Title
    e.Creator = content.Meta.Creator
    e.Language = content.Meta.Language
    
    // Need to modify file paths if content.obf file is not in root folder
    contentFilePath := ""
    rootFilePathSlice := strings.Split(container.RootFile.FullPath, "/")
    if len(rootFilePathSlice) > 1 {
        for i, path := range rootFilePathSlice {
            if i == len(rootFilePathSlice) - 1 { break }
            contentFilePath += path + string(os.PathSeparator)
        }
    }
    e.ContentFilePath = epubDest + string(os.PathSeparator) + contentFilePath

    // Only add chapters to list when they link to html
    e.Chapters = filterLinksWithCorrectExtension(content.Links, contentFilePath)

    return nil
}

func getContentOPFFile(epubDest, contentFilePath string) (EpubContentXML, error) {
    content := EpubContentXML{}

    contentXML, err := os.ReadFile(epubDest + string(os.PathSeparator) + contentFilePath)
    if err != nil {
        return content, err
    }

    err = xml.Unmarshal(contentXML, &content)
    return content, err
}

func getContainerXMLFile(epubDest string) (EpubContainerXML, error) {
    container := EpubContainerXML{}

    containerXML, err := os.ReadFile(epubDest + string(os.PathSeparator) + CONTAINERXMLLOCATION) 
    if err != nil {
        return container, err
    }

    err = xml.Unmarshal(containerXML, &container)
    return container, err
}

func validateMIMEType (epubDest string) error {
    file, err := os.ReadFile(epubDest + string(os.PathSeparator) + "mimetype")
    if err != nil {
        return err
    }

    if compare := strings.Compare(string(file), EPUBMIMETYPE); compare != 0 {
        return fmt.Errorf("EPUB could not be verified\n")
    }

    return nil
}

func unzipEpub(path, dest string) (string, error) {
    r, err := zip.OpenReader(path)
    if err != nil {
        return "", err
    }

    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    ebookPath := dest + string(os.PathSeparator) + path
    if _, err := os.Stat(ebookPath); !os.IsNotExist(err) {
        return ebookPath, nil
    }
    os.MkdirAll(ebookPath, 0755);

    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()
        
        path := filepath.Join(ebookPath, f.Name)

        // Check for ZipSlip
        if !strings.HasPrefix(path, filepath.Clean(ebookPath) + string(os.PathSeparator)) {
            return fmt.Errorf("Illegal file path: %s", path)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), 0755)
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) 
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return "", err
        }
    }

    return ebookPath, nil
}

func filterLinksWithCorrectExtension(slice []ManifestLink, contentFilePath string) (result []string) {
    for _, link := range slice {
        if strings.Compare(path.Ext(link.Link), ".html") == 0 ||
            strings.Compare(path.Ext(link.Link), ".xhtml") == 0 {
            result = append(result, contentFilePath + link.Link)
        }
    }

    return result
}
