# Tencent CVM Blaster

自动化地（使用镜像）创建腾讯云服务器，并在销毁前创建备份镜像

## 前言

常见的云服务厂商往往提供三种付费方式：**包年包月**、**按量付费**、**竞价实例**。包年包月指一次性购买以月为基本单位的时长（**预付费**），按量付费指先使用再结算 按较短的时间间隔计费（**后付费**），竞价实例和按量付费相似，但当资源紧张或出价不足时 往往可能会释放实例。不过，竞价实例的价格往往是最低的，一般比按量付费**便宜数倍**，例如：腾讯云AMD SA3实例**8核16G**的价格仅为**0.2元/小时**。

价格如此便宜，有没有一种方法让我们做到随开随用，随关随停呢？会有两个问题：

- 数据无法持久化保存：删除实例时会删除硬盘，如果保留硬盘则需要长期付费
- 流量费过于高昂：按量付费的流量价格为0.8元/GB（上行）

腾讯云提供了免费50GB的硬盘快照额度，并可以基于此快照构建镜像，镜像可以用于创建实例，使用这50G快照额度，即可实现销毁实例时保存数据

腾讯云支持传统CVM和轻量应用服务器内网互联，同地域间享有5Gbps的免费带宽，轻量机的价格已经低至2C4G8M 48元/年，我们可以将所有上行流量转发至轻量机，这样即可实现流量免费

## 如何使用

Tencent CVM Blaster实现了自动化的创建/销毁流程，仅提供三个命令：start、stop、status

首次使用时，程序会创建配置文件`config.json`，请在文件中填写基本配置信息

- start：使用默认镜像或上次备份的镜像，创建实例
- stop：创建当前实例的快照、备份镜像，随后销毁实例
- status：查看当前实例的运行状态

注意：实例信息通过实例名Instance.Name区分，一份配置文件只支持维护一个实例

```json
{
    "Account": {
        "CVMEndpoint": "cvm.tencentcloudapi.com",
        "CBSEndpoint": "cbs.tencentcloudapi.com",
        "SecretId": "AKIDCH8pKclpgJjvxxxxxxxxxxxxxxxxxxxx", // 账号Id
        "SecretKey": "VmobFfroX1ILxxxxxxxxxxxxxxxxxxxx" // 私钥
    },
    "Instance": {
        "Type": "SA2.MEDIUM4", // 创建的机型
        "Name": "cvm-blaster", // 实例名、镜像前缀
        "LoginPassword": "cvm-blaster123!A", // 登录密码(首次创建时有效)
        "DefaultImageId": "img-mmy6qctz" // 默认镜像Id(Windows server 2019)
    },
    "Network": {
        "VPC": "vpc-xxxxxxxx", // VPC Id
        "SubNet": "subnet-xxxxxxxx", // 子网Id
        "SecureGroup": "sg-xxxxxxxx" // 安全组Id
    },
    "Region": {
        "Region": "ap-beijing", // 地区
        "Zone": "ap-beijing-6" // 可用区
    }
}
```

## 相关链接

- 创建子用户：https://console.cloud.tencent.com/cam

- 实例规格及ID：https://cloud.tencent.com/document/product/213/11518#SA2

- 公共镜像列表：https://console.cloud.tencent.com/cvm/image/index?rid=1&tab=PUBLIC_IMAGE&imageType=PUBLIC_IMAGE

- VPC：https://console.cloud.tencent.com/vpc/vpc

- 安全组：https://console.cloud.tencent.com/vpc/securitygroup

  
