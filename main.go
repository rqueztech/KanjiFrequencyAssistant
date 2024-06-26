// import the main package
package main

// Import all necessary packages 
import (
    "bufio" // import bufio to scan in user input
    "fmt" // import the fmt package for printing
    "os"
    "net/url"
    "KanjiFrequencyHelper/csvoperations"
    "KanjiFrequencyHelper/utils"
    "KanjiFrequencyHelper/kanji"
    "KanjiFrequencyHelper/keigo"
    "strings"
    "sync"
)

// Builder struct to be used by the program
type Builder struct {
    strings.Builder
}

const (
    reset = iota
    transative = 1 << iota
    intransative
    naadj
    iadj
    noun
    conjunction
    adverb
)

type FlagManager struct {
    flags int
}

func (fm *FlagManager) SetFlag(flagname string) error {
    switch flagname {
        case "reset":
            fm.flags = reset
        case "transative":
            fm.flags ^= transative
        case "intransative":
            fm.flags ^= intransative
        case "naadj":
            fm.flags ^= naadj
        case "iadj":
            fm.flags ^= iadj
        case "noun":
            fm.flags ^= noun
        case "conjunction":
            fm.flags ^= conjunction
        case "adverb":
            fm.flags ^= adverb
        default:
            return fmt.Errorf("Invalid flag name")
    }
    return nil
}

func (fm *FlagManager) GetFlag() int {
    return fm.flags
}

