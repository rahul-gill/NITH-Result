package main

import (
	resultNITH "Result-NITH"
	"Result-NITH/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type ResultFetchProcessError int

const (
	_                                              = iota
	RollNumberDoesNotExist ResultFetchProcessError = iota + 1
	InvalidHtml
	UnknownParsingError
)

func (f ResultFetchProcessError) Error() string {
	switch {
	case errors.Is(f, RollNumberDoesNotExist):
		return fmt.Sprintf("Roll number doesn't exists")
	case errors.Is(f, InvalidHtml):
		return fmt.Sprintf("Html received is invalid")
	default:
		return fmt.Sprintf("Unknown error")
	}
}

// ParseResultHtml depends on the current html structure of official result website
// - table 0 => last update title table
// - table 1 => row 0 => rollNo, name, fathersName
// - tables from 2 to end, not including last one have semester data
// - two table for each semester
// - first table subjects
// - second table summary
// - last table useless
func ParseResultHtml(body io.ReadCloser) (user *resultNITH.StudentHtmlParsed, parseError error) {
	resultDoc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, UnknownParsingError
	}
	user = &resultNITH.StudentHtmlParsed{}
	invalidRoll := resultDoc.Find("h2").FilterFunction(func(index int, selection *goquery.Selection) bool {
		strings.EqualFold(selection.Text(), "Kindly Check the Roll Number")
		return true
	}).Length() > 0
	if invalidRoll {
		return nil, RollNumberDoesNotExist
	}

	tableFind := resultDoc.Find("table")
	semesters := (tableFind.Length() - 3) / 2
	if semesters < 0 || semesters >= tableFind.Length() {
		return nil, InvalidHtml
	}
	user.SemesterResults = make([]resultNITH.SemesterResult, semesters)
	tableFind.Each(func(tableIndex int, selection *goquery.Selection) {
		if tableIndex == 0 || tableIndex == tableFind.Length()-1 {
			//useless table
			return
		}
		if tableIndex == 1 {
			//student roll number, name, father's name
			selection.Find("td").Each(func(cellIndex int, selection *goquery.Selection) {
				txt := strings.Replace(selection.Text(), "ROLL NUMBER", "", -1)
				txt = strings.Replace(txt, "STUDENT NAME", "", -1)
				txt = strings.Replace(txt, "FATHER NAME", "", -1)
				txt = strings.TrimSpace(txt)
				switch cellIndex {
				case 0:
					user.RollNumber = txt
				case 1:
					user.Name = txt
				case 2:
					user.FathersName = txt
				}
			})
		} else if tableIndex%2 == 0 {
			//semester result table: subjects data
			rowFind := selection.Find("tr")
			subjectsResult := make([]resultNITH.SubjectResult, rowFind.Length()-2)
			rowFind.Each(func(rowIndex int, selection *goquery.Selection) {
				if rowIndex < 2 {
					return
				}
				//each row is a subject after row index 1
				selection.Find("td").Each(func(cellIndex int, selection *goquery.Selection) {
					text := strings.TrimSpace(selection.Text())
					switch cellIndex {
					case 1:
						subjectsResult[rowIndex-2].SubjectName = text
					case 2:
						subjectsResult[rowIndex-2].SubjectCode = text
					case 3:
						{
							subPoint, _ := strconv.Atoi(text)
							subjectsResult[rowIndex-2].SubPoint = int64(subPoint)
						}
					case 4:
						subjectsResult[rowIndex-2].Grade = text
					case 5:
						{
							subGP, _ := strconv.Atoi(text)
							subjectsResult[rowIndex-2].SubGP = int64(subGP)
						}
					}
				})
			})
			user.SemesterResults[(tableIndex-2)/2].SubjectResults = subjectsResult
			user.SemesterResults[(tableIndex-2)/2].SemesterNumber = int64((tableIndex-2)/2 + 1)
		} else {
			//semester result table: semester overall data
			selection.Find("tr td").Each(func(cellIndex int, selection *goquery.Selection) {
				equalCharPosition := strings.Index(selection.Text(), "=")
				text := strings.TrimSpace(selection.Text()[equalCharPosition+1:])
				if cellIndex == 1 {
					user.SemesterResults[(tableIndex-2)/2].SGPI, _ = strconv.ParseFloat(text, 64)
				} else if cellIndex == 3 {
					user.SemesterResults[(tableIndex-2)/2].CGPI, _ = strconv.ParseFloat(text, 64)
				}
			})
		}
	})
	if len(user.SemesterResults) <= 0 {
		return nil, InvalidHtml
	}
	user.CGPI = user.SemesterResults[len(user.SemesterResults)-1].CGPI
	return
}

