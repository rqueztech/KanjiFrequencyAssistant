package main

import (
    "encoding/csv"
    "fmt"
    "os"
)

func ReadCSV() ([]string, error) {
    var resourcesFile string = "./resources/KanjiFrequencyList.csv"

    // Open the file
    file, err := os.Open(resourcesFile)

    // Error handling
    if err != nil {
        fmt.Println("File not found...")
        return nil, err
    }

    // Close the file at the end of the function
    defer file.Close()

    var lines []string

    // Create a reader
    reader := csv.NewReader(file)

    // Read all records from the CSV
    records, err := reader.ReadAll()

    // Error handling
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    // Append each record to lines
    for _, record := range records {
        lines = append(lines, record...)
    }

    return lines, nil
}

