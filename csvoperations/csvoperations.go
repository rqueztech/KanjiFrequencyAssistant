package csvoperations

import (
    "encoding/csv"
    "fmt"
    "os"
    "unicode/utf8"
)

func ReadCSV(resourcesFile string) (map[string][]rune, error) {
    // create a map
    mymap := make(map[string][]rune)

    // Open the file
    file, err := os.Open(resourcesFile)
    if err != nil {
        fmt.Println(resourcesFile)
        return nil, fmt.Errorf("File not found: %w", err)
    }

    // Close the file at the end of the function
    defer file.Close()

    // Create a reader -> the reader is a reader interface, io 
    reader := csv.NewReader(file)

    // Read in all of the records -> returns as [][]string
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    // Append each record to lines
    for _, line := range(records) {
        if !utf8.ValidString(line[0]) || !utf8.ValidString(line[1]) {
            fmt.Println("Invalid UTF-8")
        }

        mymap[line[0]] = []rune(line[1])
    }

    return mymap , nil
}

