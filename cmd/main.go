// import the main package
package main

import (
    "fmt" // import the fmt package for printing
    "strings" // get the strings
)

// Main function
func main() {
    // Prompt the user for input
    fmt.Println("Spell out onyomi reading (hiragana, katakana, kanji): ")

    // Read the CSV file and set to string
    csv_as_string, err := ReadCSV()

    // Check if there is an error
    if err != nil {
        fmt.Println(err)
    }

    // create a csv string
    var csv_row_sliced []string

    kanjimap := make(map[string]string)

    // range through the returned string
    for _, csv_row := range(csv_as_string) {
        // create a new slice which is going to contain your split
        csv_row_sliced = strings.Split(csv_row, ",")

        kanjimap[csv_row_sliced[0]] = csv_row_sliced[1]
    }
    
    var name string

    for {
        fmt.Printf("Enter your name: ")
        fmt.Scan(&name)
        fmt.Println(kanjimap[name])
    }
}
