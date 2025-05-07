package com.cobo.callback.model;

import com.cobo.waas2.model.*;
import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;

import io.jsonwebtoken.io.IOException;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class KeyGenRequestInfoDeserializer extends JsonDeserializer<KeyGenRequestInfo> {
    @Override
    public KeyGenRequestInfo deserialize(JsonParser p, DeserializationContext ctxt) throws IOException, java.io.IOException {
        KeyGenRequestInfo result = new KeyGenRequestInfo();
        JsonNode node = p.getCodec().readTree(p);

        try {
            if (node.has("vault") && !node.get("vault").isNull()) {
                String vaultJson = node.get("vault").toString();
                result.setVault(MPCVault.fromJson(vaultJson));
            }

            if (node.has("project") && !node.get("project").isNull()) {
                String projectJson = node.get("project").toString();
                result.setProject(MPCProject.fromJson(projectJson));
            }

            if (node.has("org") && !node.get("org").isNull()) {
                String orgJson = node.get("org").toString();
                result.setOrg(OrgInfo.fromJson(orgJson));
            }

            if (node.has("target_key_share_holder_group") && !node.get("target_key_share_holder_group").isNull()) {
                String groupJson = node.get("target_key_share_holder_group").toString();
                result.setTargetKeyShareHolderGroup(KeyShareHolderGroup.fromJson(groupJson));
            }

            if (node.has("tss_request") && !node.get("tss_request").isNull()) {
                String requestJson = node.get("tss_request").toString();
                result.setTssRequest(TSSRequest.fromJson(requestJson));
            }
        } catch (Exception e) {
            log.error("Error deserializing KeyGenRequestInfo", e);
        }

        return result;
    }
}
