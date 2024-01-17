package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	"hetzner-dns/internal/config"
	"hetzner-dns/internal/hetzner"
	public_ip "hetzner-dns/internal/public-ip"
	"hetzner-dns/internal/vault"
	"hetzner-dns/internal/version"
	"hetzner-dns/pkg/logging"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Run() int {
	logger := logging.GetLogger()
	logger.Info("Application logger initialized.")
	logger.Info(""+
		"Start application...",
		logger.String("version", version.Version),
		logger.String("build_time", version.BuildTime),
		logger.String("commit", version.Commit),
	)

	appConfig := config.GetAppConfig()
	zones, err := appConfig.GetZones()
	if err != nil {
		logger.Error(err.Error())
	}

	if appConfig.UseVault {
		token, err := readTokenFromVault(logger, appConfig)
		if err != nil {
			logger.Error(err.Error())
			return 1
		}
		appConfig.HetznerToken = token
	}

	shed := gocron.NewScheduler(time.Local)

	go func() {
		every, err := strconv.Atoi(appConfig.UpdateInterval)
		if err != nil {
			logger.Error(err.Error())
			every = 60
		}
		_, err = shed.Every(every).Second().Do(func() {
			ipAddr := public_ip.GetPublicIpAddr()
			hetznerClient := hetzner.NewHetznerClient(appConfig.HetznerToken)

			resp, err := hetznerClient.GetZones()
			if err != nil {
				logger.Error(err.Error())
			}

			for _, configZone := range zones.Zones {
				for _, dnsZone := range resp.Zones {
					if dnsZone.Name == configZone.Name {
						dnsRecords, err := hetznerClient.GetRecords(dnsZone.ID)
						if err != nil {
							fmt.Println("Failure : ", err)
						}

						for _, cfgRecord := range configZone.Records {
							flag := false
							for _, dnsRecord := range dnsRecords.Records {
								if dnsRecord.Name == cfgRecord.Name {
									flag = true
									value := cfgRecord.Value
									if net.ParseIP(value) == nil {
										value = ipAddr
									}
									record := hetzner.Record{
										Value:  value,
										TTL:    cfgRecord.TTL,
										Type:   cfgRecord.Type,
										Name:   cfgRecord.Name,
										ZoneID: dnsZone.ID,
									}

									payload := new(bytes.Buffer)
									err := json.NewEncoder(payload).Encode(record)
									if err != nil {
										fmt.Println("Failure : ", err)
										break
									}
									err = hetznerClient.UpdateRecord(dnsRecord.ID, payload.Bytes())
									if err != nil {
										fmt.Println("Failure : ", err)
										break
									}
									logger.Info("Record updated",
										logger.String("zone", dnsZone.Name),
										logger.String("name", cfgRecord.Name),
										logger.String("type", cfgRecord.Type),
										logger.String("value", cfgRecord.Value),
									)
								}
							}
							if !flag {
								value := cfgRecord.Value
								if net.ParseIP(value) == nil {
									value = ipAddr
								}
								record := hetzner.Record{
									Value:  value,
									TTL:    cfgRecord.TTL,
									Type:   cfgRecord.Type,
									Name:   cfgRecord.Name,
									ZoneID: dnsZone.ID,
								}

								payload := new(bytes.Buffer)
								err := json.NewEncoder(payload).Encode(record)
								if err != nil {
									fmt.Println("Failure : ", err)
									continue
								}
								err = hetznerClient.CreateRecord(payload.Bytes())
								if err != nil {
									fmt.Println("Failure : ", err)
									continue
								}
								logger.Info("Record created",
									logger.String("zone", dnsZone.Name),
									logger.String("name", cfgRecord.Name),
									logger.String("type", cfgRecord.Type),
									logger.String("value", cfgRecord.Value),
								)
							}
						}

					}
				}
			}
		})

		if err != nil {
			logger.Error(err.Error())
		}
		shed.StartBlocking()
	}()

	shutdownChan := make(chan os.Signal, 1)
	hupChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGABRT, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(hupChan, syscall.SIGHUP)

	for {
		select {
		case <-hupChan:
			appConfig = config.GetAppConfig()
			if appConfig.UseVault {
				token, err := readTokenFromVault(logger, appConfig)
				if err != nil {
					logger.Error(err.Error())
					return 1
				}
				appConfig.HetznerToken = token
			}
			logger.Info("Reload config from env success")
		case <-shutdownChan:
			logger.Info("Shutdown Application...")
			logger.Info("Application successful shutdown")
			return 0
		}
	}
}

func readTokenFromVault(logger logging.Logger, appConfig config.AppConfig) (string, error) {
	logger.Info("Use Vault as credential storage for authentication")
	vaultClient, err := vault.NewClient(appConfig, &logger)
	if err != nil {
		logger.Error(err.Error())
	}

	data, err := vaultClient.GetClient().KVv2(appConfig.VaultSecretStoreName).Get(
		context.Background(),
		appConfig.VaultSecretStorePath,
	)
	if err != nil {
		logger.Error(
			fmt.Sprintf("error occurred when get credentials from Vault for target: %s",
				appConfig.VaultSecretStorePath),
		)
		logger.Error(err.Error())
		return "", err
	}

	return data.Data["token"].(string), nil
}
