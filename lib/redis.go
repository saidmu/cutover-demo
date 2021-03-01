package lib

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

// GetAllSecondaryReplicationGroup func
func GetAllSecondaryReplicationGroup(client *elasticache.Client) (map[string]string, error) {
	output := make(map[string]string)
	input := &elasticache.DescribeReplicationGroupsInput{}
	for {
		result, err := client.DescribeReplicationGroups(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, item := range result.ReplicationGroups {
			if item.GlobalReplicationGroupInfo == nil {
				continue
			}
			if *item.GlobalReplicationGroupInfo.GlobalReplicationGroupMemberRole != "SECONDARY" {
				continue
			}
			output[*item.ReplicationGroupId] = *item.GlobalReplicationGroupInfo.GlobalReplicationGroupId
		}
		if result.Marker == nil {
			break
		}
		input.Marker = result.Marker
	}
	return output, nil
}

// PromteToPrimary func
func PromteToPrimary(client *elasticache.Client, region string, data map[string]string) error {
	for k, v := range data {
		input := &elasticache.FailoverGlobalReplicationGroupInput{
			PrimaryRegion:             &region,
			PrimaryReplicationGroupId: aws.String(k),
			GlobalReplicationGroupId:  aws.String(v),
		}
		_, err := client.FailoverGlobalReplicationGroup(context.TODO(), input)
		if err != nil {
			return err
		}
	}
	return nil
}
