package com.cobo.callback.model;

import java.util.ArrayList;
import java.util.List;

import com.cobo.waas2.model.*;
import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;

import io.jsonwebtoken.io.IOException;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class KeySignRequestInfoDeserializer extends JsonDeserializer<KeySignRequestInfo> {
    @Override
    public KeySignRequestInfo deserialize(JsonParser p, DeserializationContext ctxt) throws IOException, java.io.IOException {
        KeySignRequestInfo result = new KeySignRequestInfo();
        JsonNode node = p.getCodec().readTree(p);

        try {
            if (node.has("org") && !node.get("org").isNull()) {
                String orgJson = node.get("org").toString();
                result.setOrg(OrgInfo.fromJson(orgJson));
            }

            if (node.has("project") && !node.get("project").isNull()) {
                String projectJson = node.get("project").toString();
                result.setProject(MPCProject.fromJson(projectJson));
            }

            if (node.has("vault") && !node.get("vault").isNull()) {
                String vaultJson = node.get("vault").toString();
                result.setVault(MPCVault.fromJson(vaultJson));
            }

            if (node.has("wallet") && !node.get("wallet").isNull()) {
                String walletJson = node.get("wallet").toString();
                result.setWallet(WalletInfo.fromJson(walletJson));
            }

            if (node.has("signer_key_share_holder_group") && !node.get("signer_key_share_holder_group").isNull()) {
                String signerGroupJson = node.get("signer_key_share_holder_group").toString();
                result.setSignerKeyShareHolderGroup(KeyShareHolderGroup.fromJson(signerGroupJson));
            }

            if (node.has("source_addresses") && !node.get("source_addresses").isNull() && node.get("source_addresses").isArray()) {
                List<AddressInfo> addresses = new ArrayList<>();
                for (JsonNode addrNode : node.get("source_addresses")) {
                    String addressJson = addrNode.toString();
                    addresses.add(AddressInfo.fromJson(addressJson));
                }
                result.setSourceAddresses(addresses);
            }

            if (node.has("transaction") && !node.get("transaction").isNull()) {
                String transactionJson = node.get("transaction").toString();
                result.setTransaction(Transaction.fromJson(transactionJson));
            }

        } catch (Exception e) {
            log.error("Error deserializing KeySignRequestInfo", e);
        }

        return result;
    }
}
