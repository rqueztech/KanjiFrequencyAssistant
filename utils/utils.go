package utils

import (
    "os"
    "os/exec"
    "runtime"
    "regexp"
    "sync"
)

type PatternCleaning struct {
    EnglishPattern *regexp.Regexp
    RomajiPattern  *regexp.Regexp
    IllegalCharacters *regexp.Regexp
    CaptureRomanCharacters *regexp.Regexp
    CaptureNonKanji *regexp.Regexp
    CaptureNonNumbers *regexp.Regexp
}

var (
    pc          *PatternCleaning
    pcInitOnce  sync.Once
)

// NewPatternCleaning creates a new instance of PatternCleaning
func NewPatternCleaning() *PatternCleaning {
    pcInitOnce.Do(func() {
        pc = &PatternCleaning{
            EnglishPattern: regexp.MustCompile(`^[A-Za-z ]+\-$`),
            RomajiPattern:  regexp.MustCompile(`^[A-Za-z ]+$`),
            IllegalCharacters: regexp.MustCompile(`^[^A-Za-z ぁ-んァ-ン]+$`),
            CaptureRomanCharacters: regexp.MustCompile(`[A-Za-z]`),
            CaptureNonKanji: regexp.MustCompile(`[^\p{Han}]`),
            CaptureNonNumbers: regexp.MustCompile(`[^\d]`),
        }
    })
    return pc
}

// GetPatternCleaning returns the singleton instance of PatternCleaning
func GetPatternCleaning() *PatternCleaning {
    return NewPatternCleaning()
}

func (pc *PatternCleaning) RemoveNotNumber(key string) string {
    return pc.CaptureNonNumbers.ReplaceAllString(key, "")
}

func (pc *PatternCleaning) RemoveNonKanji(key string) string {
    return pc.CaptureNonKanji.ReplaceAllString(key, "")
}

func (pc *PatternCleaning) IsCaptureRomanCharacters(key string) bool {
    return pc.CaptureRomanCharacters.MatchString(key)
}

func (pc *PatternCleaning) IsIllegalCharactersPattern(key string) bool {
    return pc.IllegalCharacters.MatchString(key)
}

func (pc *PatternCleaning) IsEnglishPattern(key string) bool {
    return pc.EnglishPattern.MatchString(key)
}

func (pc *PatternCleaning) IsRomajiPattern(key string) bool {
    return pc.RomajiPattern.MatchString(key)
}

// ClearScreen clears the screen
func ClearScreen() {
    var cmd *exec.Cmd

    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }

    cmd.Stdout = os.Stdout
    cmd.Run()
}

