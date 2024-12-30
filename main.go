package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func main() {
	/* define and parse arguments */
	verbose := flag.Bool("v", false, "shows what testcases failed")
	test_fn := flag.String("t", "tests.txt", "filename for tests")
	answer_fn := flag.String("a", "answers.txt", "filename for answer key")
	flag.Parse()

	/* open files with tests and answers*/
	test_file, err := os.Open(*test_fn)

	if err != nil {
		log.Fatal(err)
	}
	test_file_size := count_lines(test_file)

	answer_file, err := os.Open(*answer_fn)
	if err != nil {
		log.Fatal(err)
	}
	answer_file_size := count_lines(answer_file)

	/* check to make sure each test has a answer */
	if answer_file_size != test_file_size {
		log.Fatal("Unequal amount of tests and responses.")
	}

	/* load tests and answers */
	var tests, answers []string
	tests = make([]string, test_file_size, test_file_size)
	answers = make([]string, answer_file_size, answer_file_size)

	tests_line := bufio.NewScanner(test_file)
	answers_line := bufio.NewScanner(answer_file)

	for i := range tests {
		tests_line.Scan()
		answers_line.Scan()
		tests[i] = tests_line.Text()
		answers[i] = answers_line.Text()
	}

	test_file.Close()
	answer_file.Close()
	/* don't need the files open past here */

	/* list all dirs under working directory */
	dirs, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	/* regexp to check for .py fils in the dir */
	var valid_py = regexp.MustCompile(`.*py$`)

	/* create slice for scores */
	var scores []int
	scores = make([]int, len(dirs), len(dirs))

	/* one goroutine for each directory */
	var wait_group sync.WaitGroup
	wait_group.Add(len(dirs))

	/* for each test case */
	for n, d := range dirs {
		cur_dir := n
		cur_dir_name := d.Name()
		go func() {
			defer wait_group.Done()

			/* list all files under dir */
			programs, err := os.ReadDir("./" + cur_dir_name)
			if err != nil {
				log.Fatal(err)
			}
			if len(programs) != 0 {
				scores[cur_dir] += 0

				/* find first valid .py file under the current dir */
				index := 0
				for i, p := range programs {
					if valid_py.MatchString(p.Name()) {
						index = i
						break
					}
				}

				for test_num := range tests {
					/* run given .py files with input */
					var stdout bytes.Buffer
					var stderr bytes.Buffer

					cmd := exec.Command("python3", cur_dir_name+"/"+programs[index].Name(), tests[test_num])
					cmd.Stdout = &stdout
					cmd.Stderr = &stderr
					err = cmd.Run()
					if err != nil {
						fmt.Println(stderr.String())
						log.Fatal(err)
					}

					stdout_string := stdout.String()
					stdout_string = strings.TrimSuffix(stdout_string, "\n")

					if stdout_string != answers[test_num] {
						scores[cur_dir] += 0
						if *verbose {
							fmt.Println(cur_dir_name, "failed: test", test_num + 1)
							fmt.Println("\texpected response:", answers[test_num])
							fmt.Println("\tresponse:", stdout_string)
							fmt.Println()
						}
					} else {
						scores[cur_dir]++
					}
				}
			}
		}()
	}

	wait_group.Wait()

	for i, _ := range dirs {
		fmt.Printf("%s %d%%\n", dirs[i].Name(), (scores[i]*100)/len(tests))
	}

	return
}

func count_lines(f *os.File) int {
	lines := 0
	f.Seek(0, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines++
	}
	f.Seek(0, 0)
	return lines
}
