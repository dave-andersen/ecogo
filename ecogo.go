package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tess1o/go-ecoflow"
)

type Config struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	SerialNumber string `json:"serialNumber,omitempty"`
}

func loadConfig() Config {
	var cfg Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Warn("Error getting home directory, proceeding without config file values", "error", err)
	} else {
		configPath := filepath.Join(homeDir, ".ecoflow")
		configData, err := os.ReadFile(configPath)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Warn("Error reading config file, will rely on environment variables or defaults", "path", configPath, "error", err)
			}
		} else {
			err = json.Unmarshal(configData, &cfg)
			if err != nil {
				slog.Warn("Error parsing config file, will rely on environment variables or defaults", "path", configPath, "error", err)
				cfg = Config{}
			}
		}
	}

	// Override with environment variables if they are set
	if envAccessKey := os.Getenv("ACCESS_KEY"); envAccessKey != "" {
		cfg.AccessKey = envAccessKey
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		cfg.SecretKey = envSecretKey
	}
	if envSerialNumber := os.Getenv("SERIAL_NUMBER"); envSerialNumber != "" {
		cfg.SerialNumber = envSerialNumber
	}
	return cfg
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	config := loadConfig()

	if len(os.Args) < 2 {
		slog.Error("Mode argument is required (e.g., list, check, selfpow, or other for default parameter setting)")
		os.Exit(1)
	}
	mode := os.Args[1]

	if config.AccessKey == "" || config.SecretKey == "" {
		slog.Error("AccessKey and SecretKey are mandatory. Set them in ~/.ecoflow (JSON format) or via ACCESS_KEY and SECRET_KEY environment variables.")
		return
	}

	client := ecoflow.NewEcoflowClient(config.AccessKey, config.SecretKey)

	if mode == "list" {
		dl, err := client.GetDeviceList(context.Background())
		if err != nil {
			slog.Error("Error getting device list", "error", err)
		} else {
			slog.Info("Device list", "devices", dl)
		}
		return
	}

	if config.SerialNumber == "" {
		slog.Error("SerialNumber is mandatory for modes other than 'list'. Set it in ~/.ecoflow (JSON format) or via SERIAL_NUMBER environment variable.")
		return
	}

	if mode == "check" {
		dp, err := client.GetDeviceAllParameters(context.Background(), config.SerialNumber)
		if err != nil {
			slog.Error("Error getting device parameters", "error", err)
		} else {
			// print out device parameters in JSON
			dpJSON, err := json.MarshalIndent(dp, "", "  ")
			if err != nil {
				slog.Error("Error marshaling device parameters to JSON", "error", err)
			} else {
				os.Stdout.Write(dpJSON)
			}
		}
		return
	}

	selfpow := false
	backup_reserve := 45 // Default values
	if mode == "selfpow" {
		selfpow = true
		backup_reserve = 70
	}

	var req map[string]interface{} = make(map[string]interface{})
	req["sn"] = config.SerialNumber
	req["cmdId"] = 17
	req["dirDest"] = 1
	req["dirSrc"] = 1
	req["cmdFunc"] = 254
	req["dest"] = 2
	req["needAck"] = true

	var params map[string]interface{} = make(map[string]interface{})
	params["cfgEnergyBackup"] = map[string]interface{}{
		"energyBackupStartSoc": backup_reserve,
		"energyBackupEn":       true,
	}
	params["cfgEnergyStrategyOperateMode"] = map[string]interface{}{
		"operateSelfPoweredOpen": selfpow,
		"operateTouModeOpen":     !selfpow,
	}

	req["params"] = params

	resp, err := client.SetDeviceParameter(context.Background(), req)
	if err != nil {
		slog.Error("Error setting device parameter", "error", err)
		return
	} else {
		slog.Info("Set device parameter response", "response", resp)
	}
}
