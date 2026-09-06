package com.libreturtle.scholarbuddy.controller;

import com.libreturtle.scholarbuddy.dto.ApiResponseBody;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
public class RootController {

    @GetMapping("/")
    public ApiResponseBody<Map<String, String>> welcome() {
        return ApiResponseBody.of(Map.of("message", "Welcome to ScholarBuddy"));
    }
}
