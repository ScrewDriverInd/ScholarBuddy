package com.libreturtle.scholarbuddy.controller;

import com.libreturtle.scholarbuddy.dto.ApiResponseBody;
import com.libreturtle.scholarbuddy.dto.UserResponse;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.security.SecurityUtils;
import com.libreturtle.scholarbuddy.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/user")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/me")
    public ResponseEntity<ApiResponseBody<UserResponse>> getCurrentUser() {
        User user = SecurityUtils.getCurrentUser();
        UserResponse response = new UserResponse(
                user.getId(),
                user.getEmail(),
                user.getFullName(),
                user.getUsername(),
                user.getRoles(),
                user.getCreatedAt(),
                user.getUpdatedAt()
        );
        return ResponseEntity.ok(ApiResponseBody.of(response));
    }
}
