package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	// OutputConsole should we output this message to the console?
	OutputConsole = true

	// OutputAsReadable change the output of objects to readable strings on multiple lines.
	OutputAsReadable = true

	// OutputFilename filename to log to
	OutputFilename = ""

	// ErrorWebhook Webhook to write error output to.
	ErrorWebhook = ""

	// SingleFile Should we only show the last file
	SingleFile = false

	// SkipTime time prefix
	SkipTime = false

	// IncludeFunction ...
	IncludeFunction = true

	// SkipDate should we display without date
	SkipDate = false

	// SkipSeverity should we skip showing [INFO], [DEBUG], [WARNING]... [ERROR] is always shown
	SkipSeverity = false

	// SkipFile should we skip showing filename
	SkipFile = false

	// ApplicationTitle string
	ApplicationTitle = ""

	// OutputJSON filename to log json data to
	OutputJSON = ""

	// LenAlways ...
	LenAlways = -1

	// private output file pointer variable, keeps open.
	outputFile *os.File
	outputJSON *os.File

	// OutputDebug should the system write debug logs?
	OutputDebug = false

	// OutputInfo should the system write info logs?
	OutputInfo = true

	// OutputWarning should the system write warning logs?
	OutputWarning = true

	// IgnoredFiles is files to ignore
	IgnoredFiles []string

	// L is the internal logger
	L *log.Logger
)

type jsonLOG struct {
	Time          string      `json:"time,omitempty"`
	Files         string      `json:"files,omitempty"`
	Severity      string      `json:"severity,omitempty"`
	User          string      `json:"user,omitempty"`
	ExecutionTime string      `json:"executionTime,omitempty"`
	Message       interface{} `json:"message,omitempty"`
}

const (
	// InfoColor ...
	InfoColor = "\033[0;32m%s\033[0m"

	// NoticeColor ...
	NoticeColor = "\033[0;36m%s\033[0m"

	// WarningColor ...
	WarningColor = "\033[0;33m%s\033[0m"

	// ErrorColor ...
	ErrorColor = "\033[0;31m%s\033[0m"

	// DebugColor ...
	DebugColor = "\033[0;34m%s\033[0m"
)

// NewCustom ...
func NewCustom() {
	if SkipDate && SkipTime {
		L = log.New(os.Stdout, "", 0)
	} else if SkipDate {
		L = log.New(os.Stdout, "", log.Ltime)
	} else {
		L = log.New(os.Stdout, "", log.Ltime|log.Ldate)
	}
}

// Parse ...
func (j *jsonLOG) Parse() {
	if _, ok := j.Message.(string); ok {

		a := strings.LastIndex(j.Message.(string), "(")
		b := strings.LastIndex(j.Message.(string), ")")
		if a > 0 && b > 0 {
			submsg := j.Message.(string)[a+1 : b]
			submsgs := strings.Split(submsg, ", ")
			for _, subm := range submsgs {
				txt := strings.Split(subm, ": ")
				if len(txt) > 1 {
					switch txt[0] {
					case "User":
						j.User = txt[1]
					case "Execution Time":
						j.ExecutionTime = txt[1]
					}
				}
			}
		}
	}
}

// Printf outputs as Info
func Printf(message interface{}, args ...interface{}) {
	Info(message, args...)
}

// Debug outputs a debug message to the output and/or logfile
func Debug(message interface{}, args ...interface{}) {
	if OutputDebug {
		outputMessage(message, "DEBUG", args...)
	}
}

// Info outputs a info message to the output and/or logfile
func Info(message interface{}, args ...interface{}) {
	if OutputInfo {
		outputMessage(message, "INFO", args...)
	}
}

// Warning outputs a warning message to the output and/or logfile
func Warning(message interface{}, args ...interface{}) {
	if OutputWarning {
		outputMessage(message, "WARNING", args...)
	}
}

// Error outputs a error message to the default output and/or logfile
func Error(message interface{}, args ...interface{}) {
	outputMessage(message, "ERROR", args...)
}

// Fatalf outputs a fatal message to the output and/or logfile
func Fatalf(message interface{}, args ...interface{}) {
	outputMessage(message, "FATAL", args...)
}

// ClearLog clears all logfiles
func ClearLog() {
	if OutputFilename != "" {
		if outputFile != nil {
			Close()
		}
		os.Remove(OutputFilename)
	}

	if OutputJSON != "" {
		os.Remove(OutputJSON)
	}
}

// Close closes the logfile if it is open.
func Close() {
	if outputFile != nil {
		outputFile.Close()
	}
	outputFile = nil
}

