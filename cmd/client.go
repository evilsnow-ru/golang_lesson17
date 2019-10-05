package cmd

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"golang_lesson17/api"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	layout         = "2006-01-02 15:04:05"
	now            = "now"
	errEmptyMsg    = "Error: message is empty"
	errNoMsgParam  = "Error: no message parameter"
	errTooManyArgs = "Error: too many arguments. Message must be decorated by quotes."
	errNoIdParam   = "Error: no id param set"
)

var portFlag uint32
var hostFlag string
var eventTypeFlag int32
var dateFlag string

var eventAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adding message to storage. Usage: add [phtd options] \"message in quotes\"",
	Run: func(cmd *cobra.Command, args []string) {
		argsLen := len(args)

		if argsLen == 0 {
			exit(errNoMsgParam)
		}

		if argsLen > 1 {
			exit(errTooManyArgs)
		}

		_, ok := api.EventType_name[eventTypeFlag]

		if !ok {
			exitf("Event type %d doesn't supported. Available types are:\n%s\n", eventTypeFlag, getSupportedTypesString())
		}

		var date time.Time
		var err error

		if dateFlag == now {
			date = time.Now()
		} else {
			date, err = time.Parse(layout, dateFlag)
			if err != nil {
				exit(err.Error())
			}
		}

		msg := args[0]
		if msg == "" {
			exit(errEmptyMsg)
		}

		event := &api.Event{
			MsgId: 0,
			Date: &timestamp.Timestamp{
				Seconds: date.Unix(),
			},
			Type:        api.EventType(eventTypeFlag),
			Description: msg,
		}

		data, err := proto.Marshal(event)

		if err != nil {
			exit(err.Error())
		}

		fmt.Printf("Sending data to %s:%d\n", hostFlag, portFlag)
		response, err := http.Post(fmt.Sprintf("http://%s:%d/add", hostFlag, portFlag), "application/octet-stream", bytes.NewReader(data))
		processResponse(response, err)
	},
}

var getEventCmd = &cobra.Command{
	Use:   "get",
	Short: "Get existing event data from service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exit(errNoIdParam)
		}

		id, err := strconv.Atoi(args[0])

		if err != nil {
			exit(err.Error())
		}

		if id < 0 {
			exitf("Id must be >= 0. Actual: %d.", id)
		}

		fmt.Printf("Requesting event by id: %d\n", id)
		response, err := http.Get(fmt.Sprintf("http://%s:%d/get?id=%d", hostFlag, portFlag, id))
		processResponse(response, err)
	},
}

func processResponse(response *http.Response, err error) {
	if err != nil {
		exit(err.Error())
	}
	if response.StatusCode == 200 {
		fmt.Println("Send success")
	} else {
		fmt.Printf("Response code: %d\n", response.StatusCode)
	}
	if response.Body != nil {
		defer response.Body.Close()

		resultData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(resultData))
		}
	}
}

func init() {
	eventAddCmd.Flags().Uint32VarP(&portFlag, "port", "p", defaultPort, "Set port")
	eventAddCmd.Flags().StringVarP(&hostFlag, "host", "", "localhost", "Set host")
	eventAddCmd.Flags().Int32VarP(&eventTypeFlag, "type", "t", int32(api.EventType_UNDEFINED), "Set event type")
	eventAddCmd.Flags().StringVarP(&dateFlag, "date", "d", now, "Set event date")

	getEventCmd.Flags().Uint32VarP(&portFlag, "port", "p", defaultPort, "Set port")
	getEventCmd.Flags().StringVarP(&hostFlag, "host", "", "localhost", "Set host")
}

func exitf(template string, params ...interface{}) {
	exit(fmt.Sprintf(template, params))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func getSupportedTypesString() string {
	sortedData := make([]string, len(api.EventType_name))
	var buf strings.Builder

	for key, value := range api.EventType_name {
		sortedData[int(key)] = fmt.Sprintf("\t%d: %s\n", key, value)
	}

	for _, value := range sortedData {
		buf.WriteString(value)
	}

	return buf.String()
}
