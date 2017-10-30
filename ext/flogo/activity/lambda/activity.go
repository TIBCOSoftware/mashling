package lambda

import (
	"encoding/json"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	opentracing "github.com/opentracing/opentracing-go"
	ctx "golang.org/x/net/context"
)

var log = logger.GetLogger("activity-tibco-lambda")

const (
	ivArn       = "arn"
	ivRegion    = "region"
	ivAccessKey = "accessKey"
	ivSecretKey = "secretKey"
	ivTracing   = "tracing"
	ivPayload   = "payload"

	ovValue   = "value"
	ovResult  = "result"
	ovStatus  = "status"
	ovTracing = "tracing"
)

// LambdaResponse struct is used to store the response from Lambda
type LambdaResponse struct {
	Status  int64
	Payload []byte
}

// LambdaActivity is a App Activity implementation
type LambdaActivity struct {
	sync.Mutex
	metadata *activity.Metadata
}

// NewActivity creates a new LambdaActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &LambdaActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *LambdaActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *LambdaActivity) Eval(context activity.Context) (done bool, err error) {
	var span opentracing.Span
	if tracing := context.GetInput(ivTracing); tracing != nil {
		span = opentracing.SpanFromContext(tracing.(ctx.Context))
	}

	if span != nil {
		span = opentracing.StartSpan(
			context.TaskName(),
			opentracing.ChildOf(span.Context()))
		context.SetOutput(ovTracing, opentracing.ContextWithSpan(ctx.Background(), span))
		defer span.Finish()
	}

	setTag := func(key string, value interface{}) {
		if span != nil {
			span.SetTag(key, value)
		}
	}

	arn := context.GetInput(ivArn).(string)
	setTag("arn", arn)

	var accessKey, secretKey = "", ""
	if context.GetInput(ivAccessKey) != nil {
		accessKey = context.GetInput(ivAccessKey).(string)
	}
	if context.GetInput(ivSecretKey) != nil {
		secretKey = context.GetInput(ivSecretKey).(string)
	}

	payload := ""
	switch p := context.GetInput(ivPayload).(type) {
	case string:
		payload = p
	case map[string]interface{}:
		b, err := json.Marshal(&p)
		if err != nil {
			log.Error(err)
			setTag("error", err.Error())

			return false, err
		}
		payload = string(b)
	}
	setTag("payload", payload)

	var config *aws.Config
	region := context.GetInput(ivRegion).(string)
	if accessKey != "" && secretKey != "" {
		config = aws.NewConfig().WithRegion(region).WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, ""))
	} else {
		config = aws.NewConfig().WithRegion(region)
	}
	aws := lambda.New(session.New(config))
	setTag("region", region)

	out, err := aws.Invoke(&lambda.InvokeInput{
		FunctionName: &arn,
		Payload:      []byte(payload)})

	if err != nil {
		log.Error(err)
		setTag("error", err.Error())

		return false, err
	}

	/*
		Removing this block, as it may be useful to get a response back and not an error... For the flow logic to do something specific and continue
		if *out.StatusCode != 200 {
			err := errors.New(*out.FunctionError)
			log.Error(err)

			return true, err
		}
	*/
	response := LambdaResponse{
		Status:  *out.StatusCode,
		Payload: out.Payload,
	}

	var output map[string]interface{}
	err = json.Unmarshal(out.Payload, &output)
	if err != nil {
		log.Error(err)
		setTag("error", err.Error())
	}

	log.Debugf("Lambda response: %s", string(response.Payload))
	setTag("response", string(out.Payload))
	setTag("responseStatus", *out.StatusCode)
	context.SetOutput(ovValue, response)
	context.SetOutput(ovResult, output)
	context.SetOutput(ovStatus, *out.StatusCode)

	return true, nil
}
