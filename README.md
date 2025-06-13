# Launch Configuration to Launch Template Converter

This tool converts AWS Auto Scaling Launch Configurations to EC2 Launch Templates.

## Prerequisites

- Go 1.20 or later
- Docker
- AWS credentials with appropriate permissions

## Building

### Using Go
```bash
go build -o lc2lt .
```

### Using Docker
```bash
docker build -t lc2lt .
```

## Running

### Using Go
```bash
./lc2lt
```

### Using Docker
```bash
docker run --rm -v ~/.aws:/root/.aws -e AWS_PROFILE=qa -e AWS_REGION=us-west-2 lc2lt:latest
```

## Required AWS Permissions

The tool requires the following AWS permissions:
- `autoscaling:DescribeLaunchConfigurations`
- `ec2:CreateLaunchTemplate`

## Configuration

Edit the following variables in `lc2lt.go`:
- `launchConfigName`: Name of the Launch Configuration to convert
- `launchTemplateName`: Name for the new Launch Template
- `awsRegion`: AWS region to use

## Kubernetes Deployment

The repository includes Kubernetes deployment files:
- `lc2lt-deployment.yaml`: Deployment configuration
- `aws-config.yaml`: AWS configuration ConfigMap # LC2LT
