package main 

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	op := flag.String("op", "sum", "operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to be executed")
	flag.Parse()
	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileNames []string, op string, col int, out io.Writer) error {
	var opFunc statsFunc 
	if len(fileNames) == 0 {
		return ErrNoFiles
	}
	if col<1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, col)
	}
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg 
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)	
	}
	consolidate := make([]float64, 0)
	for _, fname := range fileNames{
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("Cannot open file: %w", err)
		}
		data, err := csv2float(f, col)
		if err != nil {
			return err 
		}
		if err := f.Close(); err != nil {
			return err 
		}
		consolidate = append(consolidate, data...)
	}
	_, err := fmt.Fprintln(out, opFunc(consolidate))
	return err
}