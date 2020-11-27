// go build . && quiz -csv="problems.csv" -limit=5

package main

import (
    "flag"                   // -help flags
    "fmt"
    "os"
    "encoding/csv"
    "strings"
    "time"                  // type timer implemented in this package
    )

func main()  {

    csvFilename := flag.String("csv","problems.csv","a csv file in format of 'question,answer'")        // add csv file in argument default file is problems.csv
    timeLimit := flag.Int("limit", 30, "the time limit for your quiz in seconds")                       // flag for 30 seconds timer
    flag.Parse()

    file, err := os.Open(*csvFilename)  //csvFilename pointer to string, *csvFilename actual value
    if err != nil{
        exit(fmt.Sprintf("Falied to Open CSV File: %s\n", *csvFilename))
    }

    // Parsing CSV
    r := csv.NewReader(file)
    lines, err := r.ReadAll()
    if err != nil{
        exit("Failed to Parse the Provided CSV file.")
    }

    problems := parseLines(lines)

    fmt.Println("There are ", len(lines), "questions and you have", *timeLimit, "seconds to solve them\n")

    // creating new timer for timelimit  (*timeLimit) is int64 convert to time.Second
    timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

    // Staring Quiz
    correct := 0
    for i, p := range(problems){

        // if we are not answering the question then also exits the quiz
        fmt.Printf("Problem #%d: %s = ", i+1, p.q)
        answerCh := make(chan string)
        go func(){                                  // goroutine running anonymous function
            var answer string
            fmt.Scanf("%s\n", &answer)
            answerCh <-strings.ToLower(answer)                       // sending answer to answerCh
        }()

        select {
        case <-timer.C:                             // receives message from timer channel then stop
            fmt.Printf("\nTime up \nYou scored %d out of %d.\n", correct, len(problems))
            return                                  // return from select and from for loop

        case answer := <-answerCh :                 // gets answer from answerCh
            if answer == p.a {
                correct++
            }
        }
    }

    fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

// convert the 2d string slice into problem struct
func parseLines(lines [][]string) []problem {
    ret := make( []problem, len(lines))                                                   //slice of length equal to problems

    for i, line := range lines {
        ret[i] = problem{ q: strings.TrimSpace(line[0]), a: strings.ToLower(strings.TrimSpace(line[1])),}      // question is column 1 and answer is column 2
    }

    return ret
}

// structure for ques and answer
type problem struct {
    q string
    a string
}

// exit message
func exit(msg string){
    fmt.Println(msg)
    os.Exit(1)
}
