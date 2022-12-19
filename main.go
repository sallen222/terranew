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
	if isGitRepo() {
		gitignore()
	} else {
		log.Println("Not a git repository. Skipping .gitignore creation.")
	}
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
}

func gitignore () {
	_, err := os.Stat(".gitignore")
	if err != nil {
		if os.IsNotExist(err) {
			ignore, err := os.Create(".gitignore")
			if err != nil {
				log.Fatal(err)
			}
			defer ignore.Close()

		}
	}
	

	var gi [15][]byte
	gi[0]=[]byte("**/.terraform/*\n")
	gi[1]=[]byte("*.tfstate\n")
	gi[2]=[]byte("*.tfstate.*\n")
	gi[3]=[]byte("crash.log\n")
	gi[4]=[]byte("crash.*.log\n")
	gi[5]=[]byte("*.tfvars\n")
	gi[6]=[]byte("*.tfvars.json\n")
	gi[7]=[]byte("override.tf\n")
	gi[8]=[]byte("override.tf.json\n")
	gi[9]=[]byte("*_override.tf\n")
	gi[10]=[]byte("*_override.tf.json\n")
	gi[11]=[]byte(".terraformrc\n")
	gi[12]=[]byte("terraform.rc\n")
	gi[13]=[]byte(".terraform.lock.hcl\n")
	gi[14]=[]byte("*.exe\n")

	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range gi {
		_, err := f.Write(item)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer f.Close()
}