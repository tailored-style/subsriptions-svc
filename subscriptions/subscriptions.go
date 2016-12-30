package subscriptions

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/satori/go.uuid"
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

const subscriptionTableName = "tailored.monthly.subscriptions"
const subscriptionTablePrimaryKey = "ID"

type Subscription struct {
	ID    string
	Name  string
	Size  string
	Email string
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
			aws.String("ID"),
			aws.String("Name"),
			aws.String("Size"),
			aws.String("Email"),
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

func CreateSubscription(name string, email string, size string) (*Subscription, error) {
	itemId := uuid.NewV4().String()
	signupDate := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID":    {S: &itemId},
			"Name":  {S: &name},
			"Email": {S: &email},
			"Size":  {S: &size},
			"SignupDate": {S: &signupDate},
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
		Size:  size}, nil
}

func parseSubscriptions(items []map[string]*dynamodb.AttributeValue) []*Subscription {
	out := make([]*Subscription, len(items))

	for i := 0; i < len(items); i++ {
		item := items[i]
		out[i] = &Subscription{
			ID:    getString(item["ID"]),
			Name:  getString(item["Name"]),
			Size:  getString(item["Size"]),
			Email: getString(item["Email"])}
	}

	return out
}

func getString(sVal *dynamodb.AttributeValue) string {
	return *sVal.S
}

