package lib

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

//ListASGWithSpecificTag will return list of asg arn carrying specific tag
func ListASGWithSpecificTag(client *autoscaling.Client, tagname, tagvalue string) ([]*string, error) {
	var output []*string
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	for {
		result, err := client.DescribeAutoScalingGroups(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, item := range result.AutoScalingGroups {
			for _, tag := range item.Tags {
				if *tag.Key == tagname && *tag.Value == tagvalue {
					output = append(output, item.AutoScalingGroupName)
				}
			}
		}
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
	}
	return output, nil
}

// ChangeASGCapacity will set a new desired capacity
func ChangeASGCapacity(client autoscaling.Client, name *string, number int32) error {
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: name,
		DesiredCapacity:      &number,
	}
	_, err := client.SetDesiredCapacity(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

// CheckASGCapacity will return the number of instances currently belonging to a ASG
func CheckASGCapacity(client *autoscaling.Client, name string) (int, error) {
	names := []string{name}
	input := &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: names}
	result, err := client.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		return 0, err
	}
	if len(result.AutoScalingGroups) != 1 {
		return 0, errors.New("Too many ASGs have been found")
	}
	return len(result.AutoScalingGroups[0].Instances), nil
}