func getInterfaceAsString(message interface{}, indent int) string {

	//name := reflect.TypeOf(message).Name()

	switch messageCast := message.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if indent > 2 {
			return fmt.Sprintf("%d\n", messageCast)
		}
		return fmt.Sprintf("%d", messageCast)

	case string:
		if indent > 2 {
			return fmt.Sprintf("%s\n", messageCast)
		}
		return fmt.Sprintf("%s", messageCast)

	default:
		/*
			if OutputAsReadable {
				var isPtr bool
				var returnValue string

				messageString := ""
				val := reflect.ValueOf(message)
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
					isPtr = true
					name = val.Type().Name()
				}

				numFields := val.Type().NumField()
				for i := 0; i < numFields; i++ {
					field := val.Type().Field(i)

					if val.CanInterface() {
						ifaceValue := val.Field(i).Interface()
						for j := 0; j < indent; j++ {
							messageString += " "
						}
						messageString += field.Name + ": " + getInterfaceAsString(ifaceValue, indent+2)
					}
				}

				if isPtr {
					returnValue = fmt.Sprintf("%s (*) {\n%s", name, messageString)
				} else {
					returnValue = fmt.Sprintf("%s {\n%s", name, messageString)
				}
				for j := 0; j < indent-2; j++ {
					returnValue += " "
				}
				returnValue += "}"
				if indent > 2 {
					returnValue += "\n"
				}
				return returnValue
			}
			return fmt.Sprintf("%+v", message)
		*/
		data, err := json.Marshal(message)
		if err != nil {
			L.Printf("err?")
			return fmt.Sprintf("%+v", message)
		}
		return fmt.Sprintf("%s", data)
	}
}

