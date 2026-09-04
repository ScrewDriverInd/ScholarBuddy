package com.libreturtle.scholarbuddy.dto;

import com.libreturtle.scholarbuddy.model.UserRole;

import java.time.Instant;
import java.util.Set;
import java.util.UUID;

public record UserResponse(
        UUID id,
        String email,
        String fullName,
        String username,
        Set<UserRole> roles,
        Instant createdAt,
        Instant updatedAt
) {
}
