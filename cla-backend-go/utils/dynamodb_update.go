package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoUpdateExpression helps build update expression
type DynamoUpdateExpression struct {
	Expression                string
	ExpressionAttributeNames  map[string]*string
	ExpressionAttributeValues map[string]*dynamodb.AttributeValue
}

// NewDynamoUpdateExpression create DynamoUpdateExpression builder
func NewDynamoUpdateExpression() *DynamoUpdateExpression {
	return &DynamoUpdateExpression{
		Expression:                "",
		ExpressionAttributeNames:  make(map[string]*string),
		ExpressionAttributeValues: make(map[string]*dynamodb.AttributeValue),
	}
}

// AddUpdateExpression the update expression if condition satisfied
func (d *DynamoUpdateExpression) AddUpdateExpression(columnUpdateExp string, condition bool) {
	if condition {
		if d.Expression == "" {
			d.Expression = "SET " + columnUpdateExp
		} else {
			d.Expression = d.Expression + ", " + columnUpdateExp
		}
	}
}

// AddAttributeName adds the attribute if condition satisfied
func (d *DynamoUpdateExpression) AddAttributeName(name, columName string, condition bool) {
	if condition {
		d.ExpressionAttributeNames[name] = aws.String(columName)
	}
}

// AddAttributeValue adds the attribute value if condition satisfied
func (d *DynamoUpdateExpression) AddAttributeValue(name string, val *dynamodb.AttributeValue, condition bool) {
	if condition {
		d.ExpressionAttributeValues[name] = val
	}
}
