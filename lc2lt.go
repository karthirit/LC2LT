package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	ctx := context.Background()

	// Use AWS_PROFILE if set, or default
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "qa"
	}
	fmt.Println("Using AWS profile:", profile)

	// ---- EDIT THESE VARIABLES AS NEEDED ----
	launchConfigName := "intcloud-apigwqaperf-eks-qa-usw2-perf-graviton-worker20250218055912566700000001-1-31"
	launchTemplateName := launchConfigName + "-lt"
	awsRegion := "us-west-2" // Change this to your desired region
	// ----------------------------------------

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	asgClient := autoscaling.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)

	// Describe the Launch Configuration
	lcOutput, err := asgClient.DescribeLaunchConfigurations(ctx, &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []string{launchConfigName},
	})
	if err != nil {
		log.Fatalf("failed to describe launch configuration %s: %v", launchConfigName, err)
	}
	if len(lcOutput.LaunchConfigurations) == 0 {
		log.Fatalf("launch configuration %s not found", launchConfigName)
	}
	lc := lcOutput.LaunchConfigurations[0]

	// Prepare Launch Template Data
	ltData := ec2Types.RequestLaunchTemplateData{
		ImageId:      lc.ImageId,
		InstanceType: ec2Types.InstanceType(aws.ToString(lc.InstanceType)),
		KeyName:      lc.KeyName,
		Monitoring: &ec2Types.LaunchTemplatesMonitoringRequest{
			Enabled: aws.Bool(lc.InstanceMonitoring != nil && aws.ToBool(lc.InstanceMonitoring.Enabled)),
		},
		IamInstanceProfile: &ec2Types.LaunchTemplateIamInstanceProfileSpecificationRequest{
			Name: lc.IamInstanceProfile,
		},
		SecurityGroups:      lc.SecurityGroups,
		UserData:            lc.UserData,
		EbsOptimized:        lc.EbsOptimized,
		BlockDeviceMappings: []ec2Types.LaunchTemplateBlockDeviceMappingRequest{},
	}

	// Convert Block Device Mappings
	for _, bdm := range lc.BlockDeviceMappings {
		var ebs *ec2Types.LaunchTemplateEbsBlockDeviceRequest
		if bdm.Ebs != nil {
			ebs = &ec2Types.LaunchTemplateEbsBlockDeviceRequest{
				DeleteOnTermination: bdm.Ebs.DeleteOnTermination,
				Encrypted:           bdm.Ebs.Encrypted,
				Iops:                bdm.Ebs.Iops,
				SnapshotId:          bdm.Ebs.SnapshotId,
				Throughput:          bdm.Ebs.Throughput,
				VolumeSize:          bdm.Ebs.VolumeSize,
			}
			if bdm.Ebs.VolumeType != nil {
				ebs.VolumeType = ec2Types.VolumeType(aws.ToString(bdm.Ebs.VolumeType))
			}
		}
		ltData.BlockDeviceMappings = append(ltData.BlockDeviceMappings, ec2Types.LaunchTemplateBlockDeviceMappingRequest{
			DeviceName:  bdm.DeviceName,
			Ebs:         ebs,
			NoDevice:    nil, // Not directly mapped from LC
			VirtualName: bdm.VirtualName,
		})
	}

	// Create the Launch Template
	_, err = ec2Client.CreateLaunchTemplate(ctx, &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String(launchTemplateName),
		LaunchTemplateData: &ltData,
		VersionDescription: aws.String("Converted from launch configuration " + launchConfigName),
	})
	if err != nil {
		log.Fatalf("Failed to create launch template: %v", err)
	}
	fmt.Printf("Launch Template %s created successfully\n", launchTemplateName)
	fmt.Println("No ASG modification performed, only launch template created.")
}
