package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type WordpressExport struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Creator     string `xml:"creator"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
}

func main() {
	var fileName string
	var outFileName string
	var maxLines int
	var startDate string
	var endDate string

	flag.StringVar(&fileName, "f", "", "Name of the XML file exported from WordPress")
	flag.StringVar(&outFileName, "o", "", "Name of the file to output to")
	flag.IntVar(&maxLines, "n", -1, "Number of entries to convert")
	flag.StringVar(&startDate, "sd", "", "Date from which to start processing entries (in YYYY-MM-DD format)")
	flag.StringVar(&endDate, "ed", "", "Date to stop processing entries at (in YYYY-MM-DD format)")

	flag.Parse()

	if fileName == "" || outFileName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	startDateObj := parseDate(startDate)
	endDateObj := parseDate(endDate)

	xmlFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var export WordpressExport
	xml.Unmarshal(byteValue, &export)

	outputFile, err := os.Create(outFileName)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	outputFile.WriteString("<html><body>\n")

	count := 0
	for _, item := range export.Channel.Items {
		if maxLines == -1 || count < maxLines {
			pubDateObj, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				fmt.Printf("Error parsing date: %s\n", err)
				continue
			}

			if (startDateObj.IsZero() || pubDateObj.After(startDateObj) || pubDateObj.Equal(startDateObj)) &&
				(endDateObj.IsZero() || pubDateObj.Before(endDateObj) || pubDateObj.Equal(endDateObj)) {

				title := "<h2>" + item.Title + "</h2>\n"
				outputFile.WriteString(title)
				outputFile.WriteString("<b>" + pubDateObj.Format("2006-01-02") + "</b>\n")

				content := strings.Replace(item.Content, "\n", "<br>", -1)
				outputFile.WriteString("<p>" + content + "</p>\n")

				count++
			}
		}
	}

	outputFile.WriteString("</body></html>")
	fmt.Printf("Conversion complete. %d entries processed.\n", count)
}

func parseDate(dateString string) time.Time {
	if dateString == "" {
		return time.Time{}
	}
	dateObj, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Printf("Error parsing date: %s\n", err)
		return time.Time{}
	}
	return dateObj
}
