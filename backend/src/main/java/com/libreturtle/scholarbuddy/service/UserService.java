package com.libreturtle.scholarbuddy.service;

import com.libreturtle.scholarbuddy.dto.UserResponse;
import com.libreturtle.scholarbuddy.exception.ApiException;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.model.UserRole;
import com.libreturtle.scholarbuddy.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;

    @Transactional
    public User createOrUpdateUser(UUID id, String email, String fullName) {
        return userRepository.findById(id)
                .map(user -> {
                    user.setEmail(email);
                    user.setFullName(fullName);
                    return userRepository.save(user);
                })
                .orElseGet(() -> {
                    User user = new User();
                    user.setId(id);
                    user.setEmail(email);
                    user.setFullName(fullName);
                    return userRepository.save(user);
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
                user.getRoles(),
                user.getCreatedAt(),
                user.getUpdatedAt()
        );
    }
}
