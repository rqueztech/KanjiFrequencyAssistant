package kanji

import (
    "sync"
    "strings"
    "regexp"
    "fmt"
    "net/url"
    "KanjiFrequencyHelper/utils"
)

// Please look into possibly using channels -> still need to learn how to do this, maybe a mutex alternative
type KanjiReadings struct {
    OnyomiMap map[string][]rune
    KunyomiMap map[string][]rune
    KunyomiWithHiragana map[string][]rune
    KanjiMeanings map[string][]rune
    Readings map[string][]rune
    FullDetailsBoth map[string][]rune
    FullDetailsKunyomi map[string][]rune
    FullDetailsOnyomi map[string][]rune

    Onyomifrequencyslice  [][]string
    Kunyomifrequencyslice [][]string
    Kunyomiwithhiraganafrequencyslice [][]string

    ShowReadings bool

    strings.Builder

    regex *regexp.Regexp
    mu sync.Mutex
}

func (kr *KanjiReadings) Lock() {
    kr.mu.Lock()
}

func (kr *KanjiReadings) Unlock() {
    kr.mu.Unlock()
}

func (kanjiOps* KanjiReadings) LoadFrequencies() {
    kanjiOps.Onyomifrequencyslice = make([][]string, 80)
    kanjiOps.Kunyomifrequencyslice = make([][]string, 30)
    kanjiOps.Kunyomiwithhiraganafrequencyslice = make([][]string, 25)

    for key, value := range kanjiOps.OnyomiMap {

        currentreadfrequency := len(value)

        if utils.GetPatternCleaning().IsCaptureRomanCharacters(key) {
            kanjiOps.Onyomifrequencyslice[currentreadfrequency] = append(kanjiOps.Onyomifrequencyslice[currentreadfrequency], key)
        }

    }
    
    for key, value := range kanjiOps.KunyomiMap {

        currentreadfrequency := len(value)

        if utils.GetPatternCleaning().IsCaptureRomanCharacters(key) {
            kanjiOps.Kunyomifrequencyslice[currentreadfrequency] = append(kanjiOps.Kunyomifrequencyslice[currentreadfrequency], key)
        }

    }
    for key, value := range kanjiOps.KunyomiWithHiragana {

        currentreadfrequency := len(value)

        if utils.GetPatternCleaning().IsCaptureRomanCharacters(key) {
            kanjiOps.Kunyomiwithhiraganafrequencyslice[currentreadfrequency] = append(kanjiOps.Kunyomiwithhiraganafrequencyslice[currentreadfrequency], key)
        }
    }
}

func (kanjiOps* KanjiReadings) PrintMap(title string, map_result []rune, userInput string) {
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
            readingString := string(kanjiOps.Readings[currentKanji])
            meaningString := string(kanjiOps.KanjiMeanings[currentKanji])
            

            readingStringsplit := strings.Split(readingString, "*")
            numberOfReadings := readingStringsplit[2]


            if kanjiOps.ShowReadings == true {
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23kanji")

                kanjilink := kanjiOps.String()
                kanjiOps.Reset()
                fmt.Printf("\n\n%s (%s) %s: %s -> %s: %s", kanjiString, userInput, numberOfReadings, meaningString, readingString, kanjilink)

                escaped = url.QueryEscape("*" + kanjiString + "*")
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23common")

                kanjiOps.Reset()

            } else {
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23kanji")

                linkOutput := kanjiOps.String()
                fmt.Printf("\n%s (%s): %s %s -> %s", kanjiString, userInput, numberOfReadings, meaningString, linkOutput)
                kanjiOps.Reset()
            }
        } 

        fmt.Printf("\nNumber of [[%s]] Readings --> : %d", userInput, len(map_result))
    }
}

func (kanjiOps* KanjiReadings) FrequencyAnalysis(userinput string) {
    utils.ClearScreen()

    if userinput == "Onyomi" {
        fmt.Println("Onyomi Frequency Report: ")
        for i := len(kanjiOps.Onyomifrequencyslice) - 1; i >= 0; i-- {
            if kanjiOps.Onyomifrequencyslice[i] != nil{
                fmt.Println(i, kanjiOps.Onyomifrequencyslice[i])
            }
        }
    } else if userinput == "Kunyomi" {
        fmt.Println("Kunyomi Frequency Report: ")
        for i := len(kanjiOps.Kunyomifrequencyslice) - 1; i >= 0; i-- {
            if kanjiOps.Kunyomifrequencyslice[i] != nil{
                fmt.Println(i, kanjiOps.Kunyomifrequencyslice[i])
            }
        }
    } else if userinput == "Kunyomiwithhiragana" {
        fmt.Println("Kunyomi with Hiragana Frequency Report: ")
        for i := len(kanjiOps.Kunyomiwithhiraganafrequencyslice) - 1; i >= 0; i-- {
            if kanjiOps.Kunyomiwithhiraganafrequencyslice[i] != nil{
                fmt.Println(i, kanjiOps.Kunyomiwithhiraganafrequencyslice[i])
            }
        }
    }   
}
