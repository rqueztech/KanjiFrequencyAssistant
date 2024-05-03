package main

import (
    "fmt"
    "os"
    "encoding/csv"
)

func ReadCSV () ([]string, error) {
    var resourcesFile string = "./resources/RussianKanjiFrequencyList.csv"

    // Open the file
    file, err := os.Open(resourcesFile)

    // Error in the code
    if err != nil {
        fmt.Println("File not found...")
        return nil, err
    }

    // create a defer to close file until the function is finished
    defer file.Close()
    
    var lines []string

    // create a reader
    reader := csv.NewReader(file)

    // Read all records out of the reader
    records, err := reader.ReadAll()

    // do error check
    if err != nil {
        fmt.Println(err)
    } 

    for _, record := range records {
        lines.append(lines, record)
    }

    return lines, nil
}
