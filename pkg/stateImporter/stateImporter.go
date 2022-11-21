package stateImporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

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

func Execute() {
	sourceStateDir := "/home/thomas/training/go/terragrunt/one"

	dat, err := os.ReadFile("./state.json")

	check(err)
	os.Chdir(sourceStateDir)

	var resourcesInformation Resources
	err2 := json.Unmarshal(dat, &resourcesInformation)

	check(err2)

	moduleArg := resourcesInformation.ResourcesList[0].Module + "." + resourcesInformation.ResourcesList[0].Type + "." + resourcesInformation.ResourcesList[0].Name
	idArg := resourcesInformation.ResourcesList[0].Instances[0].Attributes.Id
	moduleListItem := "module.s3_bucket[\"testing-them-buckets\"].aws_s3_bucket.bucket"

	fmt.Println(moduleListItem)
	fmt.Println(moduleArg)

	modules := []string{`module.s3_bucket["testing-them-buckets"].aws_s3_bucket.bucket`, `module.s3_bucket["testing-them-buckets2"].aws_s3_bucket.bucket`}

	var selectTions string

	for index, value := range modules {
		selectTions += strconv.Itoa(index) + ":" + value + "\n"
	}

	fmt.Print(selectTions)

	wordPromptContent := promptContent{
		"Please provide a word.",
		"Select module to import by index",
	}

	definition := promptGetInput(wordPromptContent)

	intVar, _ := strconv.Atoi(definition)
	fmt.Println("selected " + modules[intVar])

	if moduleArg == moduleListItem {
		//os.Chdir(destinationDir)
		fmt.Printf("terragrunt import %s %s", strconv.Quote(moduleArg), idArg)
		//doImport(moduleArg, idArg, moduleListItem)
		//doDeleteFromState(,moduleListItem)
	}
}
func doImport(moduleArg string, idArg string, moduleListItem string) {
	fmt.Println("Is equal")
	fmt.Printf("terragrunt import %s %s", strconv.Quote(moduleArg), idArg)
	os.Chdir("/home/thomas/training/go/terragrunt/two")

	cmd := exec.Command("terragrunt", "import", moduleListItem, idArg)

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
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
