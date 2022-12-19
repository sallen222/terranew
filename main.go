package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func main () {
	
	if !terraformExists() {
		log.Fatal("Terraform is not installed. Please install Terraform before running this program.")
	}
	
	err := os.Mkdir("terraform", 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	
	runInit()
	makeFiles()
}

func terraformExists() bool {
	_, err := exec.LookPath("terraform")	
	return err == nil
}

func runInit () {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = "terraform"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}

	if out.String() != "" {
		return true
	}
	
	return false
}

func makeFiles() {
	main, err := os.Create("terraform/main.tf")
	if err != nil {
		log.Fatal(err)
	}
	defer main.Close()
	
	variables, err := os.Create("terraform/variables.tf")
	if err != nil {
		log.Fatal(err)
	}
	defer variables.Close()

	output, err := os.Create("terraform/output.tf")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	tfvars, err := os.Create("terraform/terraform.tfvars")
	if err != nil {
		log.Fatal(err)
	}
	defer tfvars.Close()

	if isGitRepo() {
		_, err = os.Stat(".gitignore")
		if err != nil {
			if os.IsNotExist(err) {
				ignore, err := os.Create(".gitignore")
				if err != nil {
					log.Fatal(err)
				}
				defer ignore.Close()
				ignore.WriteString("*.tfvars\n")
			} else {
				os.WriteFile(".gitignore", []byte("*.tfvars\n"), 0777)
			}
		}
	}
}