// import the main package
package main

import (
    "fmt" // import the fmt package for printing
)

// Main function
func main() {
    // Prompt the user for input
    fmt.Println("Spell out onyomi reading (hiragana, katakana, kanji): ")

    // Read the CSV file and set to string
    csv_as_map, err := ReadCSV()

    // Check if there is an error
    if err != nil {
        fmt.Println(err)
    }
    

    for index, currentValue := range(csv_as_map) {
        fmt.Println(index, " :: ", currentValue)
    }
}
