package highlight

import (
	"slices"
	"strings"
	"unicode"
	"wicho/whim/app/consts"
)

func EditorUpdateSyntax(appData *consts.EditorConfig, appSyntax *consts.EditorSyntax, curRow *consts.EditorRow) {
	rowHl := make([]int, curRow.RenderSize)
	renderRunes := []rune(*curRow.Render)

	if appSyntax == nil {
		curRow.Highlights = rowHl
		return
	}

	keywords := appSyntax.Keywords

	singleLineCommentStart := appSyntax.SinglelineCommentStart
	mcs := appSyntax.MultilineCommentStart
	mce := appSyntax.MultilineCommentEnd

	prevSeperator := true
	inString := ""
	inComment := (curRow.Idx > 0 && appData.Row[curRow.Idx-1].HlOpenComment)

	i := 0
	for i < curRow.RenderSize {
		curRune := renderRunes[i]
		prevHl := consts.HL_NORMAL
		if i > 0 {
			prevHl = rowHl[i-1]
		} else {
			prevHl = consts.HL_NORMAL
		}

		if len(singleLineCommentStart) > 0 && inString == "" && !inComment {
			rowSubString := renderRunes[i:]
			if strings.HasPrefix(string(rowSubString), singleLineCommentStart) {
				for j := i; j < curRow.RenderSize; j++ {
					rowHl[j] = consts.HL_COMMENT
				}
				break
			}
		}

		if len(mcs) > 0 && len(mce) > 0 && inString == "" {
			rowSubString := renderRunes[i:]
			if inComment {
				rowHl[i] = consts.HL_MLCOMMENT
				rowSubString := renderRunes[i:]
				if strings.HasPrefix(string(rowSubString), mce) {
					for j := i; j < i+len(mce); j++ {
						rowHl[j] = consts.HL_COMMENT
					}
					i += len(mce)
					inComment = false
					prevSeperator = true
					continue
				} else {
					i++
					continue
				}
			} else if strings.HasPrefix(string(rowSubString), mcs) {
				for j := i; j < i+len(mcs); j++ {
					rowHl[j] = consts.HL_COMMENT
				}
				i += len(mcs)
				inComment = true
			}
		}

		if slices.Contains(appSyntax.Flags, "HL_HIGHLIGHT_STRINGS") {
			if inString != "" {
				rowHl[i] = consts.HL_STRING
				//Note: we do this for \\
				if curRune == '\\' && i+1 < curRow.RenderSize {
					rowHl[i+1] = consts.HL_STRING
					i += 2
					continue
				}
				if string(curRune) == inString {
					inString = ""
				}
				i++
				prevSeperator = true
				continue
			} else {
				if curRune == '"' || curRune == '\'' {
					inString = string(curRune)
					rowHl[i] = consts.HL_STRING
					i++
					continue
				}
			}
		}

		if slices.Contains(appSyntax.Flags, "HL_HIGHLIGHT_NUMBERS") {
			if unicode.IsDigit(curRune) && (prevSeperator || prevHl == consts.HL_NUMBER) || (curRune == '.' && prevHl == consts.HL_NUMBER) {
				rowHl[i] = consts.HL_NUMBER
				i++
				prevSeperator = false
				continue
			}
		}

		if prevSeperator {
			j := 0
			for ; j < len(keywords); j++ {
				curKeyword := keywords[j]
				isKeyword2 := strings.HasSuffix(keywords[j], "|")
				if isKeyword2 {
					curKeyword = strings.TrimSuffix(curKeyword, "|")
				}

				rowSubString := renderRunes[i:]
				afterKeywordSeperator := '\u0000'
				if i+len(curKeyword) < len(renderRunes) {
					afterKeywordSeperator = renderRunes[i+len(curKeyword)]
				}
				if strings.HasPrefix(string(rowSubString), curKeyword) && isSeperator(afterKeywordSeperator) {
					for k := i; k < i+len(curKeyword); k++ {
						if isKeyword2 {
							rowHl[k] = consts.HL_KEYWORD2
						} else {
							rowHl[k] = consts.HL_KEYWORD1
						}
					}
					i += len(curKeyword)
					break
				}
			}
			if j == len(keywords) {
				prevSeperator = false
				continue
			}
		}

		prevSeperator = isSeperator(curRune)
		i++
	}

	changed := (curRow.HlOpenComment != inComment)
	curRow.HlOpenComment = inComment
	if changed && curRow.Idx+1 < appData.NumRows {
		EditorUpdateSyntax(appData, appSyntax, appData.Row[curRow.Idx+1])
	}

	curRow.Highlights = rowHl
}

func isSeperator(aRune rune) bool {
	isRuneSpace := unicode.IsSpace(aRune)
	isEmpty := aRune == '\u0000'
	isInvalid := strings.ContainsRune(",.()+-/*=~%<>[];", aRune)
	return isRuneSpace || isEmpty || isInvalid
}

func EditorSyntaxToColor(hl int) int {
	switch hl {
	case consts.HL_NUMBER:
		return 31
	case consts.HL_MATCH:
		return 34
	case consts.HL_STRING:
		return 35
	case consts.HL_COMMENT:
		return 36
	case consts.HL_KEYWORD1:
		return 33
	case consts.HL_KEYWORD2:
		return 32
	case consts.HL_MLCOMMENT:
		return 36
	default:
		return 37
	}
}

func EditorSelectSyntaxHighlight(appData *consts.EditorConfig) {
	appData.EditorSyntax = nil
	if appData.FileName == nil {
		return
	}
	fileNameSplit := strings.Split(*appData.FileName, ".")
	if len(fileNameSplit) < 2 {
		return
	}
	extension := fileNameSplit[len(fileNameSplit)-1]

	HLDB := consts.HLDB()
	for j := 0; j < len(HLDB); j++ {
		curEditorSyntax := HLDB[j]
		for _, item := range curEditorSyntax.Filematch {
			exten := item[1:]
			if exten == extension {
				appData.EditorSyntax = &curEditorSyntax

				for filerow := 0; filerow < appData.NumRows; filerow++ {
					EditorUpdateSyntax(appData, appData.EditorSyntax, appData.Row[filerow])
				}

				return
			}
		}
	}
}