// Main function
func main() {

    // create kanji ops blank pointer
    kanjiOps := &kanji.KanjiReadings{}
    keigoOps := &keigo.KeigoReadings{
        AlreadyRead: false,
    }

    originalflag := "reset"

    fm := FlagManager{}

    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)


    // Create a waitgroup
    var wg sync.WaitGroup

    // Create  string array of filepaths, this will help us save into proper map
    filePaths := []string {
        "./resources/kanji_frequency_list_onyomi.csv",
        "./resources/kanji_frequency_list_kunyomi.csv",
        "./resources/kunyomi_with_hiragana.csv",
        "./resources/kanji_meanings.csv",
        "./resources/all_readings_string.csv",
        "./resources/keigo_mapper.csv",
        "./resources/kunyomi_by_endings.csv",
        "./resources/translator_map.csv",
        "./resources/kunyomi_transatives.csv",
        "./resources/particleset.csv",
    }

    lenFiles := len(filePaths)

    wg.Add(lenFiles)

    // Iterate through all five
    for _, filePath := range filePaths {
        go func(filePath string) {
            defer wg.Done()
            kanjiOps.Lock()
            defer kanjiOps.Unlock()
            csvMap, err := csvoperations.ReadCSV(filePath)
            
            if err != nil {
                fmt.Println("File Not Found")
                return
            }

            // Create switch statement for importing all csv files
            switch filePath {
                case "./resources/kanji_frequency_list_onyomi.csv":
                    kanjiOps.OnyomiMap = csvMap

                case "./resources/kanji_frequency_list_kunyomi.csv": 
                    kanjiOps.KunyomiMap = csvMap

                case "./resources/kunyomi_with_hiragana.csv":
                    kanjiOps.KunyomiWithHiragana = csvMap

                case "./resources/kanji_meanings.csv":
                    kanjiOps.KanjiMeanings = csvMap

                case "./resources/all_readings_string.csv":
                    // setting the csvmap directly into the kanj.KanjiReadings map
                    kanjiOps.Readings = csvMap

                case "./resources/keigo_mapper.csv":
                    keigoOps.KeigoMap = csvMap

                case "./resources/kunyomi_by_endings.csv":
                    kanjiOps.KunyomiByEndings = csvMap

                case "./resources/translator_map.csv":
                    kanjiOps.TranslatorMap = csvMap

                case "./resources/kunyomi_transatives.csv":
                    kanjiOps.KunyomiTransatives = csvMap
                
                case "./resources/particleset.csv":
                    kanjiOps.ParticleSet = csvMap
            }

        }(filePath)
    }

    // Wait for all the go routines to finish, wait on all five files
    wg.Wait()

    kanjiOps.LoadFrequencies()

    // Loop to keep the program running unless the user types in "exit"
    for {
        utils.ClearScreen() 
        fmt.Print("Select Function:\n1. Kanji Finder\n2. Keigo Finder\n3. Onyomi\n4. Kunyomi\n5. KunyomiWithHiragana\n6. Kanji Only\n7. Enter Phrase to link \n8. KunyomiByEndings\n9. Expressions\nEnter Input: ")

        scanner.Scan()
        applicationSelector := scanner.Text()
        
        userInput := ""

        
        for userInput != "exit" {
            if applicationSelector == "1" {
                utils.ClearScreen()
                // create a bool to track Readings
                

                fmt.Println("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get Readings")
                fmt.Println("Enter Input: ('exit' to quit, 'Readings' toggles verbosity: ")
                
                if kanjiOps.ShowReadings == true {
                    fmt.Println("Reading data enabled...")
                } else {
                    fmt.Println("Reading data silenced...")
                }

                scanner.Scan()
                userInput = scanner.Text()

                if userInput == "readings" {
                    kanjiOps.ShowReadings = !kanjiOps.ShowReadings
                    _ = bufio.NewScanner(os.Stdin)
                    continue
                }
                    
                // Send each string into the PrintMap
                if kanjiOps.OnyomiMap != nil && kanjiOps.OnyomiMap[userInput] != nil{
                    kanjiOps.PrintMap("Onyomi", kanjiOps.OnyomiMap[userInput], userInput)
                }

                if kanjiOps.KunyomiMap != nil && kanjiOps.KunyomiMap[userInput] != nil{
                    kanjiOps.PrintMap("Kunyomi", kanjiOps.KunyomiMap[userInput], userInput)
                }

                if kanjiOps.KunyomiWithHiragana != nil && kanjiOps.KunyomiWithHiragana[userInput] != nil{
                    kanjiOps.PrintMap("Kunyomiwithhiragana", kanjiOps.KunyomiWithHiragana[userInput], userInput)
                }

            } else if applicationSelector == "2" {
                keigoOps.ReadKeigo(userInput)
                scanner.Scan()
                userInput = scanner.Text()

                if keigoOps.KeigoMap != nil {
                    keigoOps.PrintMapKeigo(userInput)
                }
            } else if applicationSelector == "3" {
                kanjiOps.FrequencyAnalysis("Onyomi")
                userInput = "exit"
            } else if applicationSelector == "4" {
                kanjiOps.FrequencyAnalysis("Kunyomi")
                userInput = "exit"
            } else if applicationSelector == "5" {
                kanjiOps.FrequencyAnalysis("Kunyomiwithhiragana")
                userInput = "exit"
            } else if applicationSelector == "6" {
                fmt.Println("Enter Kanji Here: ")
                scanner.Scan()
                userInput = scanner.Text()

                charMap := make(map[rune]bool)

                userInput = utils.GetPatternCleaning().RemoveNonKanji(userInput)
                userInputFilter := ""

                for _, kanjitoclean := range userInput {
                    if kanjiOps.KanjiMeanings[string(kanjitoclean)] != nil {
                        userInputFilter += string(kanjitoclean)
                    }
                }

                for _, char := range userInputFilter {
                    charMap[char] = true
                }

                var removeduplicates string
                for char := range charMap {
                    removeduplicates += string(char)
                }

                for _, char := range removeduplicates {
                    currentkanji := string(char)
                    QueryEscaped := url.QueryEscape(currentkanji)
                    kanjiDefinition := string(kanjiOps.KanjiMeanings[string(char)])

                    numReadings := string(kanjiOps.Readings[string(char)])
                    slicedreadings := strings.Split(numReadings, "*")
                    
                    fmt.Println(slicedreadings[1])

                    if slicedreadings[1] == "Both" {
                        fmt.Printf("%s -> \t%s : \nLink: https://www.jisho.org/search/%s%20%23kanji", currentkanji, kanjiDefinition, QueryEscaped)
                        fmt.Println(slicedreadings[2])
                        fmt.Println(slicedreadings[4]) 
                    } else if slicedreadings[1] == "Kunyomi" {
                        fmt.Printf("%s -> \t%s : \nLink: https://www.jisho.org/search/%s%20%23kanji", currentkanji, kanjiDefinition, QueryEscaped)
                        fmt.Println("NUMBER OF READINGS")
                        fmt.Println(slicedreadings[4])
                    } else if slicedreadings[1] == "Onyomi" {
                        fmt.Printf("%s -> \t%s : \nLink: https://www.jisho.org/search/%s%20%23kanji", currentkanji, kanjiDefinition, QueryEscaped)
                        fmt.Println("NUMBER OF READIGNS")
                        fmt.Println(slicedreadings[4])
                    }

                }
            } else if applicationSelector == "7" {
                fmt.Println("Enter Kanji Here: ")
                userInput = ""

                for userInput != "exit" {
                    fmt.Println("Enter Input: ('exit' to quit)")
                    scanner.Scan()
                    userInput = scanner.Text()
                    _ = scanner.Text()

                    if userInput == "clear" {
                        utils.ClearScreen()
                    }

                    userInput = utils.GetPatternCleaning().RemoveNonKanji(userInput)

                    seen := make(map[rune]bool)
                    unique := []rune{}

                    for _, char := range(userInput) {
                        if _, ok := seen[char]; !ok {
                            seen[char] = true
                            unique = append(unique, char)
                        }
                    }
                    
                    for _, currentkanji := range string(unique) {
                        if kanjiOps.KanjiMeanings[string(currentkanji)] != nil {

                            kanjiMeanings := kanjiOps.KanjiMeanings[string(currentkanji)]
                            kanjiReadings := kanjiOps.Readings[string(currentkanji)]
                            kanjiStrings := string(kanjiMeanings)
                            kanjireadingssplit:= strings.Split(string(kanjiReadings), "*")

                            onyomiReadings := ""
                            kunyomiReadings := ""

                            fmt.Println("-----------------------------------")
                            switch kanjireadingssplit[0] {
                            case "Both":
                                kunyomiwords := strings.Split(string(kanjireadingssplit[3]), "、")

                                for _, word := range kunyomiwords {
                                    wordsplit := strings.Split(word, "－")

                                    counter := 1
                                    if len(wordsplit) > 1 {
                                        kanjiword := string(currentkanji) + wordsplit[1]
                                        transativity := kanjiOps.KunyomiTransatives[string(kanjiword)]
                                        transativitysplit := strings.Split(string(transativity), "*")

                                        if len(wordsplit) > 1 && len(transativitysplit) > 1 {
                                            fmt.Print(kanjiword, "(", wordsplit[0], wordsplit[1], ")\n", string(transativitysplit[0]))
                                        } else if (counter %3 == 0) {
                                            fmt.Println("")
                                        }

                                        counter += 1
                                    }
                                }

                                fmt.Println("")

                                onyomiReadings = kanjireadingssplit[2]
                                kunyomiReadings = kanjireadingssplit[4]
                                fmt.Println("Onyomi: ", onyomiReadings)
                                fmt.Println("Kunyomi: ", kunyomiReadings)
                                fmt.Printf("%s (%s) -> %s -> \nLink: https://www.jisho.org/search/%s\n", string(currentkanji), kanjireadingssplit[1], kanjiStrings, url.QueryEscape(userInput))
                            case "Kunyomi":
                                kunyomiwords := strings.Split(string(kanjireadingssplit[1]), "、")

                                for _, word := range kunyomiwords {
                                    wordsplit := strings.Split(word, "－")

                                    kanjiword := string(currentkanji) + wordsplit[1]
                                    transativity := kanjiOps.KunyomiTransatives[string(kanjiword)]
                                    transativitysplit := strings.Split(string(transativity), "*")

                                    if len(wordsplit) > 1{
                                        fmt.Print(kanjiword, "(", wordsplit[0], wordsplit[1], ")", " :: ", string(transativitysplit[0]))
                                    }

                                }

                                kunyomiReadings := kanjireadingssplit[2]
                                fmt.Println(kanjireadingssplit[1])
                                fmt.Println("Kunyomi: ", kunyomiReadings)
                                fmt.Printf("%s -> %s -> \nLink: https://www.jisho.org/search/%s\n", string(currentkanji), kanjiStrings, url.QueryEscape(userInput))
                            case "Onyomi":
                                onyomiReadings := kanjireadingssplit[2]
                                fmt.Println("Onyomi: ", onyomiReadings)
                                fmt.Printf("%s(%s) -> %s -> \nLink: https://www.jisho.org/search/%s\n", string(currentkanji), kanjireadingssplit[1], kanjiStrings, url.QueryEscape(userInput))
                            default:
                                fmt.Println("Default")
                            }
                        }
                    }
                }

                userInput = "exit"
            } else if applicationSelector == "8" {
                utils.ClearScreen()

                var flagNames = map[int]string{
                    reset:      "reset",
                    transative:   "transative",
                    intransative: "intransative",
                    naadj:        "naadj",
                    iadj:         "iadj",
                    noun:         "noun",
                    conjunction:  "conjunction",
                    adverb:       "adverb",
                }

                fmt.Println("Enter Kunyomi Word Ending Here: ")
                fmt.Println(originalflag)
                scanner.Scan()
                userInput = scanner.Text()
                
                hiraganatranslation := kanjiOps.TranslatorMap[userInput]

                if utils.GetPatternCleaning().IsVerbFlagsPattern(userInput) {
                    fm.SetFlag("reset")
                    fm.SetFlag(userInput)
                    fmt.Println("Flag On: ", fm.GetFlag())
                    originalflag = flagNames[fm.GetFlag()]
                }

                if userInput == "clear" {
                    utils.ClearScreen()
                    continue
                }

                if userInput == "exit" {
                    userInput = "exit"
                    break
                }

                typecounts := map[string]int {
                    "Transative Count: ": 0,
                    "Intransative Count: ": 0,
                    "Endings Count: ": 0,
                    "iAdj Count: ": 0,
                    "naAdj Count: ": 0,
                    "Adverb Count: ": 0,
                    "Conjunction Count: ": 0,
                }



                endingscount := 0

                for _, currentkanji := range(kanjiOps.KunyomiByEndings[userInput]) {
                    jointword := url.QueryEscape(string(currentkanji) + string(hiraganatranslation))
                    jointstring := string(currentkanji) + string(hiraganatranslation)
                    transatives := kanjiOps.KunyomiTransatives[jointstring]
                    
                    parts := strings.Split(string(transatives), "*")
                    
                    if len(parts) > 2 {
                        wordtype := parts[0]
                        definition := parts[1]
                        hiraganized := parts[2]
                        

                        if flagNames[fm.GetFlag()] == "reset" {
                            fmt.Printf("%s %s (%s) -> %s |%s| https://www.jisho.org/search/%s\n", string(currentkanji), string(hiraganatranslation), hiraganized, wordtype, definition, jointword)
                        } else {

                            if flagNames[fm.GetFlag()] == wordtype {
                                originalflag = wordtype
                                fmt.Printf("%s %s (%s) -> %s |%s| https://www.jisho.org/search/%s\n", string(currentkanji), string(hiraganatranslation), hiraganized, wordtype, definition, jointword)
                            }
                        }

                        typecounts[wordtype]++
                    }
                } 

                fmt.Printf("\nNumber of occurences [%s]: %s -> %s\n", userInput, endingscount)
                
                for key, value := range(typecounts) {
                    if typecounts[key] > 0 {
                        fmt.Printf("%s %s\n", key, value)
                    }
                }
            } else if applicationSelector == "9" {
                utils.ClearScreen()

                fmt.Print("Enter Particle: ")
                
                scanner.Scan()
                userInput := scanner.Text()

                particleresult, exists := kanjiOps.ParticleSet[userInput]
                if !exists {
                    fmt.Println("Particle not found")
                } else {
                    suffix := "*" + string(particleresult) + " #conj"
                    prefix := string(particleresult) + "* #conj"

                    // Escape the result
                    suffixResult := url.QueryEscape(suffix)
                    prefixResult := url.QueryEscape(prefix)

                    // Construct the full URLs with manually escaped asterisk
                    url1 := "Suffix: https://www.jisho.org/search/" + suffixResult
                    url2 := "Prefix: https://www.jisho.org/search/" + prefixResult

                    url1 = strings.ReplaceAll(url1, "+", "%20")
                    url2 = strings.ReplaceAll(url2, "+", "%20")

                    fmt.Println(url1)
                    fmt.Println(url2)

                    commonsuffix := "*" + string(particleresult) + " #exp"
                    commonprefix := string(particleresult) + "* #exp"

                    commonSuffixResult := url.QueryEscape(commonsuffix)
                    commonPrefixResult := url.QueryEscape(commonprefix)

                    commonUrl1 := "Suffix: https://www.jisho.org/search/" + commonSuffixResult
                    commonUrl2 := "Prefix: https://www.jisho.org/search/" + commonPrefixResult

                    commonUrl1 = strings.ReplaceAll(commonUrl1, "+", "%20")
                    commonUrl2 = strings.ReplaceAll(commonUrl2, "+", "%20")


                    fmt.Println(commonUrl1)
                    fmt.Println(commonUrl2)


                    particlesuffix := "*" + string(particleresult) + "*" + " #particle"

                    particleSuffixResult := url.QueryEscape(particlesuffix)

                    particleUrl1 := "Particle: https://www.jisho.org/search/" + particleSuffixResult

                    particleUrl1 = strings.ReplaceAll(particleUrl1, "+", "%20")

                    fmt.Println(particleUrl1)
                }

                
                
                userInput = "exit"
            } else {
                utils.ClearScreen()
                fmt.Println("Enter valid input")
                userInput = "exit"
            }

            fmt.Println("Press Enter to continue...")
            
            scanner.Scan()
            _ = scanner.Text()
        }
    }
}

