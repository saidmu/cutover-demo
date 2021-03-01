package lib

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// ListDBInstancesWithSpecificTag will return a list of db identifier
func ListDBInstancesWithSpecificTag(client *rds.Client, tagname, tagvalue string) ([]*string, error) {
	var output []*string
	input := &rds.DescribeDBInstancesInput{}
	for {
		result, err := client.DescribeDBInstances(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, item := range result.DBInstances {
			for _, tag := range item.TagList {
				if *tag.Key == tagname && *tag.Value == tagvalue {
					output = append(output, item.DBInstanceIdentifier)
				}
			}
		}
		if result.Marker == nil {
			break
		}
		input.Marker = result.Marker
	}
	return output, nil
}

// CheckIsReplica will check if the instance is a replica
func CheckIsReplica(client *rds.Client, identifier *string) (bool, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: identifier,
	}
	result, err := client.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		return false, err
	}
	if len(result.DBInstances) != 1 {
		return false, errors.New("To many instances have been found")
	}
	if len(result.DBInstances[0].StatusInfos) == 0 {
		return false, nil
	}
	return true, nil
}

// PromoteReplicaToPrimary will Promote a backup instqnce to primary
func PromoteReplicaToPrimary(client *rds.Client, identifier *string) error {
	input := &rds.PromoteReadReplicaInput{DBInstanceIdentifier: identifier}
	_, err := client.PromoteReadReplica(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
