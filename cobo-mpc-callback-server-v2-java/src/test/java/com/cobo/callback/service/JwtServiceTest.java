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
    private KeyPair keyPair;

    @BeforeEach
    void setUp() {
        objectMapper = new ObjectMapper();
        appConfig = new AppConfig();
        appConfig.setServiceName("test-service");
        appConfig.setTokenExpireMinutes(30);

        try {
            KeyPairGenerator keyGen = KeyPairGenerator.getInstance("RSA");
            keyGen.initialize(2048);
            keyPair = keyGen.generateKeyPair();

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

    /**
     * Generate keys with openssl in BOTH PKCS#1 and PKCS#8 formats, then exercise the full
     * Java path: loadPrivateKey -> createToken (sign with private key) -> verifyToken
     * (verify with the openssl-derived public key). Both formats must succeed.
     */
    @Test
    void testOpensslKeysPkcs1AndPkcs8() throws Exception {
        Path dir = Files.createTempDirectory("openssl-keys");

        Path priPkcs1 = dir.resolve("callback-server-pri.pem");        // PKCS#1 (openssl genrsa)
        Path priPkcs8 = dir.resolve("callback-server-pri-pkcs8.pem");  // PKCS#8 (converted)
        Path pub = dir.resolve("callback-server-pub.key");             // SPKI public key

        // 1. PKCS#1 private key (exactly what the doc's `openssl genrsa` produces)
        runOpenssl("genrsa", "-out", priPkcs1.toString(), "2048");
        // 2. public key derived from it (doc command: openssl rsa -in pri.pem -pubout)
        runOpenssl("rsa", "-in", priPkcs1.toString(), "-pubout", "-out", pub.toString());
        // 3. convert the SAME private key to PKCS#8 (no key regeneration)
        runOpenssl("pkcs8", "-topk8", "-nocrypt", "-in", priPkcs1.toString(), "-out", priPkcs8.toString());

        // sanity: files really are the two different formats
        assertTrue(new String(Files.readAllBytes(priPkcs1)).contains("-----BEGIN RSA PRIVATE KEY-----"),
                "expected PKCS#1 header");
        assertTrue(new String(Files.readAllBytes(priPkcs8)).contains("-----BEGIN PRIVATE KEY-----"),
                "expected PKCS#8 header");

        String payload = "{\"hello\":\"world\"}";

        // PKCS#1: load -> sign -> verify with openssl-derived public key
        String decodedFromPkcs1 = signAndVerify(priPkcs1.toString(), pub.toString(), payload);
        assertEquals(payload, decodedFromPkcs1, "PKCS#1 private key round-trip failed");

        // PKCS#8: load -> sign -> verify with the SAME openssl-derived public key
        String decodedFromPkcs8 = signAndVerify(priPkcs8.toString(), pub.toString(), payload);
        assertEquals(payload, decodedFromPkcs8, "PKCS#8 private key round-trip failed");

        // Both formats are the same underlying key, so the public key produces identical results
        assertEquals(decodedFromPkcs1, decodedFromPkcs8);
    }

    private String signAndVerify(String privateKeyPath, String publicKeyPath, String payload) {
        AppConfig cfg = new AppConfig();
        cfg.setServiceName("test-service");
        cfg.setTokenExpireMinutes(30);
        cfg.setServicePrivateKeyPath(privateKeyPath);
        cfg.setClientPublicKeyPath(publicKeyPath);

        JwtService svc = new JwtService(cfg);
        String token = svc.createToken(payload);
        return svc.verifyToken(token);
    }

    private void runOpenssl(String... args) throws Exception {
        java.util.List<String> cmd = new java.util.ArrayList<>();
        cmd.add("openssl");
        for (String a : args) {
            cmd.add(a);
        }
        Process p = new ProcessBuilder(cmd).redirectErrorStream(true).start();
        byte[] out = p.getInputStream().readAllBytes();
        int code = p.waitFor();
        assertEquals(0, code, "openssl " + String.join(" ", args) + " failed:\n" + new String(out));
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
