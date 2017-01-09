package subscriptions

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/satori/go.uuid"
	"os"
	"time"
)

var svc *dynamodb.DynamoDB

func init() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc = dynamodb.New(sess)
}

var subscriptionTableName = os.Getenv("SUBSCRIPTIONS_DYNAMODB_TABLE")

const columnId = "ID"
const columnName = "Name"
const columnEmail = "Email"
const columnSize = "Size"
const columnStripeToken = "StripeToken"
const columnSignupDate = "SignupDate"

const subscriptionTablePrimaryKey = columnId

const fmtSignupDate = "2006-01-02T15:04:05Z"

type Subscription struct {
	ID          string
	Name        string
	Size        string
	Email       string
	StripeToken string
}

type FetchAllSubscriptionsResult struct {
	LastEvaluatedKey *string
	Subscriptions    []*Subscription
	HasMore          bool
}

func FetchAllSubscriptions(limit int, lastKey *string) (*FetchAllSubscriptionsResult, error) {
	var limit64 int64 = int64(limit)
	var input = &dynamodb.ScanInput{
		TableName: aws.String(subscriptionTableName),
		AttributesToGet: []*string{
			aws.String(columnId),
			aws.String(columnName),
			aws.String(columnSize),
			aws.String(columnEmail),
			aws.String(columnStripeToken),
		},
		Limit: &limit64,
	}

	if lastKey != nil {
		var startKey = make(map[string]*dynamodb.AttributeValue)
		startKey[subscriptionTablePrimaryKey] = &dynamodb.AttributeValue{
			S: aws.String(*lastKey)}
		input.ExclusiveStartKey = startKey
	}

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	resp, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	lastEvaluatedKey := resp.LastEvaluatedKey[subscriptionTablePrimaryKey]

	result := &FetchAllSubscriptionsResult{
		Subscriptions: parseSubscriptions(resp.Items)}

	if lastEvaluatedKey != nil {
		s := getString(lastEvaluatedKey)
		result.LastEvaluatedKey = &s
		result.HasMore = true
	}

	return result, nil
}

func FetchSubscription(id string) (*Subscription, error) {
	var key = make(map[string]*dynamodb.AttributeValue)
	key[subscriptionTablePrimaryKey] = &dynamodb.AttributeValue{
		S: aws.String(id)}

	var input = &dynamodb.GetItemInput{
		AttributesToGet: []*string{
			aws.String(columnId),
			aws.String(columnName),
			aws.String(columnSize),
			aws.String(columnEmail),
			aws.String(columnStripeToken),
		},
		Key:       key,
		TableName: &subscriptionTableName}

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	resp, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	return parseSubscription(resp.Item), nil
}

func CreateSubscription(name string, email string, size string) (*Subscription, error) {
	itemId := uuid.NewV4().String()
	signupDate := time.Now().UTC().Format(fmtSignupDate)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			columnId:         {S: &itemId},
			columnName:       {S: &name},
			columnEmail:      {S: &email},
			columnSize:       {S: &size},
			columnSignupDate: {S: &signupDate},
		},
		TableName: aws.String(subscriptionTableName),
	}
	_, err := svc.PutItem(input)
	if err != nil {
		return nil, err
	}

	return &Subscription{
		ID:    itemId,
		Name:  name,
		Email: email,
		Size:  size,
	}, nil
}

func parseSubscriptions(items []map[string]*dynamodb.AttributeValue) []*Subscription {
	out := make([]*Subscription, len(items))

	for i := 0; i < len(items); i++ {
		out[i] = parseSubscription(items[i])
	}

	return out
}

func parseSubscription(item map[string]*dynamodb.AttributeValue) *Subscription {
	return &Subscription{
		ID:          getString(item[columnId]),
		Name:        getString(item[columnName]),
		Size:        getString(item[columnSize]),
		Email:       getString(item[columnEmail]),
		StripeToken: getString(item[columnStripeToken]),
	}
}

func getString(sVal *dynamodb.AttributeValue) string {
	if sVal == nil {
		return ""
	}

	return *sVal.S
}
