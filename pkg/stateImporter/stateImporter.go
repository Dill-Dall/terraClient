package stateImporter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"terraClient/pkg/configurationModding"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

type Resources struct {
	ResourcesList []struct {
		Module    string `json:"module"`
		Type      string `json:"type"`
		Name      string `json:"name"`
		Instances []struct {
			Attributes struct {
				Id string `json:"id"`
			}
		} `json:"instances"`
	} `json:"resources"`
}

type promptContent struct {
	errorMsg string
	label    string
}

var cyan = color.New(color.FgCyan, color.Bold)
var green = color.New(color.FgGreen, color.Bold)

var sourceStateDir string
var destinationStateDir string

func Execute(sourceStateDirInput string, destinationStateDirInput string, updateConfigFile bool, moveConfigFile bool) {
	green.Printf("\nStarting module importing script for state in %s\nto state in %s", sourceStateDirInput, destinationStateDirInput)
	sourceStateDir = sourceStateDirInput
	destinationStateDir = destinationStateDirInput

	os.Chdir(sourceStateDir)

	selectedModule := selectModuleFromStateList()

	resourcesInformation := fetchResourcesInState()

	matchFound := doImport(resourcesInformation, selectedModule, destinationStateDir)

	if !matchFound {
		log.Fatal("No values found in state for " + selectedModule)
		return
	}

	deleteStateModule(selectedModule)
	green.Printf("%s has been successfully imported into state %s from %s\n", selectedModule, destinationStateDir, sourceStateDir)

	if updateConfigFile {

		re := regexp.MustCompile(`\[\"(.*)\"\]`)
		fileName := filepath.Join("configurations", re.FindStringSubmatch(selectedModule)[1]) + ".hcl"

		var prompt promptContent
		if moveConfigFile {
			prompt = promptContent{
				"Please provide a input.",
				"Move config file to " + fileName + ":  y\\n",
			}
		} else {
			prompt = promptContent{
				"Please provide a input.",
				"Copy config file to " + fileName + ":  y\\n",
			}
		}

		definition := promptGetInput(prompt)
		if definition != "y" {
			green.Printf("Answer: "+definition, " - Chose not to add config file to state configuration directory.")
			return
		}

		if _, err := os.Stat(fileName); err != nil {
			var pwd string
			pwd, err = os.Getwd()
			fmt.Printf("Filename: %v  | was not found from src %v \n %v\n", fileName, pwd, err)
		}

		if moveConfigFile {
			configurationModding.ExecuteWithMove(fileName, destinationStateDirInput)
			green.Printf("Moved file from " + fileName + " to " + destinationStateDir)
		} else {
			configurationModding.ExecuteWithCopy(fileName, destinationStateDirInput)
			green.Printf("Copied file from " + fileName + " to " + destinationStateDir)
		}

	}

}

func selectModuleFromStateList() string {

	color.Cyan("\nCMD: terragrunt state list")
	output, err := exec.Command("terragrunt", "state", "list").CombinedOutput()

	log.Println(string(output))
	if err != nil {
		log.Fatal(err)
	}

	var listOfModules []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		listOfModules = append(listOfModules, scanner.Text())
	}

	var selections string
	if len(listOfModules) == 0 {
		color.Red("No modules found in src folder")
		os.Exit(1)
	}

	green.Printf("\nSelect a module from state list in state dir %s\n", sourceStateDir)
	for index, value := range listOfModules {
		selections += strconv.Itoa(index) + ":" + value + "\n"
	}

	fmt.Print(selections)

	wordPromptContent := promptContent{
		"Please provide a word.",
		"Select module to import by index",
	}

	definition := promptGetInput(wordPromptContent)

	intVar, _ := strconv.Atoi(definition)
	fmt.Println("selected " + listOfModules[intVar])
	fmt.Println(listOfModules[intVar])

	return listOfModules[intVar]
}

func fetchResourcesInState() Resources {
	color.Cyan("\nCMD: terragrunt state pull")
	output, err := exec.Command("terragrunt", "state", "pull").CombinedOutput()

	if err != nil {
		log.Println(string(output))
		log.Fatal(err)
	}

	var resourcesInformation Resources
	err2 := json.Unmarshal(output, &resourcesInformation)
	check(err2)
	return resourcesInformation
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func doImport(resourcesInformation Resources, selectedModule string, destinationStateDir string) bool {
	var matchFound bool
	for _, value := range resourcesInformation.ResourcesList {
		moduleArg := value.Module + "." + value.Type + "." + value.Name
		if moduleArg == selectedModule {
			idArg := value.Instances[0].Attributes.Id

			os.Chdir(destinationStateDir)

			cyan.Printf("\nAt local state directory: %s", destinationStateDir)
			cyan.Printf("\nCMD: terragrunt import %s %s\n", strconv.Quote(moduleArg), idArg)

			output, err := exec.Command("terragrunt", "import", moduleArg, idArg).CombinedOutput()
			log.Println(string(output))

			if err != nil {
				log.Fatal(err)
			}

			matchFound = true
			break
		}
	}
	return matchFound
}

func deleteStateModule(selectedModule string) {
	err := os.Chdir(sourceStateDir)
	cyan.Printf("\nAt local state directory: %s\n", sourceStateDir)
	cyan.Printf("CMD: terragrunt state rm %s ", selectedModule)
	output, err := exec.Command("terragrunt", "state", "rm", selectedModule).CombinedOutput()
	log.Println(string(output))

	if err != nil {
		log.Fatal(err)
	}
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}
