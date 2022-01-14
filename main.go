package main

import (
	"os"
	"strconv"
	"time"

	"cvm-blaster/client"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func main() {
	client.Init()

	var command string
	if len(os.Args) == 2 {
		command = os.Args[1]
	}

	switch command {
	case "start":
		start()
	case "stop":
		stop()
	case "status":
		status()
	default:
		help()
	}
}

func start() {
	existId := getInstanceId()
	if existId != "" {
		println("Instance already created: " + existId)
		return
	}

	imageId := getImageId()
	println("Backup Image id: " + imageId)
	if imageId != "" && getImageStatus() != "NORMAL" {
		println("Image not available")
		return
	}

	println("Create instance")
	instanceId := createInstance(imageId)

	for {
		time.Sleep(3 * time.Second)
		status := getInstanceStatus()
		if status == "RUNNING" {
			break
		}
		println(instanceId + ": " + status)
	}
	println("Instance created: " + instanceId)
	printInstanceInfo(getInstance())
}

func stop() {
	instanceId := getInstanceId()
	if instanceId == "" {
		println("Instance not exist")
		return
	}

	oldImageId := getImageId()

	println("Create Backup Image")
	imageId := createImage(instanceId)

	for {
		time.Sleep(5 * time.Second)
		status := getImageStatus()
		if status == "NORMAL" {
			break
		}
		println(imageId + ": " + status)
	}
	println("Image created: " + imageId)

	if oldImageId != "" {
		println("Delete old image: " + oldImageId)
		deleteImage(oldImageId)
	}

	println("Delete instance: " + instanceId)
	deleteInstance(instanceId)
}

func status() {
	instance := getInstance()
	if *instance.InstanceId == "" {
		println("Instance not running")
	} else {
		printInstanceInfo(instance)
	}
}

func help() {
	println("Command list")
	println("start:\tcreate cvm instance (use backup image)")
	println("stop:\tcreate backup image then destroy cvm instance")
	println("status:\tshow cvm instance status")
}

func createInstance(imageId string) string {
	return client.CreateInstance(imageId)
}

func getInstanceId() string {
	return *client.GetInstance().InstanceId
}

func getInstanceStatus() string {
	return *client.GetInstance().InstanceState
}

func getInstance() cvm.Instance {
	return client.GetInstance()
}

func deleteInstance(instanceId string) {
	client.DeleteInstance(instanceId)
}

func getImageId() string {
	return *client.GetImage().ImageId
}

func getImageStatus() string {
	return *client.GetImage().ImageState
}

func createImage(instanceId string) string {
	return client.CreateImage(instanceId)
}

func deleteImage(imageId string) {
	client.DeleteImage(imageId)
}

func printInstanceInfo(instance cvm.Instance) {
	println("Id:\t" + *instance.InstanceId)
	println("CPU:\t" + strconv.Itoa(int(*instance.CPU)))
	println("Memory:\t" + strconv.Itoa(int(*instance.Memory)))
	println("OS:\t" + *instance.OsName)
	println("Public ip:\t" + *instance.PublicIpAddresses[0])
	println("Private ip:\t" + *instance.PrivateIpAddresses[0])
}
