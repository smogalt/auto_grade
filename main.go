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
	"sort"
	"strings"
)

func main() {
	verbose := flag.Bool("v", false, "shows what testcases failed")
	pass_rate := flag.Bool("pr", false, "provides detailed information on the pass rates of tests")
	tests_fn := flag.String("t", "tests.txt", "filename for tests")
	answers_fn := flag.String("a", "answers.txt", "filename for answer key")
	all_fail := flag.Bool("af", false, "shows test cases where all programs failed")
	flag.Parse()

	/* open file with tests */
	tests, err := os.Open(*tests_fn)
	if err != nil {
		log.Fatal(err)
	}
	defer tests.Close()
	tests_line := bufio.NewScanner(tests)

	/* open file with answers */
	answers, err := os.Open(*answers_fn)
	if err != nil {
		log.Fatal(err)
	}
	defer answers.Close()
	answers_line := bufio.NewScanner(answers)

	/* check to make sure each test has a response */
	if count_lines(answers) != count_lines(tests) {
		log.Fatal("Unequal amount of tests and responses.")
	}

	/* list all dirs under given dir */
	dirs, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	/* regexp to check for .py fils in the dir */
	var valid_py = regexp.MustCompile(`.*py$`)

	/* create map with scores and directory names */
	var report map[string]int
	report = make(map[string]int)

	var test_report map[int]int
	test_report = make(map[int]int)
	passed := 0

	tests_num := 0
	/* for each test case */
	for tests_line.Scan() {
		passed = 0
		tests_num++
		answers_line.Scan()
		/* for each dir */
		for _, e := range dirs {
			/* list all files under e */
			programs, err := os.ReadDir("./" + e.Name())
			if err != nil {
				log.Fatal(err)
			}
			if len(programs) == 0 {
				report[e.Name()] += 0
				continue
			}

			/* find first valid .py file under e */
			index := 0
			for i, p := range programs {
				if valid_py.MatchString(p.Name()) {
					index = i
					break
				}
			}

			/* run given .py files with input */
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := exec.Command("python3.11", e.Name()+"/"+programs[index].Name(), tests_line.Text())
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				fmt.Println(stderr.String())
				log.Fatal(err)
			}

			stdout_string := stdout.String()
			stdout_string = strings.TrimSuffix(stdout_string, "\n")

			if stdout_string != answers_line.Text() {
				report[e.Name()] += 0
				if *verbose {
					fmt.Println(e.Name(), "failed: test", tests_num)
					fmt.Println("\texpected response:", answers_line.Text())
					fmt.Println("\tresponse:", stdout_string)
					fmt.Println()
				}
			} else {
				report[e.Name()]++
				passed++
			}
		}
		test_report[tests_num] += passed
	}

	/* print names and scores. extra stuff is to sort the map */
	names := make([]string, 0, len(report))
	for k := range report {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		fmt.Printf("%s %d%%\n", k, report[k]*100/tests_num)
	}

	/* print the pass rates of individual tests */
	if *pass_rate {
		fmt.Println("pass rates:")
		keys := make([]int, 0, len(test_report))
		for k := range test_report {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			fmt.Printf("\t %d %d%%\n", k, (test_report[k]*100)/len(dirs))
		}
	}

	/* print tests that all programs failed */
	if *all_fail {
		fmt.Println("unsuccessful tests:")
		for k, v := range test_report {
			if v == 0 {
				fmt.Println("\t", k)
			}
		}
	}
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
