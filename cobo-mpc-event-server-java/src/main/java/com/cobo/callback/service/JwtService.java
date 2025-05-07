package com.cobo.callback.service;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;
import java.util.Date;

import com.cobo.callback.config.AppConfig;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.Jwts;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class JwtService {
    private static final String ALGORITHM = "RSA";
    private static final String PACKAGE_DATA_CLAIM = "package_data";
    private static final String BEGIN_PUBLIC_KEY = "-----BEGIN PUBLIC KEY-----";
    private static final String END_PUBLIC_KEY = "-----END PUBLIC KEY-----";
    private static final String BEGIN_PRIVATE_KEY = "-----BEGIN PRIVATE KEY-----";
    private static final String END_PRIVATE_KEY = "-----END PRIVATE KEY-----";

    private final AppConfig config;
    private final PublicKey clientPublicKey;
    private final PrivateKey servicePrivateKey;

    public JwtService(AppConfig config) {
        this.config = config;
        try {
            this.clientPublicKey = loadPublicKey(config.getClientPublicKeyPath());
            this.servicePrivateKey = loadPrivateKey(config.getServicePrivateKeyPath());
            log.info("Successfully loaded JWT keys");
        } catch (Exception e) {
            log.error("Failed to load JWT keys", e);
            throw new RuntimeException("Failed to initialize JWT service", e);
        }
    }

    public String createToken(String data) {
        if (data == null || data.isEmpty()) {
            throw new IllegalArgumentException("Token data cannot be null or empty");
        }

        try {
            Date now = new Date();
            Date expiration = new Date(now.getTime() + (config.getTokenExpireMinutes() * 60 * 1000L));

            return Jwts.builder()
                    .expiration(expiration)
                    .issuer(config.getServiceName())
                    .claim(PACKAGE_DATA_CLAIM, Base64.getEncoder().encodeToString(data.getBytes()))
                    .signWith(servicePrivateKey, Jwts.SIG.RS256)
                    .compact();

        } catch (Exception e) {
            log.error("Failed to create token", e);
            throw new JwtException("Failed to create token: " + e.getMessage());
        }
    }

    public String verifyToken(String token) {
        if (token == null || token.isEmpty()) {
            throw new IllegalArgumentException("Token cannot be null or empty");
        }

        try {
            Claims claims = Jwts.parser()
                    .verifyWith(clientPublicKey)
                    .build()
                    .parseSignedClaims(token)
                    .getPayload();

            String packageData = claims.get(PACKAGE_DATA_CLAIM, String.class);
            if (packageData == null || packageData.isEmpty()) {
                throw new JwtException("Token missing package_data claim");
            }

            return new String(Base64.getDecoder().decode(packageData));
        } catch (JwtException e) {
            log.error("JWT validation failed", e);
            throw e;
        } catch (Exception e) {
            log.error("Failed to verify token", e);
            throw new JwtException("Invalid token: " + e.getMessage());
        }
    }

    public PublicKey loadPublicKey(String path) throws Exception {
        byte[] keyBytes = readKeyFile(path);
        try {
            String keyContent = new String(keyBytes)
                    .replace(BEGIN_PUBLIC_KEY, "")
                    .replace(END_PUBLIC_KEY, "")
                    .replaceAll("\\s", "");

            byte[] decoded = Base64.getDecoder().decode(keyContent);
            KeyFactory keyFactory = KeyFactory.getInstance(ALGORITHM);
            return keyFactory.generatePublic(new X509EncodedKeySpec(decoded));
        } catch (InvalidKeySpecException e) {
            log.error("Invalid public key format", e);
            throw new IllegalArgumentException("Invalid public key format", e);
        }
    }

    public PrivateKey loadPrivateKey(String path) throws Exception {
        byte[] keyBytes = readKeyFile(path);
        try {
            String keyContent = new String(keyBytes)
                    .replace(BEGIN_PRIVATE_KEY, "")
                    .replace(END_PRIVATE_KEY, "")
                    .replaceAll("\\s", "");

            byte[] decoded = Base64.getDecoder().decode(keyContent);
            KeyFactory keyFactory = KeyFactory.getInstance(ALGORITHM);
            return keyFactory.generatePrivate(new PKCS8EncodedKeySpec(decoded));
        } catch (InvalidKeySpecException e) {
            log.error("Invalid private key format", e);
            throw new IllegalArgumentException("Invalid private key format", e);
        }
    }

    private byte[] readKeyFile(String path) throws IOException {
        Path keyPath = Paths.get(path);
        if (!Files.exists(keyPath)) {
            throw new IOException("Key file not found: " + path);
        }
        return Files.readAllBytes(keyPath);
    }
}
