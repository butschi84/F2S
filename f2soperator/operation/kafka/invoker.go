package kafka

import (
	"butschi84/f2s/state/queue"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func invokeFunction(functionUid string, requestBody string) (result string, err error) {
	logging.Info("request to invoke a function")

	// get this function from config
	f2sfunction, errGetFunction := f2shub.F2SConfiguration.GetFunctionByUID(functionUid)
	if errGetFunction != nil {
		logging.Error(fmt.Errorf("error getting function %s from running config. abort function invocation", functionUid))
		return
	}

	// make request obj
	request := queue.F2SRequest{
		UID:           f2shub.F2SEventManager.GenerateUUID(),
		Path:          f2sfunction.Spec.Endpoint,
		Method:        f2sfunction.Spec.Method,
		ResultChannel: make(chan queue.F2SRequestResult),
		Payload:       requestBody,
	}

	// put it into queue
	logging.Info(fmt.Sprintf("[%s] add request to queue", request.UID))
	f2shub.F2SQueue.AddRequest(&request)

	// wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	select {
	case result := <-request.ResultChannel:
		logging.Info(fmt.Sprintf("[%s] Request completed [success: %v]: %s", request.UID, result.Success, result.Result))
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.Encode(result)

		return buffer.String(), nil
	case <-ctx.Done():
		logging.Warn(fmt.Sprintf("[%s] Request Timeout reached, cancelling goroutine", request.UID))
		return "", fmt.Errorf("[%s] Request Timeout reached, cancelling goroutine", request.UID)
	}
}
