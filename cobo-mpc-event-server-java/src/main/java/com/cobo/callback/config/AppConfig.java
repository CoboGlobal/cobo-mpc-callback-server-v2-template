package com.cobo.callback.config;

import java.io.File;
import java.io.IOException;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Data
@Slf4j
@NoArgsConstructor
public class AppConfig {
    private static final String DEFAULT_CONFIG_YAML = "configs/callback-server-config.yaml";

    private String serviceName = "callback-server";
    private String endpoint = "0.0.0.0:11020";
    private int tokenExpireMinutes = 2;
    private String clientPublicKeyPath = "configs/tss-node-callback-pub.key";
    private String servicePrivateKeyPath = "configs/callback-server-pri.pem";
    private boolean enableDebug = false;

    public static AppConfig loadConfig(String[] args) {
        String configPath = parseConfigPath(args);
        return loadYamlConfig(configPath);
    }

    private static String parseConfigPath(String[] args) {
        for (int i = 0; i < args.length - 1; i++) {
            if (args[i].equals("-c") || args[i].equals("--config")) {
                return args[i + 1];
            }
        }
        return DEFAULT_CONFIG_YAML;
    }

    private static AppConfig loadYamlConfig(String configPath) {
        try {
            File configFile = new File(configPath);
            if (!configFile.exists()) {
                log.warn("Config file not found: {}, using default configuration", configPath);
                return new AppConfig();
            }

            ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
            JsonNode root = mapper.readTree(configFile);
            JsonNode server = root.get("callback_server");

            if (server == null) {
                log.warn("No 'callback_server' section found in {}, using default configuration", configPath);
                return new AppConfig();
            }

            AppConfig config = new AppConfig();
            config.setServiceName(server.has("service_name") ?
                    server.get("service_name").asText() : config.getServiceName());
            config.setEndpoint(server.has("endpoint") ?
                    server.get("endpoint").asText() : config.getEndpoint());
            config.setTokenExpireMinutes(server.has("token_expire_minutes") ?
                    server.get("token_expire_minutes").asInt() : config.getTokenExpireMinutes());
            config.setClientPublicKeyPath(server.has("client_public_key_path") ?
                    server.get("client_public_key_path").asText() : config.getClientPublicKeyPath());
            config.setServicePrivateKeyPath(server.has("service_private_key_path") ?
                    server.get("service_private_key_path").asText() : config.getServicePrivateKeyPath());
            config.setEnableDebug(server.has("enable_debug") ?
                    server.get("enable_debug").asBoolean() : config.isEnableDebug());

            return config;

        } catch (IOException e) {
            log.error("Failed to load config file {}: {}", configPath, e.getMessage());
            log.info("Using default configuration...");
            return new AppConfig();
        }
    }
}
