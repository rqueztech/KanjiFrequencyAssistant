// import the main package
package main

import (
    "fmt" // import the fmt package for printing
    "bufio" // import bufio to scan in user input
    "os"
)

// Main function
func main() {
    // Prompt the user for input
    fmt.Println("Spell out onyomi reading (hiragana, katakana, kanji): ")

    scanner := bufio.NewScanner(os.Stdin)


    // Read the CSV file and set to string
    csv_as_map, err := ReadCSV()

    // Check if there is an error
    if err != nil {
        fmt.Println(err)
    }
    
    for {
        scanner.Scan()
        fmt.Print("Enter Input: ")
        userInput := scanner.Text()
        

        if val, ok := csv_as_map[userInput]; ok {
            value := len(val) / 3
            fmt.Println(val)
            fmt.Printf("Number of occurences: %d\n", value)
        } else {
            fmt.Println("DOES NOT EXIST")
        }
    }
}
