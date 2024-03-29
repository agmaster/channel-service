package middlewares

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"
)

// Convert document to plain text
func DocToText(in io.Reader, out io.Writer, logFile string) error {
	log.SetLogger("file", logFile)
	log.Trace("run middlewares.DocToText")

	cmd := exec.Command("java", "-jar", "./lib/tika-app-1.7.jar", "-t")
	stderr := bytes.NewBuffer(nil)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = stderr

	cmd.Start()
	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(500000) * time.Millisecond):
		if err := cmd.Process.Kill(); err != nil {
			return errors.New(err.Error())
		}
		<-cmdDone
		return errors.New("Command timed out")
	case err := <-cmdDone:
		if err != nil {
			return errors.New(stderr.String())
		}
	}

	return nil
}

func IsSupport(fileName string) bool {
	paths := strings.Split(fileName, "/")
	fileType := strings.Split(paths[len(paths)-1], ".")

	switch fileType[len(fileType)-1] {
	case "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "epub", "html", "xml":
		return true
	}
	return false

}

// Ref
// 1. http://rny.io/rails/elasticsearch/2013/08/05/full-text-search-for-attachments-with-rails-and-elasticsearch.html
// 2. Wrapper Apache Tika https://github.com/plimble/gika
