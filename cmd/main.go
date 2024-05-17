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
    sync.Mutex
    onyomiMap, kunyomiMap, kunyomiWithHiragana, kanjiMeanings, readings map[string][]rune

    strings.Builder
    regex *regexp.Regexp
}

func (kanjiOps* KanjiReadings) printMap(title string, map_result []rune, userInput string, readings bool) {
    // Print out the name of the function

    // Jisho link string baseline
    jishoBaseLink := "https://www.jisho.org/search/"
    
    if map_result != nil {
        fmt.Printf("============ %s ================", title)

        // Check if the value exists in the map, if not prints out DOES NOT EXIST
        for _, currentKanjiRune := range(map_result) {
            kanjiString := string(currentKanjiRune)
            escaped:= url.QueryEscape(kanjiString)

            currentKanji := string(currentKanjiRune)
            readingString := string(kanjiOps.readings[currentKanji])
            readingString = strings.ReplaceAll(readingString, "\\n", "\n")
            meaningString := string(kanjiOps.kanjiMeanings[currentKanji])

            if readings == true {
                kanjiOps.WriteString("\nLink: ")
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23kanji")

                kanjilink := kanjiOps.String()
                linkOutput := strings.ReplaceAll(kanjilink, "\\n", "\n")
                kanjiOps.Reset()
                    fmt.Printf("\n%s (%s): %s -> %s %s\n", kanjiString, userInput, meaningString, readingString, linkOutput)
                } else {
                    kanjiOps.WriteString(jishoBaseLink)
                    kanjiOps.WriteString(string(escaped))
                    kanjiOps.WriteString("%20%23kanji")

                    linkOutput := kanjiOps.String()
                    fmt.Printf("\n%s (%s): %s -> %s", kanjiString, userInput, meaningString, linkOutput)
                    kanjiOps.Reset()
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
    // create kanji ops blank pointer
    kanjiOps := &KanjiReadings{}

    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)

    // create a boolean to track readings
    var readings bool = true

    // Create a waitgroup
    var wg sync.WaitGroup

    // Create  string array of filepaths, this will help us save into proper map
    filePaths := []string {
        "./resources/KanjiFrequencyListOnyomi.csv",
        "./resources/KanjiFrequencyListKunyomi.csv",
        "./resources/KunyomiWithHiragana.csv",
        "./resources/KanjiMeanings.csv",
        "./resources/all_readings_string.csv",
    }

    lenFiles := len(filePaths)

    wg.Add(lenFiles)

    // Iterate through all four
    for _, filePath := range filePaths {
        go func(filePath string) {
            defer wg.Done()
            kanjiOps.Lock()
            defer kanjiOps.Unlock()
            csvMap, err := ReadCSV(filePath)
            
            if err != nil {
                handleError(err, "Error with file path: " + filePath)
                return
            }


            switch filePath {
                case "./resources/KanjiFrequencyListOnyomi.csv":
                    kanjiOps.onyomiMap = csvMap

                case "./resources/KanjiFrequencyListKunyomi.csv": 
                    kanjiOps.kunyomiMap = csvMap

                case "./resources/KunyomiWithHiragana.csv":
                    kanjiOps.kunyomiWithHiragana = csvMap

                case "./resources/KanjiMeanings.csv":
                    kanjiOps.kanjiMeanings = csvMap

                case "./resources/all_readings_string.csv":
                    // setting the csvmap directly into the KanjiReadings map
                    kanjiOps.readings = csvMap
            }

        }(filePath)
    }

    // Wait for all the go routines to finish, wait on all four files
    defer wg.Wait()

    // Loop to keep the program running unless the user types in "exit"
    for {
        clearScreen()
        fmt.Println("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get readings")
        fmt.Println("Enter Input: ('exit' to quit, 'readings' toggles verbosity: ")
        
        if readings == true {
            fmt.Println("Reading data enabled...")
        } else {
            fmt.Println("Reading data silenced...")
        }

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
        if kanjiOps.onyomiMap != nil {
            kanjiOps.printMap("Onyomi", kanjiOps.onyomiMap[userInput], userInput, readings)
        }

        if kanjiOps.kunyomiMap != nil {
            kanjiOps.printMap("Kunyomi", kanjiOps.kunyomiMap[userInput], userInput, readings)
        }

        if kanjiOps.kunyomiWithHiragana != nil {
            kanjiOps.printMap("Kunyomi with Hiragana", kanjiOps.kunyomiWithHiragana[userInput], userInput, readings)
        }
        
        fmt.Println("Press Enter to continue...")
        fmt.Scanln()
    }
}
