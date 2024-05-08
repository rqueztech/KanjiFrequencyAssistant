package main

import (
    "encoding/csv"
    "fmt"
    "os"
)

func ReadCSV() (map[string]string, error) {
    var resourcesFile string = "./resources/KanjiFrequencyList.csv"

    mymap := make(map[string]string)

    // Open the file
    file, err := os.Open(resourcesFile)

    // Error handling
    if err != nil {
        fmt.Println("File not found...")
        return nil, err
    }

    // Close the file at the end of the function
    defer file.Close()

    // Create a reader -> the reader is a reader interface, io 
    reader := csv.NewReader(file)

    // Read in all of the records -> returns as [][]string
    records, err := reader.ReadAll()

    // Error handling
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    // Append each record to lines
    for _, line := range(records) {
        mymap[line[0]] = line[1]
    }

    return mymap , nil
}

