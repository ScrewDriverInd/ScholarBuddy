package com.libreturtle.scholarbuddy.service;

import com.libreturtle.scholarbuddy.dto.UserResponse;
import com.libreturtle.scholarbuddy.exception.ApiException;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.model.UserRole;
import com.libreturtle.scholarbuddy.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class UserService {

    private final UserRepository userRepository;

    @Transactional
    public User createOrUpdateUser(String email, String fullName, String provider) {
        log.debug("createOrUpdateUser called - email: {}, fullName: {}, provider: {}", email, fullName, provider);

        return userRepository.findByEmail(email)
                .map(user -> {
                    log.info("User found, updating - id: {}, email: {}", user.getId(), user.getEmail());
                    user.setFullName(fullName);
                    user.setProvider(provider);
                    return userRepository.save(user);
                })
                .orElseGet(() -> {
                    log.info("User not found, creating new user - email: {}", email);
                    User user = new User();
                    user.setEmail(email);
                    user.setFullName(fullName);
                    user.setProvider(provider);
                    User saved = userRepository.save(user);
                    log.info("User created successfully - id: {}, email: {}", saved.getId(), saved.getEmail());
                    return saved;
                });
    }

    @Transactional
    public UserResponse grantAdminRole(UUID userId) {
        User user = userRepository.findById(userId)
                .orElseThrow(() -> ApiException.notFound("user was not found"));

        user.getRoles().add(UserRole.ROLE_ADMIN);
        user = userRepository.save(user);

        return mapToResponse(user);
    }

    private UserResponse mapToResponse(User user) {
        return new UserResponse(
                user.getId(),
                user.getEmail(),
                user.getFullName(),
                user.getUsername(),
                user.getProvider(),
                user.getRoles(),
                user.getCreatedAt(),
                user.getUpdatedAt()
        );
    }
}
