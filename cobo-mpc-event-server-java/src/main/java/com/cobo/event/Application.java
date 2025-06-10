package com.cobo.event;

import java.util.HashMap;
import java.util.Map;

import com.cobo.event.config.AppConfig;
import com.cobo.event.service.JwtService;
import com.cobo.waas2.model.*;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.hellokaton.blade.Blade;
import com.hellokaton.blade.annotation.*;
import com.hellokaton.blade.annotation.request.Form;
import com.hellokaton.blade.annotation.route.GET;
import com.hellokaton.blade.annotation.route.POST;
import com.hellokaton.blade.mvc.ui.ResponseType;
import com.hellokaton.blade.mvc.WebContext;

import io.jsonwebtoken.JwtException;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Path
public class Application {
    private static final ObjectMapper MAPPER = new ObjectMapper()
            .setPropertyNamingStrategy(PropertyNamingStrategies.SNAKE_CASE)
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
            .configure(DeserializationFeature.FAIL_ON_NULL_FOR_PRIMITIVES, false)
            .configure(DeserializationFeature.READ_UNKNOWN_ENUM_VALUES_AS_NULL, true)
            .configure(DeserializationFeature.ACCEPT_EMPTY_STRING_AS_NULL_OBJECT, true);

    private static JwtService jwtService;
    private static AppConfig appConfig;

    @GET(value = "/ping", responseType = ResponseType.JSON)
    public Map<String, String> ping() {
        Map<String, String> response = new HashMap<>();
        response.put("server", appConfig.getServiceName());
        response.put("timestamp", String.valueOf(System.currentTimeMillis()));
        return response;
    }

    @POST(value = "/v2/event", responseType = ResponseType.EMPTY)
    public void event(@Form String TSS_JWT_MSG) {
        try {
            if (TSS_JWT_MSG.isEmpty()) {
                throw new IllegalArgumentException("Missing TSS_JWT_MSG parameter");
            }

            String eventData = jwtService.verifyToken(TSS_JWT_MSG);
            handleEvent(eventData);

            // Set HTTP 200 for successful processing
            WebContext.response().status(200);
        } catch (JwtException e) {
            log.error("JWT verification failed: {}", e.getMessage());
            WebContext.response().status(400);
        } catch (Exception e) {
            log.error("Failed to process event: {}", e.getMessage());
            WebContext.response().status(400);
        }
    }

    private static void handleEvent(String eventJson) {
        log.info("Handle event: {}", eventJson);
        
        try {
            TSSEvent event = TSSEvent.fromJson(eventJson);
            log.info("Event: {}", event);

            if (event.getEventType() == TSSEventType.PING) {
                log.info("Ping event: {}", event);
                return;
            }

            TSSEventData eventData = event.getData();
            log.info("Event data: {}", eventData);

            assert eventData != null;
            if (eventData.getActualInstance() instanceof TSSKeyGenEventData) {
                TSSKeyGenEventData keyGenEventData = (TSSKeyGenEventData) eventData.getActualInstance();
                log.info("Key gen request: {}", keyGenEventData);

                TSSKeyGenExtra extra = TSSKeyGenExtra.fromJson(String.valueOf(keyGenEventData.getExtraInfo()));
                log.info("Key gen extra: {}", extra);
            } else if (eventData.getActualInstance() instanceof TSSKeySignEventData) {
                TSSKeySignEventData keySignEventData = (TSSKeySignEventData) eventData.getActualInstance();
                log.info("Key sign request: {}", keySignEventData);

                TSSKeySignExtra extra = TSSKeySignExtra.fromJson(String.valueOf(keySignEventData.getExtraInfo()));
                log.info("Key sign extra: {}", extra);
            } else if (eventData.getActualInstance() instanceof TSSKeyReshareEventData) {
                TSSKeyReshareEventData keyReshareEventData = (TSSKeyReshareEventData) eventData.getActualInstance();
                log.info("Key reshare request: {}", keyReshareEventData);

                TSSKeyReshareExtra extra = TSSKeyReshareExtra.fromJson(String.valueOf(keyReshareEventData.getExtraInfo()));
                log.info("Key reshare extra: {}", extra);
            } else if (eventData.getActualInstance() instanceof TSSKeyShareSignEventData) {
                TSSKeyShareSignEventData keyShareSignEventData = (TSSKeyShareSignEventData) eventData.getActualInstance();
                log.info("Key share sign request: {}", keyShareSignEventData);

                TSSKeyShareSignExtra extra = TSSKeyShareSignExtra.fromJson(String.valueOf(keyShareSignEventData.getExtraInfo()));
                log.info("Key share sign extra: {}", extra);
            }
        

        // Add your event handle logic here


        } catch (Exception e) {
            log.error("Failed to parse event: {}", e.getMessage());
        }
    }

    public static void main(String[] args) {
        appConfig = AppConfig.loadConfig(args);
        log.info("Loaded configuration: {}", appConfig);

        jwtService = new JwtService(appConfig);

        String[] hostAndPort = appConfig.getEndpoint().split(":");
        int port = Integer.parseInt(hostAndPort[1]);
        log.info("Starting server on port {}", port);

        Blade.create().listen(port).start(Application.class, args);
    }
}
