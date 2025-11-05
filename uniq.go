package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Params struct {
	C          bool
	D          bool
	U          bool
	Num_fields int
	Num_chars  int
	I          bool
}

func countOccurences(lines []string, s string, p Params, start int) int {
	var count = 0
	for i := start; i < len(lines); i++ {
		if checkLines(lines[i], p) == s {
			count++
		}
	}
	return count
}

func checkLines(s string, p Params) string {
	res := s
	var n = p.Num_fields
	if n > 0 {
		fields := strings.Fields(res)
		if len(fields) > n {
			res = strings.Join(fields[n:], " ")
		} else {
			res = ""
		}
	}
	if p.Num_chars > 0 {
		if len(res) > p.Num_chars {

			res = res[p.Num_chars:]
		} else {
			res = ""
		}
	}
	if p.I {
		res = strings.ToLower(res)
	}
	return res
}

func checkCDU(lines []string, p Params) []string {
	if len(lines) == 0 {
		return nil
	} else {
		var res []string
		pastlines := make(map[string]bool)
		for i, line := range lines {
			s := checkLines(line, p)
			if pastlines[s] == false {
				count := countOccurences(lines, s, p, i)
				switch {
				case p.C:
					res = append(res, fmt.Sprintf("%d %s", count, line))
				case p.D && count > 1:
					res = append(res, line)
				case p.U && count == 1:
					res = append(res, line)
				case p.D == false && p.U == false:
					res = append(res, line)
				}
				pastlines[s] = true
			}
		}
		return res
	}
}

func checkFiles(p Params, inputfile, outputfile string) error {
	input := os.Stdin
	if inputfile != "" {
		file, err := os.Open(inputfile)
		if err != nil {
			return err
		}
		defer file.Close()
		input = file
	}

	scanner := bufio.NewScanner(input)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	result := checkCDU(lines, p)

	output := os.Stdout
	if outputfile != "" {
		file, err := os.Create(outputfile)
		if err != nil {
			return err
		}
		defer file.Close()
		output = file
	}

	for _, line := range result {
		fmt.Fprintln(output, line)
	}

	return nil
}

func main() {
	var params Params

	flag.BoolVar(&params.C, "c", false, "подсчитать количество встречаний строки во входных данных. Вывести это число перед строкой отделив пробелом")
	flag.BoolVar(&params.D, "d", false, "Вывести только те строки, которые повторились во входных данных")
	flag.BoolVar(&params.U, "u", false, "вывести только те строки, которые не повторились во входных данных")
	flag.BoolVar(&params.I, "i", false, "не учитывать регистр букв")
	flag.IntVar(&params.Num_fields, "f", 0, "не учитывать первые num_fields полей в строке")
	flag.IntVar(&params.Num_chars, "s", 0, "не учитывать первые num_chars символов в строке")

	flag.Parse()

	if (params.C && params.D) || (params.D && params.U) || (params.C && params.U) {
		fmt.Println("Параметры -c -d -u взаимнозаменяемы, нет смысла использовать их параллельно")
	}

	data := flag.Args()

	var inputfile, outputfile string

	if len(data) > 0 {
		inputfile = data[0]
	}
	if len(data) > 1 {
		outputfile = data[1]
	}
	if len(data) > 2 {
		fmt.Fprintln(os.Stderr, "слишком много аргументов")
		os.Exit(1)
	}

	if err := checkFiles(params, inputfile, outputfile); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

}
