// import the main package
package main

import (
    "fmt" // import the fmt package for printing
    "bufio" // import bufio to scan in user input
    "os"
    "os/exec"
    "runtime"
    "net/url"
    "sync"
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

func printMap(title string, map_result []rune, userInput string, meaning_map *map[string][]rune) {
    // Print out the name of the function
    fmt.Printf("============ %s ================", title)
    
    // Check if the value exists in the map, if not prints out DOES NOT EXIST
    for _, currentRune := range(map_result) {
        escaped:= url.QueryEscape(string(currentRune))
        
        kanjiString := string(currentRune)
        meaningString := string((*meaning_map)[kanjiString])

        fmt.Printf("\n%s -> %s : https://www.jisho.org/search/%s%%20%%23kanji", kanjiString, meaningString, escaped)
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


// Main function
func main() {
    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)

    // Create a waitgroup
    var wg sync.WaitGroup

    // create four different maps of type map[string][]rune. these will return the map returned by ReadCSV
    var onyomiMap, kunyomiMap, kunyomiWithHiragana, kanjiMeanings map[string][]rune

    // Create  string array of filepaths, this will help us save into proper map
    filePaths := []string {
        "./resources/KanjiFrequencyListOnyomi.csv",
        "./resources/KanjiFrequencyListKunyomi.csv",
        "./resources/KunyomiWithHiragana.csv",
        "./resources/KanjiMeanings.csv",
    }

    lenFiles := len(filePaths)

    wg.Add(lenFiles)

    // Iterate through all four
    for _, filePath := range filePaths {
        go func(filePath string) {
            defer wg.Done()
            csvMap, err := ReadCSV(filePath)
            
            if err != nil {
                handleError(err, "File Wrong: " + filePath)
                return
            }
            switch filePath {
            case "./resources/KanjiFrequencyListOnyomi.csv":
                onyomiMap = csvMap

            case "./resources/KanjiFrequencyListKunyomi.csv": 
                kunyomiMap = csvMap

            case "./resources/KunyomiWithHiragana.csv":
                kunyomiWithHiragana = csvMap

            case "./resources/KanjiMeanings.csv":
                kanjiMeanings = csvMap
            }

        }(filePath)
    }

    defer wg.Wait()

    //handleError(err, "csv_as_onyomi_map")
    //handleError(err, "csv_as_kunyomi_map")
    //handleError(err, "csv_as_kunyomi_hiragana_map")
    //handleError(err, "kanjiMeanings")

    

    // Loop to keep the program running unless the user types in "exit"
    for {
        clearScreen()
        fmt.Print("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get readings")
        fmt.Print("romaji - prints both onyomi and kunyomi")
        fmt.Println("hiragana - prints kunyomi with hiragana")
        fmt.Println("katakana - prints onyomi")
        fmt.Print("Enter Input: (type 'exit' to quit)")
        scanner.Scan()
        userInput := scanner.Text()

        if userInput == "exit" {
            fmt.Println("Exiting the program...")
            break
        }
       
        // Send each string into the printMap
        if onyomiMap != nil {
            printMap("Onyomi", onyomiMap[userInput], userInput, &kanjiMeanings)
        }

        if kunyomiMap != nil {
            printMap("Kunyomi", kunyomiMap[userInput], userInput, &kanjiMeanings)
        }

        if kunyomiWithHiragana != nil {
            printMap("Kunyomi with Hiragana", kunyomiWithHiragana[userInput], userInput, &kanjiMeanings)
        }

        fmt.Println("Press enter to continue...")
        scanner.Scan()
    }
}
