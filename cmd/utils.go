package cmd

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/agilebits/urlreader"
)

func isStdinAvailable() (bool, error) {
	s, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	return s.Mode()&os.ModeNamedPipe != 0, nil
}

func getURL(args []string) (string, error) {
	useStdin, err := isStdinAvailable()
	if err != nil {
		return "", err
	}

	if useStdin {
		return "", nil
	}

	if len(args) == 0 {
		return "", errors.New("missing file name or url")
	}

	return args[0], nil
}

func isFileURL(url string) bool {
	return strings.HasPrefix(url, "file://") || !strings.Contains(url, "://")
}

func open(url string) (io.Reader, error) {
	var reader io.Reader

	if url == "" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		var err error
		reader, err = urlreader.Open(url)
		if err != nil {
			return nil, err
		}
	}

	return reader, nil
}

func read(url string) ([]byte, error) {
	reader, err := open(url)
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if closer, ok := reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func write(path string, body []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := file.Write(body); err != nil {
		return err
	}

	if err := file.Truncate(int64(len(body))); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