// StoreStudentInDb this function will contain the logic to associate roll number with branch and batch year
func StoreStudentInDb(sqlDB *sql.DB, queries *db.Queries, students []resultNITH.StudentHtmlParsed) error {
	ctx := context.Background()

	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	for studentIndex, student := range students {
		//student table
		batch, branch := resultNITH.GetBatchAndBranch(student.RollNumber)
		_, err := qtx.CreateStudent(ctx, db.CreateStudentParams{
			RollNumber:     student.RollNumber,
			Name:           student.Name,
			FathersName:    student.FathersName,
			Batch:          batch,
			Branch:         branch,
			LatestSemester: int64(len(student.SemesterResults)),
			Cgpi:           student.CGPI,
		})
		if err != nil {
			println("Error when creating student entry for rollNumber ", student.RollNumber, ": ", err)
		}
		//subject table
		for _, sem := range student.SemesterResults {
			for _, subject := range sem.SubjectResults {
				_, err := qtx.CreateSubject(ctx, db.CreateSubjectParams{
					Code:    subject.SubjectCode,
					Name:    subject.SubjectName,
					Credits: subject.SubPoint,
				})
				if err != nil {
					println("Error when creating subject entry for rollNumber ", student.RollNumber, "and subject ", subject.SubjectName, ": ", err)
				}
			}
		}
		//semester result
		for semNumber, sem := range student.SemesterResults {
			_, err := qtx.CreateSemesterResultData(ctx, db.CreateSemesterResultDataParams{
				StudentRollNumber: student.RollNumber,
				Semester:          int64(semNumber + 1),
				Cgpi:              sem.CGPI,
				Sgpi:              sem.SGPI,
			})
			if err != nil {
				println("Error when creating semester entry for rollNumber ", student.RollNumber, "and semNumber ", semNumber, ": ", err)
			}
		}
		//subject result data
		for semNumber, sem := range student.SemesterResults {
			for _, subject := range sem.SubjectResults {
				_, err := qtx.CreateSubjectResultData(ctx, db.CreateSubjectResultDataParams{
					StudentRollNumber: student.RollNumber,
					SubjectCode:       subject.SubjectCode,
					Grade:             subject.Grade,
					SubGp:             subject.SubGP,
					Semester:          int64(semNumber + 1),
				})
				if err != nil {
					println("Error when creating subject result entry for rollNumber ", student.RollNumber, "and subject ", subject.SubjectName, ": ", err.Error())
				}
			}
		}
		println("Finished student Number ", studentIndex, "out of ", len(students))
	}
	return tx.Commit()
}

func getResultHtml(rollNumber string) (io.ReadCloser, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Jar: cookieJar,
	}
	path := resultNITH.GetUrlForRollNumber(rollNumber)

	//get tokens
	formPageResponse, err := httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	formPageDoc, err := goquery.NewDocumentFromReader(formPageResponse.Body)
	if err != nil {
		return nil, err
	}
	csrfToken, exists := formPageDoc.Find("[name=CSRFToken]").Attr("value")
	if !exists {
		return nil, fmt.Errorf("CSRFToken not found")
	}
	verToken, exists := formPageDoc.Find("[name=RequestVerificationToken]").Attr("value")
	if !exists {
		return nil, fmt.Errorf("RequestVerificationToken not found")
	}
	//get result html
	data := url.Values{
		"RollNumber":               {rollNumber},
		"CSRFToken":                {csrfToken},
		"RequestVerificationToken": {verToken},
		"B1":                       {"Submit"},
	}
	postReq, err := http.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	postReq.Header.Set("DNT", "1")
	postReq.Header.Set("Content-Type", " application/x-www-form-urlencoded")
	postReq.AddCookie(formPageResponse.Cookies()[0])
	resp, err := httpClient.Do(postReq)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func getResultsFromWeb(forOnlyBatch *int) []resultNITH.StudentHtmlParsed {
	//build an array of roll numbers
	rollNumbers := resultNITH.GenRollNumbers(forOnlyBatch)
	println("Total roll numbers to process: ", len(rollNumbers))
	var doneRollNumbers int32 = 0
	//build an array of student objects that contain result
	var students []resultNITH.StudentHtmlParsed

	processNext := func(rollNumber string) (*resultNITH.StudentHtmlParsed, error) {
		resultHtml, err := getResultHtml(rollNumber)
		if err != nil {
			err = fmt.Errorf("error for rollNumber %s: %w in getResultHtml", rollNumber, err)
			return nil, err
		}
		student, err := ParseResultHtml(resultHtml)
		if err == nil && student != nil {
			return student, nil
		} else {
			err = fmt.Errorf("error for rollNumber %s: %w\n", rollNumber, err)
			return nil, err
		}
	}
	const maxRetries = 5
	var retryNum = 0
	for _, rollNumber := range rollNumbers {
		retryNum = 0
	retryLoop:
		for retryNum <= maxRetries {
			var sleepDuration = retryNum + 1 + rand.Intn(2)
			log.Printf("Next fetch after %d seconds\n", sleepDuration)
			time.Sleep(time.Second * time.Duration(sleepDuration))
			student, err := processNext(rollNumber)
			if student != nil {
				students = append(students, *student)
				atomic.AddInt32(&doneRollNumbers, 1)
				log.Printf("Success for rollNumber %s; Done: %d/%d", rollNumber, doneRollNumbers, len(rollNumbers))
				break retryLoop
			} else if errors.Is(err, RollNumberDoesNotExist) {
				atomic.AddInt32(&doneRollNumbers, 1)
				log.Printf("Skipping rollNumber %s, invalid roll number; Done: %d/%d", rollNumber, doneRollNumbers, len(rollNumbers))
				break retryLoop
			}
			retryNum += 1
			log.Printf("Unknown error for roll number %s, will retry: %t", rollNumber, retryNum <= maxRetries)
		}

	}
	return students
}

func main() {
	var yearToFetch *int = nil
	if len(os.Args) > 1 {
		fmt.Printf("Got batch argument: %s\n", os.Args[1])
		yearToFetchVal, err := strconv.Atoi(os.Args[1])
		if err != nil || yearToFetchVal < 18 || yearToFetchVal > 23 {
			fmt.Errorf("Invalid argument; %w\n", err)
			os.Exit(1)
		}
		yearToFetch = &yearToFetchVal
	}
	students := getResultsFromWeb(yearToFetch)
	println("\n\nFinished fetching students\n")
	alreadyCreated := 1
	dbObj, queries := resultNITH.GetDbQueriesForNewDb("result.db", alreadyCreated == 1)

	err := StoreStudentInDb(dbObj, queries, students)
	if err != nil {
		println("Error in storeTheStudentInDb")
		println(err)
	}
}
