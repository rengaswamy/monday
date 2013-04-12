package monday

import (
    "unicode"
    "unicode/utf8"
)

// dateStringLayoutItem represents one word or set of delimiters between words.
// This is an abstraction level above date raw character string of date representation.
//
// Example: "1  February / 2013" -> 
//           dateStringLayoutItem { item: "1",        isWord: true,  isKnown: false }
//           dateStringLayoutItem { item: "  ",       isWord: false, isKnown: false }
//           dateStringLayoutItem { item: "February", isWord: true,  isKnown: true  }
//           dateStringLayoutItem { item: " / ",      isWord: false, isKnown: false }
//           dateStringLayoutItem { item: "2013",     isWord: true,  isKnown: false }
type dateStringLayoutItem struct {
    item string
    isWord bool  // true if this is a sequence of letters/digits (as opposed to a sequence of non-letters like delimiters)
    isDigit bool // true if this is a sequence only containing digits
    isKnown bool // true if this is a word and the word is a recongnized month, day, etc.
}

// extractLetterSequence extracts first word (sequence of letters ending with a non-letter)
// starting with the specified index and wraps it to dateStringLayoutItem according to the type
// of the word.
func extractLetterSequence(originalStr string, index int, locale Locale) (it dateStringLayoutItem) {
    letters := ""

    bytesToParse := []byte(originalStr[index:])
    runeCount := utf8.RuneCount(bytesToParse)

    var isWord bool
    var isDigit bool

    for i := 0; i < runeCount; i++ {
        rune, runeSize := utf8.DecodeRune(bytesToParse)
        bytesToParse = bytesToParse[runeSize:]

        if i == 0 {
            isWord = unicode.IsLetter(rune)
            isDigit = unicode.IsDigit(rune)
        } else {
            if (isWord && (!unicode.IsLetter(rune) && !unicode.IsDigit(rune))) ||
               (isDigit && !unicode.IsDigit(rune)) ||
               (!isWord && unicode.IsLetter(rune)) ||
               (!isDigit && unicode.IsDigit(rune)) {
                break
            }
        }

        letters += string(rune)
    }

    it.item = letters
    it.isWord = isWord
    it.isDigit = isDigit

    if isWord {
        vals, ok := knownWords[locale][letters]

        if ok {
            for _, v := range vals {
                if v == letters {
                    it.isKnown = true
                    break
                }
            }
        }
    }

    return
}

// stringToLayoutItems transforms raw date string (like "2 Mar 2012") into
// a set of dateStringLayoutItems, which are more convenient to work with
// in other analysis modules.
func stringToLayoutItems(dateStr string, locale Locale) (seqs []dateStringLayoutItem) {
    i := 0

    for i < len(dateStr) {
        seq := extractLetterSequence(dateStr, i, locale)
        i += len(seq.item)
        seqs = append(seqs, seq)
    }

    return
}

func layoutToString(li []dateStringLayoutItem) (s string) {
    for _, v := range li {
        s += v.item
    }
    return
}