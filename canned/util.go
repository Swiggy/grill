package canned

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func skipReaper() bool {
	val, _ := strconv.ParseBool(os.Getenv("TESTCONTAINERS_RYUK_DISABLED"))
	return val
}

func getEnvString(variable string, defaultValue string) string {
	val := os.Getenv(variable)
	if val == "" {
		return defaultValue
	}
	return val
}

func getAWSConfig() (string, string, string) {
	accessKey := getEnvString("AWS_ACCESS_KEY_ID", "awsaccesskey")
	secretKey := getEnvString("AWS_SECRET_ACCESS_KEY", "awssecretkey")
	region := getEnvString("AWS_REGION", "ap-southeast-1")

	return accessKey, secretKey, region
}

func getBasicAuth() string {
	username := os.Getenv("CONTAINER_REGISTRY_USERNAME")
	password := os.Getenv("CONTAINER_REGISTRY_PASSWORD")
	if len(username) == 0 || len(password) == 0 {
		return ""
	}

	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encoded, _ := json.Marshal(auth)

	return base64.URLEncoding.EncodeToString(encoded)
}

func getConsecutiveFreePorts(count, startPort, endPort int) ([]int, error) {
	// 10 tries to find a free consecutive ports
	for i := 0; i < 10; i++ {
		port := rand.Intn(endPort-startPort) + startPort
		if ports, ok := checkConsecutiveFreePorts(port, count); ok {
			return ports, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("%d free ports are not available in range %d-%d", count, startPort, endPort))
}

func checkConsecutiveFreePorts(startingPort int, count int) ([]int, bool) {
	var freePorts []int
	for i := 0; i < count; i++ {
		port := startingPort + i
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return freePorts, false
		}
		if ln != nil {
			ln.Close()
			freePorts = append(freePorts, port)
		}
	}
	return freePorts, true
}
