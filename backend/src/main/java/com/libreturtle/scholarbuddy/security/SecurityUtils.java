package com.libreturtle.scholarbuddy.security;

import com.libreturtle.scholarbuddy.exception.ApiException;
import com.libreturtle.scholarbuddy.model.User;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

public class SecurityUtils {

    private SecurityUtils() {}

    public static User getCurrentUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (authentication == null) {
            throw ApiException.unauthorized("valid authentication is required");
        }

        Object principal = authentication.getPrincipal();
        if (principal instanceof CustomOidcUser customOidcUser) {
            return customOidcUser.getUser();
        }
        if (principal instanceof CustomOAuth2User customOAuth2User) {
            return customOAuth2User.getUser();
        }

        throw ApiException.unauthorized("valid authentication is required");
    }
}
