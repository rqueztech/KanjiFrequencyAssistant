// import the main package
package main

import (
    "fmt" // import the fmt package for printing
    "bufio" // import bufio to scan in user input
    "os"
    "os/exec"
    "runtime"
    "net/url"
)

func clearScreen() {
    var cmd *exec.Cmd

    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }

    cmd.Stdout = os.Stdout
    cmd.Run()
}

func printMap(title string, map_result []rune, userInput string) {
    // Print out the name of the function
    fmt.Printf("============ %s ================", title)
    
    // Check if the value exists in the map, if not prints out DOES NOT EXIST
    for _, currentRune := range(map_result) {
        escaped:= url.QueryEscape(string(currentRune))
        
        kanjiString := string(currentRune)

        fmt.Printf("\n%s -> https://www.jisho.org/search/%s%%20%%23kanji", kanjiString, escaped)
    }    

    fmt.Printf("\nNumber of occurences: %d\n", len(map_result))
}

// create function to handle error
func handleError(err error, message string) {
    if err != nil {
        fmt.Println("Error: ", message)
        fmt.Println("Error: ", err)
    }
}

// Create a struct to define maps needed to import
type KanjiMaps struct {

}

// Main function
func main() {
    // Prompt the user for input
    fmt.Println("Spell out onyomi reading (hiragana, katakana, kanji): ")

    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)
    
    // Read the CSV file and set to hashmaps (Onyomi, Kunyomi, and Kunyomi with hiragana(verbs, adverbs, adjectives etc...))
    csv_as_onyomi_map, err := ReadCSV("./resources/KanjiFrequencyListOnyomi.csv")
    handleError(err, "csv_as_onyomi_map")

    csv_as_kunyomi_map, err := ReadCSV("./resources/KanjiFrequencyListKunyomi.csv")
    handleError(err, "csv_as_kunyomi_map")

    csv_as_kunyomi_hiragana_map, err := ReadCSV("./resources/KunyomiWithHiragana.csv")
    handleError(err, "csv_as_kunyomi_hiragana_map")
    
    // Loop to keep the program running unless the user types in "exit"
    for {
        clearScreen()
        fmt.Print("Enter Input: (type 'exit' to quit)")
        scanner.Scan()
        userInput := scanner.Text()

        if userInput == "exit" {
            fmt.Println("Exiting the program...")
            break
        }
        
        // Use the input put in by the user to do checks to see if the value exists
        onyomi_result := csv_as_onyomi_map[userInput]
        kunyomi_result := csv_as_kunyomi_map[userInput]
        kunyomi_hiragana_result := csv_as_kunyomi_hiragana_map[userInput]

        // Send each string into the printMap
        if onyomi_result != nil {
            printMap("Onyomi", onyomi_result, userInput)
        }

        if kunyomi_result != nil {
            printMap("Kunyomi", kunyomi_result, userInput)
        }

        if kunyomi_hiragana_result != nil {
            printMap("Kunyomi with Hiragana", kunyomi_hiragana_result, userInput)
        }

        fmt.Println("Press enter to continue...")
        scanner.Scan()
    }
}
