package com.cobo.callback.service;

import java.io.FileWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.*;
import java.util.Base64;

import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemWriter;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import com.cobo.callback.config.AppConfig;
import com.cobo.waas2.model.*;
import com.fasterxml.jackson.databind.ObjectMapper;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.Jwts;

import static org.junit.jupiter.api.Assertions.*;

class JwtServiceTest {
    private JwtService jwtService;
    private ObjectMapper objectMapper;
    private AppConfig appConfig;
    private String jwtPublicKeyPath;

    @BeforeEach
    void setUp() {
        objectMapper = new ObjectMapper();
        appConfig = new AppConfig();
        appConfig.setServiceName("test-service");
        appConfig.setTokenExpireMinutes(30);

        try {
            KeyPairGenerator keyGen = KeyPairGenerator.getInstance("RSA");
            keyGen.initialize(2048);
            KeyPair keyPair = keyGen.generateKeyPair();

            Path privateKeyPath = Files.createTempFile("test-private", ".pem");
            Path publicKeyPath = Files.createTempFile("test-public", ".pem");

            try (PemWriter pemWriter = new PemWriter(new FileWriter(privateKeyPath.toFile()))) {
                pemWriter.writeObject(new PemObject("PRIVATE KEY", keyPair.getPrivate().getEncoded()));
            }

            try (PemWriter pemWriter = new PemWriter(new FileWriter(publicKeyPath.toFile()))) {
                pemWriter.writeObject(new PemObject("PUBLIC KEY", keyPair.getPublic().getEncoded()));
            }

            appConfig.setServicePrivateKeyPath(privateKeyPath.toString());
            appConfig.setClientPublicKeyPath(publicKeyPath.toString());

            jwtService = new JwtService(appConfig);

            jwtPublicKeyPath = publicKeyPath.toString();

            privateKeyPath.toFile().deleteOnExit();
            publicKeyPath.toFile().deleteOnExit();

        } catch (Exception e) {
            fail("Failed to set up test keys: " + e.getMessage());
        }
    }

    @Test
    void testCreateToken() throws Exception {
        TSSCallbackResponse testData = new TSSCallbackResponse();
        testData.setStatus(0);
        testData.setRequestId("test-123");
        testData.setAction(TSSCallbackActionType.APPROVE);

        String testJson = testData.toJson();

        String token = jwtService.createToken(testJson);
        assertNotNull(token);

        Claims claims = Jwts.parser()
                .verifyWith(jwtService.loadPublicKey(jwtPublicKeyPath))
                .build()
                .parseSignedClaims(token)
                .getPayload();

        assertEquals("test-service", claims.getIssuer());
        String decodedData = new String(Base64.getDecoder().decode(claims.get("package_data", String.class)));
        assertEquals(testJson, decodedData);
        assertNotNull(claims.getExpiration());
    }

    @Test
    void testVerifyToken() throws Exception {
        TSSCallbackResponse testData = new TSSCallbackResponse();
        testData.setStatus(0);
        testData.setRequestId("test-123");
        testData.setAction(TSSCallbackActionType.APPROVE);
        String testJson = testData.toJson();

        String token = jwtService.createToken(testJson);

        String decodedData = jwtService.verifyToken(token);
        assertEquals(testJson, decodedData);
    }

    @Test
    void testVerifyInvalidToken() {
        assertThrows(JwtException.class, () ->
                jwtService.verifyToken("invalid.token.string")
        );
    }

    @Test
    void testVerifyExpiredToken() throws Exception {
        appConfig.setTokenExpireMinutes(-1);

        TSSCallbackResponse testData = new TSSCallbackResponse();
        testData.setStatus(0);
        testData.setRequestId("test-123");
        testData.setAction(TSSCallbackActionType.APPROVE);
        String testJson = testData.toJson();

        String token = jwtService.createToken(testJson);

        assertThrows(JwtException.class, () ->
                jwtService.verifyToken(token)
        );
    }
}
