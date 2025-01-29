package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/interpolation"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// applyCmd 解析 docker-compose.yml
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Parse docker-compose.yml and run",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			filePath = "docker-compose.yml"
		}

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Fatalf("docker-compose.yml does not exist: %s", filePath)
		}

		// 解析 Compose 文件
		services, err := parseComposeFile(filePath)
		if err != nil {
			log.Fatalf("Parse Compose failed: %v", err)
		}

		// 啟動容器（按照 depends_on 順序）
		started := make(map[string]bool)
		visiting := make(map[string]bool)
		for _, service := range services {
			if err := startContainerWithPodman(service, started, services, visiting); err != nil {
				log.Printf("Container %s launch failed: %v", service.Name, err)
			}
		}

		fmt.Println("✅ All containers started successfully")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("file", "f", "docker-compose.yml", "Compose YAML file")
}

// 解析 compose.yml 並返回服務配置
func parseComposeFile(filePath string) ([]types.ServiceConfig, error) {
	// 讀取 YAML 檔案
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 轉換為 map[string]interface{}
	var composeFile map[string]interface{}
	if err := yaml.Unmarshal(data, &composeFile); err != nil {
		return nil, err
	}

	// 取得環境變數
	envs := map[string]string{}
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envs[parts[0]] = parts[1]
		}
	}

	// 解析 Compose 配置
	cfg, err := loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{Filename: filePath, Config: composeFile},
		},
		Environment: envs,
	}, func(options *loader.Options) {
		options.SkipNormalization = true
		options.Interpolate = &interpolation.Options{}
		options.SkipValidation = false
	})
	if err != nil {
		return nil, err
	}

	return cfg.Services, nil
}

// 用 Podman 啟動容器，確保 `depends_on` 的容器先啟動
func startContainerWithPodman(service types.ServiceConfig, started map[string]bool, allServices []types.ServiceConfig, visiting map[string]bool) error {
	// 先處理 `depends_on`
	for depName := range service.DependsOn {
		// 避免循環依賴（如果 `visiting` 中已經有這個容器）
		if visiting[depName] {
			return fmt.Errorf("detected circular dependency: %s depends on %s, but it is already in the startup process", service.Name, depName)
		}

		// 如果 `depends_on` 容器還沒啟動，則先啟動
		if !started[depName] {
			// 在所有服務中尋找 `depName`
			found := false
			for _, s := range allServices {
				if s.Name == depName {
					fmt.Printf("Wait for dependency %s launching...\n", depName)

					// 標記當前容器，避免循環依賴s
					visiting[service.Name] = true

					// 啟動 `depends_on` 容器
					if err := startContainerWithPodman(s, started, allServices, visiting); err != nil {
						fmt.Printf("Cannot launch dependent container %s: %v\n", depName, err)
						return err
					}

					// 依賴已啟動，取消標記
					delete(visiting, service.Name)
					found = true
					break
				}
			}

			// 如果 `depends_on` 服務在 `docker-compose.yml` 中找不到
			if !found {
				return fmt.Errorf("dependent container %s is not dependent in `docker-compose.yml`", depName)
			}
		}
	}

	// 如果容器已經存在，先刪除
	if isContainerRunning(service.Name) {
		fmt.Printf("Container %s already exists, removing it...\n", service.Name)
		exec.Command("podman", "stop", service.Name).Run()
		exec.Command("podman", "rm", service.Name).Run()
	}

	// 構建 Podman run 命令
	args := []string{"run", "-d", "--name", service.Name}

	// 設置環境變數
	for key, value := range service.Environment {
		if value != nil {
			args = append(args, "-e", fmt.Sprintf("%s=%s", key, *value))
		} else {
			args = append(args, "-e", key)
		}
	}

	// 設置端口
	for _, port := range service.Ports {
		args = append(args, "-p", fmt.Sprintf("%s:%s", port.Published, strconv.Itoa(int(port.Target))))
	}

	// 設置網路
	if service.NetworkMode != "" {
		args = append(args, "--network", service.NetworkMode)
	}

	// 設置 volumes
	for _, volume := range service.Volumes {
		args = append(args, "-v", fmt.Sprintf("%s:%s", volume.Source, volume.Target))
	}

	// 設置 restart 策略
	if service.Restart != "" {
		args = append(args, "--restart", service.Restart)
	}

	// 設置 image
	args = append(args, service.Image)

	// 執行 Podman 命令
	cmd := exec.Command("podman", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution error: %s", stderr.String())
	}

	fmt.Printf("Container %s started successfully\n", service.Name)
	started[service.Name] = true
	return nil
}

// 檢查容器是否已存在
func isContainerRunning(containerName string) bool {
	cmd := exec.Command("podman", "ps", "-a", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Unable to check containers: %v", err)
		return false
	}

	containerList := strings.Split(string(output), "\n")
	for _, name := range containerList {
		if name == containerName {
			return true
		}
	}
	return false
}
