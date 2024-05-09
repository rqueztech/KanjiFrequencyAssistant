// import the main package
package main

import (
    "fmt" // import the fmt package for printing
    "bufio" // import bufio to scan in user input
    "os"
    "net/url"
)

func printMap(title string, datamap map[string][]rune, userInput string) {
    // Print out the name of the function
    fmt.Printf("============ %s ================", title)
    
    // Check if the value exists in the map, if not prints out DOES NOT EXIST
    if val, ok := datamap[userInput]; ok {
        for _, currentRune := range(val) {
            escaped:= url.QueryEscape(string(currentRune))
            
            fmt.Printf("\n%s: https://www.jisho.org/search/%s%%20%%23kanji", string(currentRune), escaped)
        }    

        fmt.Printf("\nNumber of occurences: %d\n", len(val))
    } else {
        fmt.Println("DOES NOT EXIST")
    }


}

// Main function
func main() {
    // Prompt the user for input
    fmt.Println("Spell out onyomi reading (hiragana, katakana, kanji): ")

    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)
    
    // Read the CSV file and set to hashmaps (Onyomi, Kunyomi, and Kunyomi with hiragana(verbs, adverbs, adjectives etc...))
    csv_as_onyomi_map, err := ReadCSV("./resources/KanjiFrequencyListOnyomi.csv")
    csv_as_kunyomi_map, err := ReadCSV("./resources/KanjiFrequencyListKunyomi.csv")
    csv_as_kunyomi_hiragana_map, err := ReadCSV("./resources/KunyomiWithHiragana.csv")

    // Check if there is an error
    if err != nil {
        fmt.Println(err)
    }

    // Loop to keep the program running unless the user types in "exit"
    for {
        fmt.Print("Enter Input: (type 'exit' to quit)")
        scanner.Scan()
        userInput := scanner.Text()

        if userInput == "exit" {
            fmt.Println("Exiting the program...")
            break
        }

        printMap("Onyomi", csv_as_onyomi_map, userInput)
        printMap("Kunyomi", csv_as_kunyomi_map, userInput)
        printMap("Kunyomi with Hiragana", csv_as_kunyomi_hiragana_map, userInput)
    }
}
