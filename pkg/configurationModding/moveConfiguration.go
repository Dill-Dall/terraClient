package configurationModding

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ExecuteWithCopy(srcFile string, destinationDir string) {
	Execute(srcFile, destinationDir, true)
}

func ExecuteWithMove(srcFile string, destinationDir string) {
	Execute(srcFile, destinationDir, false)
}

func Execute(configurationSrc string, destinationDir string, moveFile bool) {
	var destinationFile string
	var destConfigDir string

	if strings.Contains(destinationDir, ".hcl") {
		destinationFile = destinationDir
		destConfigDir = filepath.Dir(destinationFile)
	} else {
		if !strings.Contains(destinationDir, "configurations") {
			destConfigDir = filepath.Join(destinationDir, "configurations")
		}
		destinationFile = filepath.Join(destConfigDir, filepath.Base(configurationSrc))
	}

	os.MkdirAll(destConfigDir, 0755)
	parentOfSrc := filepath.Dir(configurationSrc)
	log.Println("Parent:" + parentOfSrc)

	createConfigFileInDest(moveFile, configurationSrc, destinationFile)

	copyTerragruntFile(configurationSrc, destinationFile)

}

func copyTerragruntFile(srcFile string, destinationFile string) {

	srcTerragruntFile := filepath.Join(filepath.Dir(filepath.Dir(srcFile)), "terragrunt.hcl")
	destTerragruntFile := filepath.Join(filepath.Dir(filepath.Dir(destinationFile)), "terragrunt.hcl")

	var updateImportLinesString string
	if _, err := os.Stat(destTerragruntFile); err == nil {
		fmt.Printf("File exists\n")
		updateImportLinesString = updateImportLines(destTerragruntFile)
	} else {
		fmt.Printf("File doesnt exists\n")
		updateImportLinesString = updateImportLines(srcTerragruntFile)
	}

	err := ioutil.WriteFile(destTerragruntFile, []byte(updateImportLinesString), 0755)
	simpeError(err)
}

func updateImportLines(terraGruntFile string) string {
	filesInDestConfigFolder, err := ioutil.ReadDir(filepath.Join(filepath.Dir(terraGruntFile), "configurations"))
	if err != nil {
		log.Fatalln(err)
	}

	input, err := ioutil.ReadFile(terraGruntFile)
	if err != nil {
		log.Fatalln(err)
	}

	var filesToImport string

	for _, value := range filesInDestConfigFolder {
		filesToImport += "\t\tread_terragrunt_config(\"configurations/" + value.Name() + "\").inputs,\n"
	}
	re := regexp.MustCompile(`merge\((\n.*)*\)`)

	return re.ReplaceAllString(string(input), `merge(
`+filesToImport+`	)`)

}

func createConfigFileInDest(moveFile bool, configurationSrc string, destinationFile string) {
	var err error
	if moveFile {
		srcFile, err := os.Open(configurationSrc)
		check(err)
		defer srcFile.Close()

		destFile, err := os.Create(destinationFile)
		simpeError(err)

		_, err = io.Copy(destFile, srcFile)
		simpeError(err)

	} else {
		err = os.Rename(configurationSrc, destinationFile)
		simpeError(err)
	}
}

func simpeError(err error) {
	if err != nil {
		log.Println(err)
	}
}
func check(err error) {
	if err != nil {
		fmt.Printf("Error : %s", err.Error())
		os.Exit(1)
	}
}
