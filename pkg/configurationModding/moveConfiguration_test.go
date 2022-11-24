package configurationModding

import (
	//"regexp"

	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestReg(t *testing.T) {
	re := regexp.MustCompile(`\[\"(.*)\"\]`)
	selectedModule := `module.s3_bucket["testing-them-buckets"].aws_s3_bucket.bucket`
	fileName := filepath.Join("configurations", re.FindStringSubmatch(selectedModule)[1]) + ".hcl"
	t.Log("file is !", fileName)

	if fileName != filepath.Join("configurations", "testing-them-buckets")+".hcl" {
		t.Fatalf(`Wrong regexp: ` + fileName)
	}

}

func TestMoveConfigurationFile(t *testing.T) {
	log.Println(`Execute("./test/src/configurations/testing-them-buckets.hcl", "./test/dest")`)

	ExecuteWithMove("./test/src/configurations/testing-them-buckets.hcl", "./test/dest/")

	if _, err := os.Stat("./test/dest/configurations/testing-them-buckets.hcl"); err != nil {
		t.Fatalf(`File is not moved to test/dest/configurations/testing-them-buckets.hcl`)
	}

	if _, err := os.Stat("./test/src/configurations/testing-them-buckets.hcl"); err == nil {
		t.Fatalf(`File is still in src`)
	}

	ExecuteWithMove("./test/dest/configurations/testing-them-buckets.hcl", "./test/src/")

	if _, err := os.Stat("./test/src/configurations/testing-them-buckets.hcl"); err != nil {
		t.Fatalf(`File is not moved to test/src/configurations/testing-them-buckets.hcl`)
	}

	if _, err := os.Stat("./test/dest/configurations/testing-them-buckets.hcl"); err == nil {
		t.Fatalf(`File is still in src`)
	}
	os.RemoveAll("./test/dest/")
}

func TestCopyConfigurationFile(t *testing.T) {
	log.Println(`Execute("./test/src/configurations/testing-them-buckets.hcl", "./test/dest")`)

	ExecuteWithCopy("./test/src/configurations/testing-them-buckets.hcl", "./test/dest/")

	if _, err := os.Stat("./test/dest/configurations/testing-them-buckets.hcl"); err != nil {
		t.Fatalf(`File is not moved to test/dest/configurations/testing-them-buckets.hcl`)
	}

	os.RemoveAll("./test/dest/")
}
