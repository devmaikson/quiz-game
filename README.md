# quiz-game

This is a little quiz game from: <https://gophercises.com>

This go program reads a csv file with questions and answers from a CSV file. Users can answer the questions from the CLI and get the results right after the timeout happened or all questions are answered before the timeout.

Example csv:

```text
5+5,10
7+3,10
1+1,2
8+3,11
1+2,3
8+6,14
3+1,4
1+4,5
5+1,6
2+3,5
3+3,6
2+4,6
5+2,7
"Juicy company?","apple"
```

Commandline help:

```bash
go run main.go -help
  -f string
        use to specify a CSV filename (default "examples/problems.csv")
  -shuffle
        boolean to enable shuffling of the questions
  -t int
        time in seconds before the timeout (default 30)
```
