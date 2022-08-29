package helpers

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gopkg.in/yaml.v3"
)

func CommonDownloadTranslations(dir string, authFilePath string, sheetID string, readRange string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	sheetService, err := sheets.NewService(ctx, option.WithCredentialsFile(authFilePath))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client | err=%v", err)
	}

	comutils.EchoWithTime("Loading Sheet `%v` values...", sheetID)
	resp, err := sheetService.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet | err=%v", err)
	}
	if len(resp.Values) < 2 {
		return fmt.Errorf("no data in the translation file")
	}

	type (
		Message struct {
			Other string `json:"other" yaml:"other"`
		}
		MessageTable map[string]Message
	)

	getLangCode := func(value string) string {
		var (
			leftIdx  = strings.Index(value, "[") + 1
			rightIdx = strings.Index(value, "]")
		)
		return strings.ToLower(value[leftIdx:rightIdx])
	}

	var (
		metaColCount = 2
		headerRow    = resp.Values[0]
	)

	comutils.EchoWithTime("Parsing Sheet `%v` values...", sheetID)
	var (
		msgTableMap    = make(map[string]MessageTable, len(headerRow))
		idxLangCodeMap = make(map[int]string, len(headerRow))
	)
	for i, header := range headerRow {
		if i < metaColCount {
			continue
		}
		langCode := getLangCode(header.(string))
		idxLangCodeMap[i] = langCode
		msgTableMap[langCode] = make(MessageTable, len(resp.Values))
	}
	for _, row := range resp.Values[1:] {
		key := row[0].(string)
		for i := metaColCount; i < len(row); i++ {
			var (
				langCode = idxLangCodeMap[i]
				msgTable = msgTableMap[langCode]
				value    = row[i].(string)
			)
			msgTable[key] = Message{
				Other: strings.TrimSpace(value),
			}
		}
	}

	comutils.EchoWithTime("Writing translations files...")
	writtenCount := 0
	for langCode, msgTable := range msgTableMap {
		if len(msgTable) == 0 {
			continue
		}
		emptyCount := 0
		for _, value := range msgTable {
			if value.Other == "" {
				emptyCount++
			}
		}
		if (emptyCount/len(msgTable))*100 > 20 {
			continue
		}
		writeFile, err := os.Create(path.Join(dir, fmt.Sprintf("%s.yaml", langCode)))
		comutils.PanicOnError(err)

		fileYAML, err := yaml.Marshal(msgTable)
		comutils.PanicOnError(err)
		if _, err := writeFile.Write(fileYAML); err != nil {
			return err
		}

		comutils.EchoWithTime("- %v: %v messages.", langCode, len(msgTable))
		writtenCount++
	}
	comutils.EchoWithTime("Written %v translations files.", writtenCount)

	return nil
}
