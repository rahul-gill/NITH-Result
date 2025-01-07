package main

import (
	resultNITH "Result-NITH"
	"Result-NITH/db"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func santizeRollNumber(rollNumberIn string, queries *db.Queries) (*string, error) {
	ctx := context.Background()
	student, err := queries.GetStudent(ctx, strings.ToUpper(rollNumberIn))
	if err != nil {
		return nil, err
	}
	return &student.RollNumber, nil
}

func santizeBranch(branchIn string) (*string, error) {
	if len(branchIn) < 3 {
		return nil, fmt.Errorf("Branch input is invalid")
	}
	code := branchIn[:3]
	branchName, found := resultNITH.BranchCodesToNames[code]
	if found {
		return &branchName, nil
	} else {
		return nil, fmt.Errorf("Branch input is invalid")
	}
}

func main() {
	if len(os.Args) <= 2 {
		fmt.Errorf("Not enough arguments provided")
		os.Exit(1)
	}
	fmt.Printf("Got roll_number argument: %s\n", os.Args[1])
	fmt.Printf("Got branch argument: %s\n", os.Args[2])

	alreadyCreated := 1
	_, queries := resultNITH.GetDbQueriesForNewDb("result.db", alreadyCreated == 1)
	rollNumber, err := santizeRollNumber(os.Args[1], queries)
	if err != nil {
		log.Fatal(err)
	}
	branch, err := santizeBranch(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	err = queries.UpdateStudentBranch(context.Background(), db.UpdateStudentBranchParams{
		Branch:     *branch,
		RollNumber: *rollNumber,
	})
	if err != nil {
		log.Fatal(err)
	}
}
