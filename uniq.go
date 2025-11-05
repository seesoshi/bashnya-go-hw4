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

func processLine(line string, p Params) string {
	result := line

	// Обработка -f (поля)
	if p.Num_fields > 0 {
		fields := strings.Fields(result)
		if p.Num_fields < len(fields) {
			result = strings.Join(fields[p.Num_fields:], " ")
		} else {
			result = ""
		}
	}
	// Обработка -s (символы)
	if p.Num_chars > 0 {
		if p.Num_chars < len(result) {
			result = result[p.Num_chars:]
		} else {
			result = ""
		}
	}

	// Обработка -i (регистр)
	if p.I {
		result = strings.ToLower(result)
	}

	return result
}

func uniqLines(lines []string, p Params) []string {
	if len(lines) == 0 {
		return nil
	}

	// Собираем полную статистику
	type lineInfo struct {
		count     int
		firstLine string
	}
	groups := make(map[string]*lineInfo)

	for _, line := range lines {
		key := processLine(line, p)
		if info, exists := groups[key]; exists {
			info.count++
		} else {
			groups[key] = &lineInfo{
				count:     1,
				firstLine: line,
			}
		}
	}

	var result []string
	seen := make(map[string]bool)

	for _, line := range lines {
		key := processLine(line, p)
		if seen[key] {
			continue
		}
		seen[key] = true

		info := groups[key]

		// Применяем фильтры
		switch {
		case p.C:
			result = append(result, fmt.Sprintf("%d %s", info.count, info.firstLine))
		case p.D:
			if info.count > 1 {
				result = append(result, info.firstLine)
			}
		case p.U:
			if info.count == 1 {
				result = append(result, info.firstLine)
			}
		default:
			result = append(result, info.firstLine)
		}
	}

	return result
}

func processFiles(p Params, inputfile, outputfile string) error {
	// Вход
	input := os.Stdin
	if inputfile != "" {
		file, err := os.Open(inputfile)
		if err != nil {
			return err
		}
		defer file.Close()
		input = file
	}

	// Чтение
	scanner := bufio.NewScanner(input)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Обработка
	result := uniqLines(lines, p)

	// Выход
	output := os.Stdout
	if outputfile != "" {
		file, err := os.Create(outputfile)
		if err != nil {
			return err
		}
		defer file.Close()
		output = file
	}

	// Запись
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
		fmt.Fprintln(os.Stderr, "Параметры -c -d -u взаимнозаменяемы, нет смысла использовать их параллельно")
		os.Exit(1)
	}

	args := flag.Args()
	var inputfile, outputfile string

	if len(args) > 0 {
		inputfile = args[0]
	}
	if len(args) > 1 {
		outputfile = args[1]
	}
	if len(args) > 2 {
		fmt.Fprintln(os.Stderr, "слишком много аргументов")
		os.Exit(1)
	}

	if err := processFiles(params, inputfile, outputfile); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}
