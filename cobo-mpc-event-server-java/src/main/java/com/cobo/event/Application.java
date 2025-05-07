package com.cobo.event;

import java.util.HashMap;
import java.util.Map;

import com.cobo.event.config.AppConfig;
import com.cobo.event.service.JwtService;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.hellokaton.blade.Blade;
import com.hellokaton.blade.annotation.*;
import com.hellokaton.blade.annotation.request.Form;
import com.hellokaton.blade.annotation.route.GET;
import com.hellokaton.blade.annotation.route.POST;
import com.hellokaton.blade.mvc.ui.ResponseType;

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
            processEvent(eventData);
        } catch (JwtException e) {
            handleError(e.getMessage());
        } catch (Exception e) {
            log.error("Failed to process event", e);
            handleError(e.getMessage());
        }
    }

    private void handleError(String message) {
        log.error("Error: {}", message);
    }

    private static void processEvent(String event) {
        log.info("Processing event: {}", event);
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
