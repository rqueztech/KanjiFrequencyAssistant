// import the main package
package main

import (
    "fmt" // import the fmt package for printing
    "bufio" // import bufio to scan in user input
    "os"
    "net/url"
)

// Main function
func main() {
    // Prompt the user for input
    fmt.Println("Spell out onyomi reading (hiragana, katakana, kanji): ")

    scanner := bufio.NewScanner(os.Stdin)

    // Read the CSV file and set to hashmap
    csv_as_onyomi_map, err := ReadCSV("./resources/KanjiFrequencyList.csv")
    csv_as_kunyomi_map, err := ReadCSV("./resources/KanjiFrequencyListKunyomi.csv")

    csv_as_kunyomi_hiragana_map, err := ReadCSV("./resources/KunyomiWithHiragana.csv")

    // Check if there is an error
    if err != nil {
        fmt.Println(err)
    }

    for {
        fmt.Print("Enter Input: ")
        scanner.Scan()
        userInput := scanner.Text()

        fmt.Println("============ ONYOMI ================")
        if val, ok := csv_as_onyomi_map[userInput]; ok {
            for _, currentRune := range(val) {
                escaped:= url.QueryEscape(string(currentRune))
                
                fmt.Printf("\n%s: https://www.jisho.org/search/%s%%20%%23kanji", string(currentRune), escaped)
            }    

            fmt.Printf("\nNumber of occurences: %d\n", len(val))
        } else {
            fmt.Println("DOES NOT EXIST")
        }

        fmt.Println("============= KUNYOMI ===============")
        if val, ok := csv_as_kunyomi_map[userInput]; ok {
            for _, currentRune := range(val) {
                escaped:= url.QueryEscape(string(currentRune))
                
                fmt.Printf("\n%s: https://www.jisho.org/search/%s%%20%%23kanji", string(currentRune), escaped)
            }    

            fmt.Printf("\nNumber of occurences: %d\n", len(val))
        } else {
            fmt.Println("DOES NOT EXIST")
        }

        fmt.Println("============= KUNYOMI HIRAGANA MIX ===============")
        if val, ok := csv_as_kunyomi_hiragana_map[userInput]; ok {
            for _, currentRune := range(val) {
                escaped:= url.QueryEscape(string(currentRune))
                
                fmt.Printf("\n%s: https://www.jisho.org/search/%s%%20%%23kanji", string(currentRune), escaped)
            }    

            fmt.Printf("\nNumber of occurences: %d\n", len(val))
        } else {
            fmt.Println("DOES NOT EXIST")
        }
    }
}
