package client

import (
	"github.com/google/uuid"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var cvmClient *cvm.Client

func Init() {
	initConfig()

	credential := common.NewCredential(Conf.Account.SecretId, Conf.Account.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = Conf.Account.Endpoint
	client, _ := cvm.NewClient(credential, Conf.Region.Region, cpf)

	cvmClient = client
}

func CreateInstance(imageId string) string {
	request := cvm.NewRunInstancesRequest()
	request.InstanceChargeType = common.StringPtr("SPOTPAID")
	request.Placement = &cvm.Placement{
		Zone:      common.StringPtr(Conf.Region.Zone),
		ProjectId: common.Int64Ptr(0),
	}
	request.InstanceType = common.StringPtr(Conf.Instance.Type)
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: common.StringPtr("CLOUD_PREMIUM"),
		DiskSize: common.Int64Ptr(50),
	}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:            common.StringPtr(Conf.Network.VPC),
		SubnetId:         common.StringPtr(Conf.Network.SubNet),
		AsVpcGateway:     common.BoolPtr(false),
		Ipv6AddressCount: common.Uint64Ptr(0),
	}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(1),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceCount = common.Int64Ptr(1)
	request.InstanceName = common.StringPtr(Conf.Instance.Name)
	if imageId == "" {
		request.LoginSettings = &cvm.LoginSettings{
			Password: common.StringPtr(Conf.Instance.LoginPassword),
		}
		imageId = Conf.Instance.DefaultImageId
	}
	request.ImageId = common.StringPtr(imageId)
	request.SecurityGroupIds = common.StringPtrs([]string{Conf.Network.SecureGroup})
	request.EnhancedService = &cvm.EnhancedService{
		SecurityService: &cvm.RunSecurityServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		MonitorService: &cvm.RunMonitorServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		AutomationService: &cvm.RunAutomationServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
	}
	request.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{
		SpotOptions: &cvm.SpotMarketOptions{
			MaxPrice: common.StringPtr("1000"),
		},
	}

	response, err := cvmClient.RunInstances(request)
	if err != nil {
		panic(err)
	}

	return *response.Response.InstanceIdSet[0]
}

func GetInstance() cvm.Instance {
	request := cvm.NewDescribeInstancesRequest()
	request.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("zone"),
			Values: common.StringPtrs([]string{Conf.Region.Zone}),
		},
		{
			Name:   common.StringPtr("instance-name"),
			Values: common.StringPtrs([]string{Conf.Instance.Name}),
		},
	}

	response, err := cvmClient.DescribeInstances(request)
	if err != nil {
		panic(err)
	}

	if len(response.Response.InstanceSet) == 0 {
		instanceId := ""
		return cvm.Instance{
			InstanceId: &instanceId,
		}
	}

	return *response.Response.InstanceSet[0]
}

func DeleteInstance(instanceId string) {
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{instanceId})

	_, err := cvmClient.TerminateInstances(request)
	if err != nil {
		panic(err)
	}
}

func CreateImage(instanceId string) string {
	request := cvm.NewCreateImageRequest()
	request.InstanceId = common.StringPtr(instanceId)
	request.ImageName = common.StringPtr(Conf.Instance.Name + uuid.New().String())
	request.ForcePoweroff = common.StringPtr("TRUE")
	response, err := cvmClient.CreateImage(request)
	if err != nil {
		panic(err)
	}

	return *response.Response.ImageId
}

func GetImage() cvm.Image {
	request := cvm.NewDescribeImagesRequest()
	request.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("image-name"),
			Values: common.StringPtrs([]string{Conf.Instance.Name}),
		},
	}

	response, err := cvmClient.DescribeImages(request)
	if err != nil {
		panic(err)
	}

	if len(response.Response.ImageSet) == 0 {
		imageId := ""
		return cvm.Image{
			ImageId: &imageId,
		}
	}
	return *response.Response.ImageSet[0]
}

func DeleteImage(id string) {
	request := cvm.NewDeleteImagesRequest()
	request.ImageIds = common.StringPtrs([]string{id})
	_, err := cvmClient.DeleteImages(request)
	if err != nil {
		panic(err)
	}
}
