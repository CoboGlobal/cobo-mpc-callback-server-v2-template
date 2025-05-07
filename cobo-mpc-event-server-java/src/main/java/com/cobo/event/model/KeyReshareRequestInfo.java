package com.cobo.callback.model;

import com.cobo.waas2.model.*;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

import lombok.Data;

@Data
@JsonDeserialize(using = KeyReshareRequestInfoDeserializer.class)
public class KeyReshareRequestInfo {
    private OrgInfo org;
    private MPCProject project;
    private MPCVault vault;
    private KeyShareHolderGroup sourceKeyShareHolderGroup;
    private KeyShareHolderGroup targetKeyShareHolderGroup;
    private TSSRequest tssRequest;
}
