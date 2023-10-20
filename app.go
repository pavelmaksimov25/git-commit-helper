package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// get pwd
	_ = getPwd()

	// using cli get branch name
	branch := getCurrentBranchName()

	// create a {ticket-number} from branch name
	ticketNum := "no-ticket"
	if branch != "" {
		ticketNum = getTicketNumFromBranch(branch)
	}

	// ask for commit message
	commitMessage := readCommitMessage()

	// todo :: 1. add chatgpt integration.
	// 		2. Send git diff to chatgpt to generate commit message.
	//		3. Check if git diff is too long and show message if it is.
	// 		4. User picks offered commit or writes its own.
	// 		5. Add commit message validator as pre-commit hook.

	// do a commit for all the files
	if confirmCommit(ticketNum, commitMessage) {
		commit(commitMessage)
	}
}

func getPwd() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path
}

func getCurrentBranchName() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	stdout, err := cmd.Output()

	if err != nil {
		return ""
	}

	return string(stdout)
}

func getTicketNumFromBranch(branch string) string {
	parts := strings.Split(branch, "/")
	ticketAndShortDesc := parts[len(parts)-1]
	parts = strings.Split(ticketAndShortDesc, "-")

	return strings.ToUpper(parts[0]) + "-" + parts[1]
}

func readCommitMessage() string {
	fmt.Println("Type commit short description and hit Enter.")

	reader := bufio.NewReader(os.Stdin)
	commitMessage, _ := reader.ReadString('\n')
	commitMessage = formatCommitMessage(commitMessage)

	return commitMessage
}

func formatCommitMessage(commitMessage string) string {
	strings.Trim(commitMessage, "\n, ")
	messageParts := strings.Split(commitMessage, " ")
	messageParts[0] = cases.Title(language.English).String(messageParts[0])

	return strings.Join(messageParts, " ")
}

func confirmCommit(ticketNum string, commitMessage string) bool {
	fmt.Println("Commit message is: " + ticketNum + ": " + commitMessage)
	fmt.Println("Commit? (Yes/No) [Yes] > ")
	r := bufio.NewReader(os.Stdin)
	res, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	resTrimmed := 'y' //commit by default
	if len(res) >= 2 {
		resTrimmed = int32(strings.ToLower(strings.TrimSpace(res))[0])
		println(resTrimmed)
	}

	return resTrimmed == 'y'
}

func commit(commitMessage string) {
	cmd := exec.Command("git", "commit", "-am", commitMessage, "-n")
	stdout, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(stdout))
}
