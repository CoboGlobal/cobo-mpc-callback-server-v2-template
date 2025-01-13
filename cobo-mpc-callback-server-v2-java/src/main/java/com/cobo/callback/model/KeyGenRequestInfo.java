package com.cobo.callback.model;

import com.cobo.waas2.model.*;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

import lombok.Data;

@Data
@JsonDeserialize(using = KeyGenRequestInfoDeserializer.class)
public class KeyGenRequestInfo {
    private OrgInfo org;
    private MPCProject project;
    private MPCVault vault;
    private KeyShareHolderGroup targetKeyShareHolderGroup;
    private TSSRequest tssRequest;
}
