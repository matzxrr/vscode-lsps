package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type LSPConfig struct {
	Name            string
	RepoURL         string
	RepoBuildCmds   []string
	ServerPath      string
	ServerBuildCmds []string
	ServerOutDir    string
	BinaryName      string
	ModulePath      string
	VersionKey      string
}

var lsps = map[string]LSPConfig{
	"eslint": {
		Name:            "eslint",
		RepoURL:         "https://github.com/microsoft/vscode-eslint.git",
		RepoBuildCmds:   []string{"npm install"},
		ServerPath:      "server",
		ServerBuildCmds: []string{"npm install", "npm run webpack"},
		ServerOutDir:    "out",
		BinaryName:      "eslint-lsp",
		ModulePath:      "github.com/matzxrr/vscode-lsps/eslint",
		VersionKey:      "version",
	},
}

func printAvailableLSPs() {
	fmt.Println("Available LSPs:")
	for name := range lsps {
		fmt.Println("-", name)
	}
}

func main() {

	lspName := flag.String("lsp", "", "Language server to extract (required)")
	workDir := flag.String("workdir", os.TempDir(), "Working directory")

	flag.Parse()

	if *lspName == "" {
		printAvailableLSPs()
		log.Fatal("Error: -lsp flag is required")
	}

	config, ok := lsps[*lspName]
	if !ok {
		log.Printf("Error: Unknown LSP %q", *lspName)
		printAvailableLSPs()
		os.Exit(1)
	}

	lspWorkDir := filepath.Join(*workDir, config.Name)
	if err := os.MkdirAll(lspWorkDir, 0755); err != nil {
		log.Fatalf("Error creating working directory: %v", err)
	}

	repoDir := filepath.Join(lspWorkDir, "repo")
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		log.Printf("Cloning %s", config.RepoURL)
		cmd := exec.Command("git", "clone", config.RepoURL, repoDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error cloning repository: %v", err)
		}
	} else {
		log.Printf("Updating existing repository...")
		cmd := exec.Command("git", "-C", repoDir, "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error updating repository: %v", err)
		}
	}

	for _, cmd := range config.RepoBuildCmds {
		log.Printf("Running Repo Build Command: %s", cmd)
		parts := strings.Fields(cmd)
		command := exec.Command(parts[0], parts[1:]...)
		command.Dir = repoDir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			log.Fatalf("Error running repo build command: %v", err)
		}
	}

	serverDir := filepath.Join(repoDir, config.ServerPath)
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		log.Fatalf("Server directory not found: %s", serverDir)
	}

	for _, cmd := range config.ServerBuildCmds {
		log.Printf("Running Server Build Command: %s", cmd)
		parts := strings.Fields(cmd)
		command := exec.Command(parts[0], parts[1:]...)
		command.Dir = serverDir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			log.Fatalf("Error running server build command: %v", err)
		}
	}

	outputDir := filepath.Join(lspWorkDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	serverOutDir := filepath.Join(serverDir, config.ServerOutDir)

	wildcardOut := filepath.Join(serverOutDir, "*")
	items, err := filepath.Glob(wildcardOut)
	if err != nil {
		log.Fatalf("Error finding files to copy: %v", err)
	}

	if len(items) == 0 {
		log.Fatalf("Error: No files found in '%s' to copy", wildcardOut)
	}
	
	log.Printf("Copying assets from out directory to output")
	for _, item := range items {
		command := exec.Command("cp", "-r", item, outputDir)
		command.Dir = outputDir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			log.Fatalf("Error running copy output command: %v", err)
		}
	}
}
