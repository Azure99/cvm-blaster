package client

import (
	"github.com/google/uuid"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var cvmClient *cvm.Client
var cbsClient *cbs.Client

func Init() {
	initConfig()

	credential := common.NewCredential(Conf.Account.SecretId, Conf.Account.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = Conf.Account.CVMEndpoint
	cvmClient, _ = cvm.NewClient(credential, Conf.Region.Region, cpf)

	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = Conf.Account.CBSEndpoint
	cbsClient, _ = cbs.NewClient(credential, Conf.Region.Region, cpf)
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
	var privateIp []string
	if Conf.Network.PrivateIp != "" {
		privateIp = append(privateIp, Conf.Network.PrivateIp)
	}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:              common.StringPtr(Conf.Network.VPC),
		SubnetId:           common.StringPtr(Conf.Network.SubNet),
		AsVpcGateway:       common.BoolPtr(false),
		Ipv6AddressCount:   common.Uint64Ptr(0),
		PrivateIpAddresses: common.StringPtrs(privateIp),
	}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(1),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceCount = common.Int64Ptr(1)
	request.InstanceName = common.StringPtr(Conf.Instance.Name)
	request.LoginSettings = &cvm.LoginSettings{}
	if imageId == "" {
		imageId = Conf.Instance.DefaultImageId
		request.LoginSettings.Password = common.StringPtr(Conf.Instance.LoginPassword)
	} else {
		request.LoginSettings.KeepImageLogin = common.StringPtr("TRUE")
	}
	request.ImageId = common.StringPtr(imageId)
	request.SecurityGroupIds = common.StringPtrs([]string{Conf.Network.SecureGroup})
	request.EnhancedService = &cvm.EnhancedService{
		SecurityService: &cvm.RunSecurityServiceEnabled{
			Enabled: common.BoolPtr(Conf.Service.Security),
		},
		MonitorService: &cvm.RunMonitorServiceEnabled{
			Enabled: common.BoolPtr(Conf.Service.Monitor),
		},
		AutomationService: &cvm.RunAutomationServiceEnabled{
			Enabled: common.BoolPtr(Conf.Service.Automation),
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
		return cvm.Instance{
			InstanceId: common.StringPtr(""),
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
		return cvm.Image{
			ImageId: common.StringPtr(""),
		}
	}
	return *response.Response.ImageSet[0]
}

func DeleteSnapshotAndImage(snapshotId string) {
	request := cbs.NewDeleteSnapshotsRequest()
	request.SnapshotIds = common.StringPtrs([]string{snapshotId})
	request.DeleteBindImages = common.BoolPtr(true)
	_, err := cbsClient.DeleteSnapshots(request)
	if err != nil {
		panic(err)
	}
}
