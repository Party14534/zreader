# TODO

- [X] Open and read the container.xml file
    - Unmarshal the container.xml file and create structs to hold the data
- [X] Use the container.xml file to find and open the content.opf file
    - Create structs to hold unmarshalled data
- [X] Load a Page into Gio
    - [X] Make page pretty
    - [X] Add controls
- [X] Create parser for html files
    - [X] Add support for inline elements
    - [X] Add html element types to parsed html elements
        - Use material.LabelStyle struct to hold different types
    - [ ] Add styling to parsed html elements
- [X] Add support for unicode decimal code in html
- [X] Add support for images
- [X] Add darkmode and lightmode
- [X] Remove unnecessary spaces in ereader
- [X] Don't show metadata tags
- [X] Save already parsed chunks in memory
- [X] Don't render chunks with no text
- [X] Treat unknown elements as inline
- [X] Store unzipped epubs in a specific install location regardless of where the program is run
- [X] Add support for saving reading history
    - Chapter number
    - Page Number
    - Font size
    - Font
    - etc
- [ ] Update readme.md
- [ ] Add main menu ??
    - Change font
    - Load epubs
    - etc
- [ ] PDF support?!?

## Resources and Notes
- [Epub Resource](https://opensource.com/article/22/8/epub-file) 
- [Gio Resource](https://gioui.org/doc/learn)
