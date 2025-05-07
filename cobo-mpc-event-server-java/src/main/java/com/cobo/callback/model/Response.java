package com.cobo.callback.model;

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
public class Response {
    private int status;
    private String requestId;
    private String action;
    private String errStr;

    public static final int STATUS_OK = 0;
    public static final int STATUS_INVALID_REQUEST = 10;
    public static final int STATUS_INVALID_TOKEN = 20;
    public static final int STATUS_INTERNAL_ERROR = 30;

    public static final String ACTION_APPROVE = "APPROVE";
    public static final String ACTION_REJECT = "REJECT";
}
