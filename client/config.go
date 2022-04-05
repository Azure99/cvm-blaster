package client

import (
	"encoding/json"
	"io/ioutil"
)

const configFile = "config.json"

var Conf Config

type Account struct {
	CVMEndpoint string
	CBSEndpoint string
	SecretId    string
	SecretKey   string
}

type Instance struct {
	Type           string
	Name           string
	LoginPassword  string
	DefaultImageId string
}

type Network struct {
	VPC         string
	SubNet      string
	SecureGroup string
	PrivateIp   string
}

type Service struct {
	Monitor    bool
	Security   bool
	Automation bool
}

type Region struct {
	Region string
	Zone   string
}

type Config struct {
	Account  Account
	Instance Instance
	Network  Network
	Service  Service
	Region   Region
}

func initConfig() {
	rawConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		Conf = getDefaultConfig()
		content, _ := json.MarshalIndent(Conf, "", "    ")
		if err := ioutil.WriteFile(configFile, content, 0644); err != nil {
			panic(err)
		}
		return
	}

	if err := json.Unmarshal(rawConfig, &Conf); err != nil {
		panic(err)
	}
}

func getDefaultConfig() Config {
	return Config{
		Account: Account{
			CVMEndpoint: "cvm.tencentcloudapi.com",
			CBSEndpoint: "cbs.tencentcloudapi.com",
			SecretId:    "AKIDCH8pKclpgJjvxxxxxxxxxxxxxxxxxxxx",
			SecretKey:   "VmobFfroX1ILxxxxxxxxxxxxxxxxxxxx",
		},
		Instance: Instance{
			Type:           "SA2.MEDIUM4",
			Name:           "cvm-blaster",
			LoginPassword:  "cvm-blaster123!A",
			DefaultImageId: "img-mmy6qctz",
		},
		Network: Network{
			VPC:         "vpc-xxxxxxxx",
			SubNet:      "subnet-xxxxxxxx",
			SecureGroup: "sg-xxxxxxxx",
			PrivateIp:   "",
		},
		Service: Service{
			Monitor:    true,
			Security:   true,
			Automation: true,
		},
		Region: Region{
			Region: "ap-beijing",
			Zone:   "ap-beijing-6",
		},
	}
}
