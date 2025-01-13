package com.cobo.callback.model;

import java.util.List;

import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.annotation.JsonNaming;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
public class KeySignDetail {
    private String groupId;
    private String rootPubKey;
    private List<String> usedNodeIds;
    private List<String> bip32PathList;
    private List<String> msgHashList;
    private List<String> tweakList;
    private SignatureType signatureType;
    private TssProtocol tssProtocol;
    private String taskId;
    private String bizTaskId;
}
