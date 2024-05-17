// import the main package
package main

import (
    "bufio" // import bufio to scan in user input
    "fmt" // import the fmt package for printing
    "net/url"
    "os"
    "os/exec"
    "regexp"
    "runtime"
    "strings"
    "sync"
)

// Read the from the os/exec package and create custom clearscreen
func clearScreen() {
    // Create a pointer to the exec.Cmd struct, which is used to build the command to be used.
    var cmd *exec.Cmd

    // Check the runtime of the OS and create the command accordingly
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }

    // Set the Stdout to the os.Stdout
    cmd.Stdout = os.Stdout

    // Run the command that was created
    cmd.Run()
}

type KanjiReadings struct {
    readings map[string][]rune
    strings.Builder
    regex *regexp.Regexp
}

func (kanjiOps* KanjiReadings) printMap(title string, map_result []rune, userInput string, meaning_map *map[string][]rune, readings bool) {
    // Print out the name of the function

    // Jisho link string baseline
    jishoBaseLink := "https://www.jisho.org/search/"
    
    if map_result != nil {
        fmt.Printf("============ %s ================", title)

        // Check if the value exists in the map, if not prints out DOES NOT EXIST
        for _, currentKanjiRune := range(map_result) {
            kanjiString := string(currentKanjiRune)
            escaped:= url.QueryEscape(kanjiString)
            
            kanjiOps.WriteString(jishoBaseLink)
            kanjiOps.WriteString(string(escaped))
            kanjiOps.WriteString("%20%23kanji")
            
            kanjilink := kanjiOps.String()

            kanjiOps.Reset()

            currentKanji := string(currentKanjiRune)
            meaningString := string((*meaning_map)[currentKanji])

            fmt.Printf("\n%s -> %s (%s): %s\n", kanjiString, meaningString, userInput, kanjilink)
            
            if readings == true {
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23kanji")

                readings := string(kanjiOps.readings[kanjiString])
                readings = strings.ReplaceAll(readings, "\\n", "\n")
                fmt.Println(readings)
                fmt.Printf("", )
            }
        } 

        fmt.Printf("\nNumber of [[%s]] readings --> : %d\n", userInput, len(map_result))
    }
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

    // create a boolean to track readings
    var readings bool = true

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
        "./resources/all_readings_string.csv",
    }

    var kanjiReadings *KanjiReadings

    lenFiles := len(filePaths)

    wg.Add(lenFiles)

    // Iterate through all four
    for _, filePath := range filePaths {
        go func(filePath string) {
            defer wg.Done()
            csvMap, err := ReadCSV(filePath)
            
            if err != nil {
                handleError(err, "Error with file path: " + filePath)
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

                case "./resources/all_readings_string.csv":
                    // setting the csvmap directly into the KanjiReadings map
                    kanjiReadings = &KanjiReadings {
                        readings: csvMap,
                    }

            }

        }(filePath)
    }

    // Wait for all the go routines to finish, wait on all four files
    defer wg.Wait()

    // Loop to keep the program running unless the user types in "exit"
    for {
        clearScreen()
        fmt.Print("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get readings")
        fmt.Print("romaji - prints both onyomi and kunyomi")
        fmt.Println("hiragana - prints kunyomi with hiragana")
        fmt.Println("katakana - prints onyomi")
        fmt.Print("Enter Input: (type 'exit' to quit, 'readings' to print out the readings as well)")
        scanner.Scan()
        userInput := scanner.Text()

        if userInput == "exit" {
            fmt.Println("Exiting the program...")
            break
        } else if userInput == "readings" {
            readings = !readings
            fmt.Println("Reading data silenced...")
            _ = bufio.NewScanner(os.Stdin)
            continue
        }
       
        // Send each string into the printMap
        if onyomiMap != nil {
            kanjiReadings.printMap("Onyomi", onyomiMap[userInput], userInput, &kanjiMeanings, readings)
        }

        if kunyomiMap != nil {
            kanjiReadings.printMap("Kunyomi", kunyomiMap[userInput], userInput, &kanjiMeanings, readings)
        }

        if kunyomiWithHiragana != nil {
            kanjiReadings.printMap("Kunyomi with Hiragana", kunyomiWithHiragana[userInput], userInput, &kanjiMeanings, readings)
        }
        
        fmt.Println("Press Enter to continue...")
        fmt.Scanln()
    }
}
