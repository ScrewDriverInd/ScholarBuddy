package com.libreturtle.scholarbuddy.security;

import com.libreturtle.scholarbuddy.exception.ApiException;
import com.libreturtle.scholarbuddy.model.User;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

public class SecurityUtils {

    private SecurityUtils() {}

    public static User getCurrentUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (authentication == null || !(authentication.getPrincipal() instanceof CustomOAuth2User principal)) {
            throw ApiException.unauthorized("valid authentication is required");
        }
        return principal.getUser();
    }
}
