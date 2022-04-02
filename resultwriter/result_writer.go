package resultwriter

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// WriteResultToHTML 結果をhtmlファイルに書き込む
func WriteResultToHTML(similarGroupList []SimilarImageGroupInfo, outFileName string) error {
	if len(similarGroupList) == 0 {
		return nil
	}
	file, err := os.Create(fmt.Sprintf("%s.html", outFileName))
	defer file.Close()
	if err != nil {
		return err
	}

	htmlTagList := []string{}

	template :=
		`<html>
	<body>
		%s
	</body>
	<style>
	body {
		margin: 30px;
	}
	.pair-files-block {
		padding: 10px;
		margin-bottom: 10px;
		box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
	}
	.colored {
		background-color: #EEEEEE;
	}
	</style>
</html>
	`
	blockTemplate :=
		`<div class='pair-files-block %s'>
		<p class='group-files-name'>%s</p>
	%s
</div>
	`
	imageTagTemplate := `<img src='%s' width='20%%' />`
	for i, pairInfo := range similarGroupList {
		imageTagList := []string{
			fmt.Sprintf(imageTagTemplate, pairInfo.AbsoluteFilePath),
		}
		imageFilePathList := []string{}
		blockClass := ""
		if i%2 == 0 {
			blockClass = "colored"
		}
		for _, similarImage := range pairInfo.SimilarImages {
			imageFilePathList = append(imageFilePathList, similarImage.AbsoluteFilePath)
			imageTagList = append(imageTagList, fmt.Sprintf(imageTagTemplate, similarImage.AbsoluteFilePath))
		}
		imageFilePathText := fmt.Sprintf("%s - [%s]", pairInfo.AbsoluteFilePath, strings.Join(imageFilePathList, ", "))
		htmlTagList = append(htmlTagList, fmt.Sprintf(blockTemplate, blockClass, imageFilePathText, strings.Join(imageTagList, "\n")))
	}
	file.WriteString(fmt.Sprintf(template, strings.Join(htmlTagList, "\n")))
	return nil
}

// WriteResultToCSV 結果をcsvファイルに書き込む
func WriteResultToCSV(similarGroupList []SimilarImageGroupInfo, outFileName string) error {
	if len(similarGroupList) == 0 {
		return nil
	}
	file, err := os.Create(fmt.Sprintf("%s.csv", outFileName))
	defer file.Close()
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)

	writer.Write([]string{
		"origin file name",
		"target file name",
		"origin absolute file path",
		"target absolute file path",
		"similarity",
	})
	for _, pairInfo := range similarGroupList {
		for _, similarImage := range pairInfo.SimilarImages {
			writer.Write([]string{
				pairInfo.OriginFile.Name(),
				similarImage.File.Name(),
				pairInfo.AbsoluteFilePath,
				similarImage.AbsoluteFilePath,
				strconv.FormatFloat(float64(similarImage.Similarity), 'f', 4, 32),
			})
		}
	}
	writer.Flush()
	return nil
}
