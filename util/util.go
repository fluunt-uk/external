package util

import (
	"crypto/rand"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"gitlab.com/projectreferral/util/client/rabbitmq"
	"gitlab.com/projectreferral/util/client/s3"
	"gitlab.com/projectreferral/util/client/s3/models"
	dynamo_lib "gitlab.com/projectreferral/util/pkg/dynamodb"
	"log"
)

type ServiceConfigs struct {
	Environment  	string
	Region       	string
	Table        	string
	SearchParam  	string
	GenericModel 	interface{}
	BrokerUrl    	string
	Port         	string
}

//internal specific configs are loaded at runtime
func (sc *ServiceConfigs) LoadEnvConfigs() {

	log.Printf("Environment: %s\n", sc.Environment)
	log.Printf("Running on %s\n", sc.Port)
}

//DynamoDB configs
func (sc *ServiceConfigs) LoadDynamoDBConfigs() *dynamo_lib.Wrapper {

	switch sc.Environment {
	case "UAT":
		sc.Table = "uat-" + sc.Table
	case "PROD":
		sc.Table = "prod-" + sc.Table
	default:
		sc.Table = "dev-" + sc.Table
	}

	dynamoDBInstance := &dynamo_lib.Wrapper{
		GenericModel: &sc.GenericModel,
		SearchParam:  &sc.SearchParam,
		Table:        &sc.Table,
		Credentials:  sc.generateCredentials(),
		Region:       &sc.Region,
	}
	return dynamoDBInstance
}

//RabbitMQ configs
func (sc *ServiceConfigs) LoadRabbitMQConfigs() *rabbitmq.DefaultQueueClient {

	client := &rabbitmq.DefaultQueueClient{}
	client.SetupURL(sc.BrokerUrl)

	return client
}

//S3 configs
func (sc *ServiceConfigs) LoadS3BucketConfigs(cfg *models.S3Configs) *s3.DefaultBucketClient {

	client := &s3.DefaultBucketClient{}
	client.SetConfigs(cfg)
	return client
}

//Create a single instance of DynamoDB connection
func (sc *ServiceConfigs) generateCredentials() *credentials.Credentials {

	c := credentials.NewSharedCredentials("", "default")

	return c
}

func NewUUID() (string,error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}