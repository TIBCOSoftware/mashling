package field

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("expr-mapper-field")

type MappingField struct {
	HasSpecialField bool
	HasArray        bool
	Fields          []string
}

//GetAllspecialFields get all fields that conainer special fields
func GetAllspecialFields(path string) ([]string, error) {
	var re = regexp.MustCompile(`(\[\"(.*?)\"\])|(\[\'(.*?)\'\])`)
	var fields []string
	var lastIndex = 0
	matches := re.FindAllStringIndex(path, -1)
	for i, match := range matches {
		//log.Debugf("Mathing index %d", match)
		//log.Debugf("Mathing string %s", path[match[0]:match[1]])

		if i == 0 && lastIndex == 0 {
			startPart := trimDot(path[:match[0]])
			if startPart != "" {
				if strings.Index(startPart, ".") > 0 {
					fields = append(fields, strings.Split(startPart, ".")...)
				} else {
					fields = append(fields, startPart)
				}
			}
		} else if lastIndex > 0 {
			//between match string
			if match[0] > lastIndex {
				missingPart := path[lastIndex:match[0]]
				if missingPart != "" {
					missingPart = trimDot(missingPart)
					//Array index part then append to last one
					if missingPart != "" {
						if strings.Index(missingPart, ".") > 0 {
							misspartsArray := strings.Split(missingPart, ".")
							if strings.HasPrefix(missingPart, "[") {
								fields[len(fields)-1] = fields[len(fields)-1] + misspartsArray[0]
								misspartsArray = misspartsArray[1:]
							}
							fields = append(fields, misspartsArray...)
						} else {
							if strings.HasPrefix(missingPart, "[") {
								fields[len(fields)-1] = fields[len(fields)-1] + missingPart
							} else {
								fields = append(fields, missingPart)
							}
						}
					}
				}
			}
		}

		matchStr := path[match[0]+1 : match[1]-1]
		if HasQuote(matchStr) {
			fields = append(fields, RemoveQuote(matchStr))
		} else {
			fields = append(fields, matchStr)
		}
		lastIndex = match[1]

		//Last matches and append rest
		if i == len(matches)-1 {
			lastPart := trimDot(path[lastIndex:])
			if lastPart != "" {
				if lastPart != "" {
					//Array index part then append to last one
					if strings.Index(lastPart, ".") > 0 {
						lastPartArray := strings.Split(lastPart, ".")
						if strings.HasPrefix(lastPart, "[") {
							fields[len(fields)-1] = fields[len(fields)-1] + lastPartArray[0]
							lastPartArray = lastPartArray[1:]
						}
						fields = append(fields, lastPartArray...)
					} else {
						if strings.HasPrefix(lastPart, "[") {
							fields[len(fields)-1] = fields[len(fields)-1] + lastPart
						} else {
							fields = append(fields, lastPart)
						}
					}
				}
			}
		}
	}
	return fields, nil
}

func HasQuote(quoteStr string) bool {
	if strings.HasPrefix(quoteStr, `"`) && strings.HasSuffix(quoteStr, `"`) {
		return true
	}

	if strings.HasPrefix(quoteStr, `'`) && strings.HasSuffix(quoteStr, `'`) {
		return true
	}

	return false
}

func RemoveQuote(quoteStr string) string {
	if HasQuote(quoteStr) {
		if strings.HasPrefix(quoteStr, `"`) || strings.HasPrefix(quoteStr, `'`) {
			quoteStr = quoteStr[1 : len(quoteStr)-1]
		}
	}
	return quoteStr
}

func HasArray(path string) bool {
	var re = regexp.MustCompile(`\[(.*?)\]`)
	for _, match := range re.FindAllString(path, -1) {
		if match != "" && len(match) > 0 {
			nameInBracket := match[1 : len(match)-1]
			_, err := strconv.Atoi(nameInBracket)
			if err != nil {
				continue
			}
			return true
		}
	}
	return false
}

func HasSpecialFields(path string) bool {
	var re = regexp.MustCompile(`(\[\"(.*?)\"\])|(\[\'(.*?)\'\])`)
	for _, match := range re.FindAllString(path, -1) {
		if match != "" && len(match) > 0 {
			nameInBracket := match[1 : len(match)-1]
			if HasQuote(nameInBracket) {
				return true
			}
			_, err := strconv.Atoi(nameInBracket)
			if err != nil {
				return true
			}
		}
	}
	return false
}

func trimDot(str string) string {
	if str != "" {
		str = strings.TrimPrefix(str, ".")
		return strings.TrimSuffix(str, ".")
	}
	return str
}
