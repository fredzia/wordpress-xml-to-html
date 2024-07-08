### Description

This code is based on [wordpressXML2HTML](https://github.com/Petess/wordpressXML2HTML).

#### Output format

| WordPress element | HTML element |
|-------------------|--------------|
| Post title        | `<h2>`       |
| Post date         | `<b>`        |
| Post content      | `<p>`        |

### How to run

1. Place the XML file inside this project folder.

2. Run the following command in your shell from the project's root folder:

   ```sh
   go run main.go -f FILE -o OUTPUT [-n NUMBER] [-sd STARTDATE] [-ed ENDDATE]
   ```

   For example:

   ```sh
   go run main.go -f wordpress.xml -o new-file.html -n 10 -sd 2017-01-01 -ed 2017-12-31
   ```

   This script supports optional filtering by date and limiting the number of entries processed. Input and output file names are mandatory.

   Available parameters:

   - `-f`: Name of the XML file exported from WordPress.
   - `-o`: Name of the output HTML file.
   - `-n`: Number of entries to convert.
   - `-sd`: Start date for processing entries (format: `YYYY-MM-DD`).
   - `-ed`: End date for processing entries (format: `YYYY-MM-DD`).

#### Dependencies

To run this code locally, you need to [install Go](https://go.dev/dl/).

### Warnings

#### WordPress pages

Currently, there is no way to separate pages from posts, so they will be formatted as equal elements in the output file. If necessary, manually remove pages you do not wish to include.

#### Post dates

Posts without a publication date in the XML file will be ignored and excluded from the output file.

#### Media files

Media files should be exported separately and manually added to the HTML file.

### Code explanation

- **Struct Definitions (`WordpressExport`, `Channel`, `Item`)**:
  - These structs define the structure of the XML data exported from WordPress. Each `Item` represents a post with fields such as `Title`, `Link`, `PubDate`, `Description`, and `Content`.

- **Command-line Flags**:
  - The `flag` package is used to parse command-line arguments (`-f`, `-o`, `-n`, `-sd`, `-ed`) corresponding to the XML input file, output HTML file, maximum number of entries to process, start date, and end date respectively.

- **Date Parsing**:
  - The `parseDate` function converts date strings in `YYYY-MM-DD` format to `time.Time` objects. It handles errors and returns a zero `time.Time{}` object if the date cannot be parsed.

- **XML Parsing**:
  - The XML file (`fileName`) is read and parsed into a `WordpressExport` struct using `xml.Unmarshal`.

- **Processing and Output**:
  - For each post (`Item`) in `WordpressExport.Channel.Items`:
    - The `PubDate` is parsed into a `time.Time` object.
    - Entries are filtered based on the provided `startDate` and `endDate`.
    - Valid entries are formatted into HTML and written to the output file (`outFileName`).

- **HTML Generation**:
  - The HTML file begins with `<html><body>` and ends with `</body></html>`.
  - Each post's title (`<h2>`), date (`<b>`), and content (`<p>`) are formatted accordingly.

- **Error Handling**:
  - Errors related to file operations (`os.Open`, `os.Create`), XML parsing (`xml.Unmarshal`), and date parsing (`time.Parse`) are handled with error messages printed to standard output.