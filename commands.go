package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (git *GitCommand) WithOptions(options Options) {
	git.options = options
}

func (git *GitCommand) Verify() error {

	out, err := exec.Command("git", "--version").CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("GIT: %s", out)
	return nil
}

// git clone [--mirror] git@github.com:whatever [WorkDir]]
func (git *GitCommand) Exec() error {

	cmdArgs := []string{"clone"}
	if git.options.GithubOptions.Mirror {
		cmdArgs = append(cmdArgs, "--mirror")
	}
	cmdArgs = append(cmdArgs, git.options.Repository)
	cmdArgs = append(cmdArgs, git.options.RepositoryPathWithWorkDir())

	logCommand("git", cmdArgs)
	return runCommand("git", cmdArgs)
}

func (tar *TarCommand) WithOptions(options Options) {
	tar.options = options
}

func (tar *TarCommand) Verify() error {

	out, err := exec.Command("tar", "--version").CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("TAR: %s", out)
	return nil
}

func (tar *TarCommand) Exec() error {

	cmdArgs := []string{"czv"}
	if len(tar.options.Workdir) != 0 {
		workDir := tar.options.Workdir
		if !strings.HasSuffix(workDir, "/") {
			workDir = workDir + "/"
		}
		cmdArgs = append(cmdArgs, []string{"-C", workDir}...)
	}
	archiveFile := tar.options.RepositoryPathWithWorkDir() + ".tar.gz"
	cmdArgs = append(cmdArgs, []string{"-f", archiveFile, tar.options.RepositoryName()}...)

	logCommand("tar", cmdArgs)
	return runCommand("tar", cmdArgs)
}

func (s3 *S3Command) WithOptions(options Options) {
	s3.options = options
}

func (s3 *S3Command) Verify() error {

	out, err := exec.Command("aws", "--version").CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("AWS: %s", out)
	return nil
}

func (s3 *S3Command) Exec() error {

	cmdArgs := []string{"s3", "cp", s3.options.RepositoryPathWithWorkDir() + ".tar.gz", s3.options.ArchiveTarget()}

	logCommand("aws", cmdArgs)
	return runCommand("aws", cmdArgs)
}

func (rm *CleanCommand) WithOptions(options Options) {
	rm.options = options
}

func (rm *CleanCommand) Verify() error {

	_, err := exec.Command("which", "rm").CombinedOutput()
	return err
}

func (rm *CleanCommand) Exec() error {

	cmdArgs1 := []string{"-rf", rm.options.RepositoryPathWithWorkDir() + ".tar.gz"}
	logCommand("rm", cmdArgs1)
	if err := runCommand("rm", cmdArgs1); err != nil {
		return err
	}

	cmdArgs2 := []string{"-rf", rm.options.RepositoryPathWithWorkDir()}
	logCommand("rm", cmdArgs2)
	return runCommand("rm", cmdArgs2)
}

func logCommand(cmd string, args []string) {
	fmt.Println(append([]string{"Run: ", cmd}, args...))
}

func runCommand(cmd string, args []string) error {
	osCmd := exec.Command(cmd, args...)
	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr
	return osCmd.Run()
}