// outputMessage outputs the log to file and/or console.
func outputMessage(message interface{}, messageType string, args ...interface{}) {
	var err error
	var files string

	if L == nil {
		NewCustom()
	}

	if SingleFile {

		ptr, file, no, _ := runtime.Caller(2)
		file = file[(strings.LastIndex(file, "/") + 1):]
		//programCounters := make([]uintptr, 2)
		//n := runtime.Callers(2, programCounters)
		if IncludeFunction {
			frames := runtime.CallersFrames([]uintptr{ptr})
			f, _ := frames.Next()
			function := f.Function[(strings.LastIndex(f.Function, ".") + 1):]
			files = fmt.Sprintf("%v:%d:%s()", file, no, function)
		} else {
			files = fmt.Sprintf("%v:%d", file, no)
		}
	} else {

		for i := 2; i < 10; i++ {
			ptr, file, no, _ := runtime.Caller(i)
			file = file[(strings.LastIndex(file, "/") + 1):]
			if file == "" ||
				file == "<autogenerated>" ||
				file == "asm_amd64.s" ||
				file == "server.go" ||
				file == "mux.go" ||
				file == "proc.go" ||
				file == "testing.go" {
				break
			}
			if file != "go-logger.go" {
				ignore := false
				for _, ignoredFile := range IgnoredFiles {
					if ignoredFile == file {
						ignore = true
						break
					}
				}
				if !ignore {
					if files != "" {
						files = "," + files
					}
					if IncludeFunction && i == 2 {
						frames := runtime.CallersFrames([]uintptr{ptr})
						f, _ := frames.Next()
						function := f.Function[(strings.LastIndex(f.Function, ".") + 1):]
						files = fmt.Sprintf("%s:%d%s:%s()", file, no, files, function)
					} else {
						files = fmt.Sprintf("%s:%d%s", file, no, files)
					}
				}
			}
		}
	}

	if OutputConsole {

		empty := ""
		if LenAlways > 0 {

			textLen := len(fmt.Sprintf("[%s]", files))
			if textLen > LenAlways {
				files = files[textLen-LenAlways:]
			}
			for i := LenAlways - textLen; i > 0; i-- {
				empty += " "
			}
		}

		fileAndDate := fmt.Sprintf("[%s]%s", files, empty)
		if SkipSeverity == false {
			fileAndDate = fmt.Sprintf("[%s] [%s]%s", files, messageType, empty)
		}

		if SkipFile {
			fileAndDate = strings.Replace(fileAndDate, fmt.Sprintf("[%s] ", files), "", 1)
		}

		info := ""
		switch messageType {
		case "ERROR":
			info = fmt.Sprintf(ErrorColor, fmt.Sprintf("%s%s", fileAndDate, empty))
		case "INFO":
			info = fmt.Sprintf(InfoColor, fmt.Sprintf("%s%s", fileAndDate, empty))
		case "DEBUG":
			info = fmt.Sprintf(DebugColor, fmt.Sprintf("%s%s", fileAndDate, empty))
		case "WARNING":
			info = fmt.Sprintf(WarningColor, fmt.Sprintf("%s%s", fileAndDate, empty))
		case "NOTICE":
			info = fmt.Sprintf(NoticeColor, fmt.Sprintf("%s%s", fileAndDate, empty))
		}

		switch messageCast := message.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			L.Printf("%s %d", info, messageCast)
			break
		default:
			if len(args) > 0 {
				L.Printf("%s %s", info, fmt.Sprintf(message.(string), args...))
			} else {
				L.Printf("%s %+v", info, message)
			}
			break
		}
	}

	y, m, d := time.Now().Local().Date()
	hh, mm, ss := time.Now().Local().Clock()

	if OutputJSON != "" {
		msg := ""
		var jlog jsonLOG
		if outputJSON == nil {
			if _, err := os.Stat(OutputJSON); os.IsNotExist(err) {
				outputJSON, err = os.OpenFile(OutputJSON, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
				jlog = jsonLOG{Time: fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", y, m, d, hh, mm, ss), Files: "go-logger.go:190", Severity: "INFO", Message: "Starting new JSON log file"}
				data, _ := json.Marshal(jlog)
				//msg = fmt.Sprintf("[\n\t{\n\t\t\"time\": \"%04d-%02d-%02d %02d:%02d:%02d\",\n\t\t\"files\": \"%s\",\n\t\t\"severity\": \"%s\",\n\t\t\"message\": \"%s\"\n\t}\n]", y, m, d, hh, mm, ss, "go-logger.go:180", "INFO", "Starting new JSON logfile.")
				msg = fmt.Sprintf("[\n\t%s\n]", string(data))
				outputJSON.WriteString(msg)
				outputJSON.Close()
				outputJSON = nil
			}
			outputJSON, err = os.OpenFile(OutputJSON, os.O_WRONLY, 0755)
			if err != nil {
				messageTmp := fmt.Sprintf("Unable to open file %s", OutputJSON)
				OutputJSON = ""
				outputMessage(messageTmp, "ERROR")
			}
		}
		if outputJSON != nil {
			outputJSON.Seek(-2, 2)

			if len(args) > 0 {
				jlog = jsonLOG{Time: fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", y, m, d, hh, mm, ss), Files: files, Severity: messageType, Message: fmt.Sprintf(message.(string), args...)}
				//msg = fmt.Sprintf(",\n\t{\n\t\t\"time\": \"%04d-%02d-%02d %02d:%02d:%02d\",\n\t\t\"files\": \"%s\",\n\t\t\"severity\": \"%s\",\n\t\t\"message\": \"%s\"\n\t}\n]", y, m, d, hh, mm, ss, files, messageType, fmt.Sprintf(message.(string), args...))
			} else {
				//messageValue := getInterfaceAsString(message, 2)
				jlog = jsonLOG{Time: fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", y, m, d, hh, mm, ss), Files: files, Severity: messageType, Message: message}
				//msg = fmt.Sprintf(",\n\t{\n\t\t\"time\": \"%04d-%02d-%02d %02d:%02d:%02d\",\n\t\t\"files\": \"%s\",\n\t\t\"severity\": \"%s\",\n\t\t\"message\": \"%s\"\n\t}\n]", y, m, d, hh, mm, ss, files, messageType, messageValue)
			}
			jlog.Parse()
			data, _ := json.Marshal(jlog)
			msg = fmt.Sprintf(",\n\t%s\n]", string(data))
			outputJSON.WriteString(msg)
			outputJSON.Sync()
		}
	}

	if OutputFilename != "" {
		if outputFile == nil {
			outputFile, err = os.OpenFile(OutputFilename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
			if err != nil {
				messageTmp := fmt.Sprintf("Unable to open file %s", OutputFilename)
				OutputFilename = ""
				outputMessage(messageTmp, "ERROR")
			}
		}
		if outputFile != nil {
			if len(args) > 0 {
				outputFile.WriteString(fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d] [%s] [%s] :: %s\n", y, m, d, hh, mm, ss, files, messageType, fmt.Sprintf(message.(string), args...)))
			} else {
				messageValue := getInterfaceAsString(message, 2)
				outputFile.WriteString(fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d] [%s] [%s] :: %s\n", y, m, d, hh, mm, ss, files, messageType, messageValue))
			}
			outputFile.Sync()
		}
	}

	if ErrorWebhook != "" {

		if messageType == "ERROR" {
			func(message interface{}, ErrorWebhook string) {
				text := ""
				if len(args) > 0 {
					text = fmt.Sprintf("%s", fmt.Sprintf(message.(string), args...))
				} else {
					text = fmt.Sprintf("%+v", message)
				}

				card := CreateMessageCard(text)
				card.AddSectionWithFacts("", true, map[string]string{"Application": ApplicationTitle, "Files": files})
				card.AddSectionWithText("", false, text)

				data, err := json.Marshal(card)
				if err != nil {
					L.Printf("[%s] [%s] :: %+v", "go-logger...", "DEBUG", err)
					return
				}

				log.Printf("%+v", string(data))

				req, err := http.NewRequest("POST", ErrorWebhook, bytes.NewBuffer(data))
				req.Header.Set("Content-Type", "application/json")

				/*
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						log.Printf("[%s] [%s] :: %+v", "go-logger...", "DEBUG", err)
						return
					}
					defer resp.Body.Close()
				*/
			}(message, ErrorWebhook)
		}
	}
}
