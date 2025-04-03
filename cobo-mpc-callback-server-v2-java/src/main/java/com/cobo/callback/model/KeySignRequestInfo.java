package com.cobo.callback.model;

import java.util.List;

import com.cobo.waas2.model.*;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

import lombok.Data;

@Data
@JsonDeserialize(using = KeySignRequestInfoDeserializer.class)
public class KeySignRequestInfo {
    private OrgInfo org;
    private MPCProject project;
    private MPCVault vault;
    private WalletInfo wallet;
    private KeyShareHolderGroup signerKeyShareHolderGroup;
    private List<AddressInfo> sourceAddresses;
    private Transaction transaction;
}
