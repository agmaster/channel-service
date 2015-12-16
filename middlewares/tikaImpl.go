package middlewares

import (
	"bytes"
	"errors"
	Logger "github.com/astaxie/beego/logs"
	"io"
	"os/exec"
	"strings"
	"time"
)

var log = Logger.NewLogger(10000)
var logFileName = `{"filename":"channel-service.log"}`

// Convert document to plain text
func DocToText(in io.Reader, out io.Writer) error {
	log.SetLogger("file", logFileName)
	log.Trace("DocToText: Convert document to plain text")

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
