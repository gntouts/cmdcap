package cmdcap

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createLogDir(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		err := os.Mkdir(path, 0755)
		catch(err)
	}
	return path
}

func countFiles(path string) int {
	files, err := ioutil.ReadDir(path)
	catch(err)
	return len(files)
}

func fname(path string) string {
	files := countFiles(path)
	filesn := strconv.Itoa(files)
	return path + "log-" + filesn + ".log"
}

func createLogFile(path string) string {
	filename := fname(path)
	file, err := os.Create(filename)
	catch(err)
	defer file.Close()
	return file.Name()
}

func CaptureCmd(path string) {
	path = path + "logs"/
	createLogDir(path)
	logFile := createLogFile(path)

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0666)
	catch(err)
	defer file.Close()

	argsWithProg := os.Args

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, argsWithProg)

	w.Flush()
}
