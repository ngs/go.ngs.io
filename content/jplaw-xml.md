---
title: jplaw-xml
import_path: go.ngs.io/jplaw-xml
repo_url: https://github.com/ngs/go-jplaw-xml
description: Go struct for Japanese Standard Law XML Schema (法令標準XMLスキーマ)
version: v0.0.5
documentation_url: https://pkg.go.dev/go.ngs.io/jplaw-xml
license: MIT
author: ngs
created_at: 2025-02-08T20:20:46Z
updated_at: 2025-08-12T14:37:21Z
---

# go-jplaw-xml

Go library for parsing [Japanese Standard Law XML Schema (法令標準XMLスキーマ)][xmldoc] documents.

This library provides comprehensive Go struct definitions that fully implement the Japanese Law XML Schema v3, enabling parsing and manipulation of Japanese legal documents in XML format.

## Features

- **Complete XSD Implementation**: Full coverage of the Japanese Law XML Schema v3
- **Type-Safe Parsing**: Strongly-typed Go structs for all XML elements
- **Comprehensive Coverage**: Supports all legal document structures including:
  - Law titles and metadata (Era, Year, Law numbers)
  - Table of Contents (TOC) with nested chapters and sections  
  - Main provisions with articles, paragraphs, and items
  - Supplementary provisions and appendices
  - Complex nested structures (up to 10 levels of subitems)
  - Tables, figures, and formatting elements
  - Ruby text and multilingual support

## Installation

```sh
go get go.ngs.io/jplaw-xml
```

## Usage

### Basic Parsing

```go
package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"go.ngs.io/jplaw-xml"
)

func main() {
	file, err := os.Open("path/to/law.xml")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var law jplaw.Law
	err = xml.Unmarshal(data, &law)
	if err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	// Access law metadata
	log.Println("Law Number:", law.LawNum)
	log.Println("Era:", law.Era)
	log.Println("Year:", law.Year)
	
	// Access law title
	if law.LawBody.LawTitle != nil {
		log.Println("Title:", law.LawBody.LawTitle.Content)
		log.Println("Kana:", law.LawBody.LawTitle.Kana)
	}

	// Iterate through table of contents
	if law.LawBody.TOC != nil {
		for _, chapter := range law.LawBody.TOC.TOCChapter {
			log.Printf("Chapter %s: %s", chapter.Num, chapter.ChapterTitle.Content)
		}
	}

	// Access main provisions
	for _, chapter := range law.LawBody.MainProvision.Chapter {
		log.Printf("Chapter: %s", chapter.ChapterTitle.Content)
		for _, article := range chapter.Article {
			if article.ArticleTitle != nil {
				log.Printf("  Article %s: %s", article.Num, article.ArticleTitle.Content)
			}
		}
	}
}
```

### Creating XML Documents

```go
law := jplaw.Law{
	Era:     jplaw.EraReiwa,
	Year:    5,
	Num:     "1",
	LawType: jplaw.LawTypeAct,
	Lang:    jplaw.LanguageJapanese,
	LawNum:  "令和五年法律第一号",
	LawBody: jplaw.LawBody{
		LawTitle: &jplaw.LawTitle{
			Content: "テスト法",
			Kana:    "てすとほう",
		},
		MainProvision: jplaw.MainProvision{
			Article: []jplaw.Article{
				{
					Num: "1",
					ArticleTitle: &jplaw.ArticleTitle{
						Content: "目的",
					},
					Paragraph: []jplaw.Paragraph{
						{
							Num: 1,
							ParagraphNum: jplaw.ParagraphNum{
								Content: "1",
							},
							ParagraphSentence: jplaw.ParagraphSentence{
								Sentence: []jplaw.Sentence{
									{
										Content: "この法律は、テストを目的とする。",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

// Marshal to XML
data, err := xml.MarshalIndent(law, "", "  ")
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(data))
```

## Supported Elements

The library supports all elements defined in the Japanese Law XML Schema including:

- **Law Structure**: Law, LawBody, LawTitle
- **Table of Contents**: TOC, TOCChapter, TOCSection, TOCSubsection
- **Provisions**: MainProvision, SupplProvision, Article, Paragraph
- **Items**: Item, Subitem1-10, List, Sublist1-3
- **Content**: Sentence, Column, Ruby, Line
- **Tables**: TableStruct, Table, TableRow, TableColumn
- **Figures**: FigStruct, Fig, NoteStruct, StyleStruct
- **Appendices**: AppdxTable, AppdxNote, AppdxStyle, AppdxFormat

## Constants

The library provides typed constants for common values:

```go
// Eras
jplaw.EraShowa
jplaw.EraHeisei  
jplaw.EraReiwa

// Law Types
jplaw.LawTypeAct
jplaw.LawTypeImperialOrdinance
jplaw.LawTypeCabinetOrder

// Languages
jplaw.LanguageJapanese
jplaw.LanguageEnglish

// Writing Modes
jplaw.WritingModeVertical
jplaw.WritingModeHorizontal
```

## Testing

Run the test suite:

```sh
go test -v
```

The tests include comprehensive validation against real Japanese law XML documents.

## Schema Compliance

This library implements the complete [Japanese Standard Law XML Schema v3][xmldoc] as published by the Japanese government. It has been tested against real legal documents including the Aviation Law (航空法) and provides full coverage of all schema elements and attributes.

## Author

[Atsushi Nagase]

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

[Atsushi Nagase]: https://ngs.io/
[xmldoc]: https://laws.e-gov.go.jp/docs/law-data-basic/419a603-xml-schema-for-japanese-law/

