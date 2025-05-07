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
public class KeyReshareDetail {
    private String oldGroupId;
    private String rootPubKey;
    private CurveType curve;
    private List<String> usedNodeIds;
    private int oldThreshold;
    private int newThreshold;
    private List<String> newNodeIds;
    private String taskId;
    private String bizTaskId;
}
