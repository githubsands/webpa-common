package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/wrp"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// Failures is a map type which records device routing failures
type Failures map[Interface]error

func (df Failures) Add(d Interface, deviceError error) {
	df[d] = deviceError
}

func (df Failures) WriteJSON(output io.Writer) (count int, err error) {
	_, err = fmt.Fprintf(output, `{"errors": [`)
	if err != nil {
		return
	}

	separator := ""
	for d, deviceError := range df {
		if deviceError != nil {
			count++
			fmt.Fprintf(output, `{"id": "%s", "key": "%s", error: "%s"}%s`, d.ID(), d.Key(), deviceError, separator)
			separator = ","
		}
	}

	_, err = fmt.Fprintf(output, `]}`)
	return
}

func (df Failures) MarshalJSON() (data []byte, err error) {
	var buffer bytes.Buffer
	_, err = df.WriteJSON(&buffer)
	data = buffer.Bytes()
	return
}

func (df Failures) WriteResponse(response http.ResponseWriter) error {
	if len(df) > 0 {
		var buffer bytes.Buffer
		if count, err := df.WriteJSON(&buffer); err != nil {
			return err
		} else if count > 0 {
			// only write the JSON out if any devices actually had errors, as opposed
			// to devices mapped to nil errors
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusInternalServerError)
			_, err := buffer.WriteTo(response)
			return err
		}
	}

	return nil
}

// NewTranscodingHandler produces an http.Handler that decodes the body of a request as a something other than
// Msgpack, e.g. JSON.  The exact format is determined by the supplied decoder.
func NewTranscodingHandler(decoderPool *wrp.DecoderPool, router Router) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		message := new(wrp.Message)
		if err := decoderPool.Decode(message, request.Body); err != nil {
			http.Error(
				response,
				fmt.Sprintf("Could not decode WRP message: %s", err),
				http.StatusBadRequest,
			)

			return
		}

		failures := make(Failures)
		if _, count, err := router.Route(message, nil, failures.Add); err != nil {
			http.Error(
				response,
				fmt.Sprintf("Could not route WRP message: %s", err),
				http.StatusBadRequest,
			)
		} else if count == 0 {
			response.WriteHeader(http.StatusNotFound)
		} else {
			failures.WriteResponse(response)
		}
	})
}

// NewMsgpackHandler produces an http.Handler that decodes the body of a request as a Msgpack WRP message
// and dispatches that message via the supplied Router.
func NewMsgpackHandler(decoderPool *wrp.DecoderPool, router Router) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(
				response,
				fmt.Sprintf("Unable to read request body: %s", err),
				http.StatusBadRequest,
			)

			return
		}

		message := new(wrp.Message)
		if err := decoderPool.DecodeBytes(message, body); err != nil {
			http.Error(
				response,
				fmt.Sprintf("Could not decode WRP message: %s", err),
				http.StatusBadRequest,
			)

			return
		}

		failures := make(Failures)
		if _, count, err := router.Route(message, body, failures.Add); err != nil {
			http.Error(
				response,
				fmt.Sprintf("Could not route WRP message: %s", err),
				http.StatusBadRequest,
			)
		} else if count == 0 {
			response.WriteHeader(http.StatusNotFound)
		} else {
			failures.WriteResponse(response)
		}
	})
}

// NewConnectHandler produces an http.Handler that allows devices to connect
// to a specific Manager.
func NewConnectHandler(connector Connector, responseHeader http.Header, logger logging.Logger) http.Handler {
	if logger == nil {
		logger = logging.DefaultLogger()
	}

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		device, err := connector.Connect(response, request, responseHeader)
		if err != nil {
			logger.Error("Failed to connect device: %s", err)
		} else {
			logger.Debug("Connected device: %s", device.ID())
		}
	})
}

// NewDeviceListHandler returns an http.Handler that renders a JSON listing
// of the devices within a manager.
func NewDeviceListHandler(manager Manager, logger logging.Logger) http.Handler {
	if logger == nil {
		logger = logging.DefaultLogger()
	}

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		flusher := response.(http.Flusher)
		response.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(response, `{"device": [`); err != nil {
			logger.Error("Unable to write content: %s", err)
			return
		}

		devices := make(chan Interface, 100)
		finish := new(sync.WaitGroup)
		finish.Add(1)

		// to minimize the time we hold the read lock on the Manager, spawn a goroutine
		// that collects devices and inserts them into an output buffer
		go func() {
			defer finish.Done()

			needsDelimiter := false
			for d := range devices {
				if needsDelimiter {
					io.WriteString(response, ",")
				}

				needsDelimiter = true
				if data, err := json.Marshal(d); err != nil {
					message := fmt.Sprintf("Unable to marshal device [%s] as JSON: %s", d.ID(), err)
					logger.Error(message)
					fmt.Fprintf(response, `"%s"`, message)
				} else {
					response.Write(data)
				}

				flusher.Flush()
			}
		}()

		manager.VisitAll(func(d Interface) {
			devices <- d
		})

		close(devices)
		finish.Wait()
		io.WriteString(response, `]}`)
		flusher.Flush()
	})
}