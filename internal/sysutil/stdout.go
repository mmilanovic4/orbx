package sysutil

import "os"

// WriteStdout writes data to stdout byte for byte, so decoded and binary
// output survives pipes and redirects unchanged. Only on a terminal is a
// missing trailing newline added, to keep the prompt off the output line.
func WriteStdout(data []byte) error {
	return writeOutput(os.Stdout, data)
}

func writeOutput(f *os.File, data []byte) error {
	if _, err := f.Write(data); err != nil {
		return err
	}
	if len(data) > 0 && data[len(data)-1] == '\n' {
		return nil
	}
	if stat, err := f.Stat(); err == nil && stat.Mode()&os.ModeCharDevice != 0 {
		_, err := f.Write([]byte{'\n'})
		return err
	}
	return nil
}
