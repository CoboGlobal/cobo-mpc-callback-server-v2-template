package com.cobo.callback.verify;

import com.cobo.callback.model.*;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class TssVerifier implements Verifier {
    private static final ObjectMapper MAPPER = new ObjectMapper()
            .setPropertyNamingStrategy(PropertyNamingStrategies.SNAKE_CASE)
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
            .configure(DeserializationFeature.FAIL_ON_NULL_FOR_PRIMITIVES, false)
            .configure(DeserializationFeature.READ_UNKNOWN_ENUM_VALUES_AS_NULL, true)
         .configure(DeserializationFeature.ACCEPT_EMPTY_STRING_AS_NULL_OBJECT, true);

    public static TssVerifier create() {
        return new TssVerifier();
    }

    @Override
    public String verify(Request request) {
        if (request == null) {
            return "request is nil";
        }

        try {
            switch (request.getRequestType()) {
                case TYPE_PING:
                    log.debug("Got ping request");
                    return null;
                case TYPE_KEY_GEN:
                    return handleKeyGen(request.getRequestDetail(), request.getExtraInfo());
                case TYPE_KEY_SIGN:
                    return handleKeySign(request.getRequestDetail(), request.getExtraInfo());
                case TYPE_KEY_RESHARE:
                    return handleKeyReshare(request.getRequestDetail(), request.getExtraInfo());
                default:
                    return "not support to process request type " + request.getRequestType();
            }
        } catch (Exception e) {
            log.error("Failed to verify request", e);
            return e.getMessage();
        }
    }

    private String handleKeyGen(String requestDetail, String extraInfo) {
        try {
            if (requestDetail == null || requestDetail.isEmpty() || extraInfo == null || extraInfo.isEmpty()) {
                return "request detail or extra info is empty";
            }
            log.debug("key gen original detail:\n{}\nrequest info:\n{}", requestDetail, extraInfo);

            KeyGenDetail detail = MAPPER.readValue(requestDetail, KeyGenDetail.class);
            KeyGenRequestInfo requestInfo = MAPPER.readValue(extraInfo, KeyGenRequestInfo.class);

            log.debug("key gen class detail:\n{}\nrequest info:\n{}", detail, requestInfo);

            // key gen logic add here

            return null;
        } catch (Exception e) {
            log.error("Failed to handle key gen", e);
            return "failed to handle key gen: " + e.getMessage();
        }
    }

    private String handleKeySign(String requestDetail, String extraInfo) {
        try {
            if (requestDetail == null || requestDetail.isEmpty() || extraInfo == null || extraInfo.isEmpty()) {
                return "request detail or extra info is empty";
            }
            log.debug("key sign original detail:\n{}\nrequest info:\n{}", requestDetail, extraInfo);

            KeySignDetail detail = MAPPER.readValue(requestDetail, KeySignDetail.class);
            KeySignRequestInfo requestInfo = MAPPER.readValue(extraInfo, KeySignRequestInfo.class);

            log.debug("key sign class detail:\n{}\nrequest info:\n{}", detail, requestInfo);

            // key sign logic add here

            return null;
        } catch (Exception e) {
            log.error("Failed to handle key sign", e);
            return "failed to handle key sign: " + e.getMessage();
        }
    }

    private String handleKeyReshare(String requestDetail, String extraInfo) {
        try {
            if (requestDetail == null || requestDetail.isEmpty() || extraInfo == null || extraInfo.isEmpty()) {
                return "request detail or extra info is empty";
            }
            log.debug("key reshare original detail:\n{}\nrequest info:\n{}", requestDetail, extraInfo);

            KeyReshareDetail detail = MAPPER.readValue(requestDetail, KeyReshareDetail.class);
            KeyReshareRequestInfo requestInfo = MAPPER.readValue(extraInfo, KeyReshareRequestInfo.class);

            log.debug("key reshare class detail:\n{}\nrequest info:\n{}", detail, requestInfo);

            // key reshare logic add here

            return null;
        } catch (Exception e) {
            log.error("Failed to handle key reshare", e);
            return "failed to handle key reshare: " + e.getMessage();
        }
    }
}
