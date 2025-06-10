package com.cobo.callback.verify;

import com.cobo.waas2.model.*;
import lombok.extern.slf4j.Slf4j;

import java.util.Objects;

@Slf4j
public class TssVerifier implements Verifier {

    public static TssVerifier create() {
        return new TssVerifier();
    }

    @Override
    public String verify(TSSCallbackRequest request) {
        if (request == null) {
            return "request is nil";
        }

        try {
            TSSCallbackRequestType requestType = Objects.requireNonNull(request.getRequestType());
            switch (requestType) {
                case PING:
                    log.debug("Got ping request");
                    return null;
                case KEYGEN:
                    return handleKeyGen(request.getRequestDetail(), request.getExtraInfo());
                case KEYSIGN:
                    return handleKeySign(request.getRequestDetail(), request.getExtraInfo());
                case KEYRESHARE:
                    return handleKeyReshare(request.getRequestDetail(), request.getExtraInfo());
                default:
                    return "not support to process request type " + requestType;
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

            TSSKeyGenRequest detail =  TSSKeyGenRequest.fromJson(requestDetail);
            TSSKeyGenExtra extra = TSSKeyGenExtra.fromJson(extraInfo);

            log.debug("key gen class detail:\n{}\nextra:\n{}", detail, extra);

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

            TSSKeySignRequest detail = TSSKeySignRequest.fromJson(requestDetail);
            TSSKeySignExtra extra = TSSKeySignExtra.fromJson(extraInfo);

            log.debug("key sign class detail:\n{}\nextra:\n{}", detail, extra);

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

            TSSKeyReshareRequest detail = TSSKeyReshareRequest.fromJson(requestDetail);
            TSSKeyReshareExtra extra = TSSKeyReshareExtra.fromJson(extraInfo);

            log.debug("key reshare class detail:\n{}\nextra:\n{}", detail, extra);

            // key reshare logic add here

            return null;
        } catch (Exception e) {
            log.error("Failed to handle key reshare", e);
            return "failed to handle key reshare: " + e.getMessage();
        }
    }
}
