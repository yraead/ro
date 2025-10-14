# CSV Plugin

The CSV plugin provides operators for reading and writing CSV (Comma-Separated Values) data using Go's `encoding/csv` package.

## Installation

```bash
go get github.com/samber/ro/plugins/encoding/csv
```

## Operators

### NewCSVReader

Creates an observable that reads CSV data from a `csv.Reader`.

```go
import (
    "encoding/csv"
    "strings"
    "github.com/samber/ro"
    rocsv "github.com/samber/ro/plugins/encoding/csv"
)

csvData := `name,age,city
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago`

reader := csv.NewReader(strings.NewReader(csvData))
observable := rocsv.NewCSVReader(reader)

subscription := observable.Subscribe(ro.PrintObserver[[]string]())
defer subscription.Unsubscribe()

// Output:
// Next: [name age city]
// Next: [Alice 30 New York]
// Next: [Bob 25 Los Angeles]
// Next: [Charlie 35 Chicago]
// Completed
```

### NewCSVWriter

Creates an operator that writes CSV data to a `csv.Writer` and returns the number of rows written.

```go
import (
    "bytes"
    "encoding/csv"
    "github.com/samber/ro"
    rocsv "github.com/samber/ro/plugins/encoding/csv"
)

var buf bytes.Buffer
writer := csv.NewWriter(&buf)

data := ro.Just(
    []string{"name", "age", "city"},
    []string{"Alice", "30", "New York"},
    []string{"Bob", "25", "Los Angeles"},
    []string{"Charlie", "35", "Chicago"},
)

observable := ro.Pipe1(
    data,
    rocsv.NewCSVWriter(writer),
)

subscription := observable.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()

// Output:
// Next: 4
// Completed
```

## Configuration Options

### Custom Delimiters

You can configure the CSV reader and writer to use custom delimiters:

```go
// Reading with custom delimiter
csvData := `name;age;city
Alice;30;New York
Bob;25;Los Angeles`

reader := csv.NewReader(strings.NewReader(csvData))
reader.Comma = ';' // Use semicolon as delimiter
observable := rocsv.NewCSVReader(reader)

// Writing with custom delimiter
var buf bytes.Buffer
writer := csv.NewWriter(&buf)
writer.Comma = ';' // Use semicolon as delimiter

data := ro.Just(
    []string{"name", "age", "city"},
    []string{"Alice", "30", "New York"},
)

observable := ro.Pipe1(
    data,
    rocsv.NewCSVWriter(writer),
)
```

### Quoted Fields

The CSV plugin automatically handles quoted fields:

```go
// Reading quoted fields
csvData := `name,age,city
"Alice Smith",30,"New York, NY"
"Bob Johnson",25,"Los Angeles, CA"`

reader := csv.NewReader(strings.NewReader(csvData))
observable := rocsv.NewCSVReader(reader)

// Writing quoted fields
var buf bytes.Buffer
writer := csv.NewWriter(&buf)

data := ro.Just(
    []string{"name", "age", "city"},
    []string{"Alice Smith", "30", "New York, NY"},
    []string{"Bob Johnson", "25", "Los Angeles, CA"},
)

observable := ro.Pipe1(
    data,
    rocsv.NewCSVWriter(writer),
)
```

## Error Handling

Both `NewCSVReader` and `NewCSVWriter` handle errors gracefully:

### Reading Errors

```go
csvData := `name,age,city
Alice,30,New York
Bob,25,"Los Angeles, CA"
Charlie,35,Chicago`

reader := csv.NewReader(strings.NewReader(csvData))
observable := rocsv.NewCSVReader(reader)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value []string) {
            // Handle successful row reading
        },
        func(err error) {
            // Handle CSV reading error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

### Writing Errors

```go
var buf bytes.Buffer
writer := csv.NewWriter(&buf)

data := ro.Just(
    []string{"name", "age", "city"},
    []string{"Alice", "30", "New York"},
    []string{"Bob", "25", "Los Angeles"},
)

observable := ro.Pipe1(
    data,
    rocsv.NewCSVWriter(writer),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value int) {
            // Handle successful write count
        },
        func(err error) {
            // Handle CSV writing error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Roundtrip Examples

Demonstrate roundtrip CSV reading and writing:

```go
inputData := `name,age,city
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago`

// Read from string
reader := csv.NewReader(strings.NewReader(inputData))
readObservable := rocsv.NewCSVReader(reader)

// Write to buffer
var buf bytes.Buffer
writer := csv.NewWriter(&buf)
writeObservable := ro.Pipe1(
    readObservable,
    rocsv.NewCSVWriter(writer),
)

subscription := writeObservable.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()

// Output:
// Next: 4
// Completed
```

## Real-world Example

Here's a practical example that processes CSV data with transformations:

```go
import (
    "bytes"
    "encoding/csv"
    "strings"
    "github.com/samber/ro"
    rocsv "github.com/samber/ro/plugins/encoding/csv"
)

// Process CSV data with transformations
pipeline := ro.Pipe4(
    // Read CSV data
    rocsv.NewCSVReader(csv.NewReader(strings.NewReader(`
name,age,city
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago
    `))),
    // Skip header row
    ro.Skip[[]string](1),
    // Transform data
    ro.Map(func(row []string) []string {
        // Convert age to integer and back to string for validation
        return []string{row[0], row[1], strings.ToUpper(row[2])}
    }),
    // Write transformed data
    rocsv.NewCSVWriter(csv.NewWriter(&bytes.Buffer{})),
)

subscription := pipeline.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses Go's standard `encoding/csv` package for all operations
- Error handling is built into both reading and writing operators
- CSV reading is streaming and memory-efficient for large files
- CSV writing automatically flushes data on completion or error
- Consider the size of your CSV data when processing large files
- Use appropriate delimiters and quoting for your data format
- The writer returns the count of successfully written rows 