package com.libreturtle.scholarbuddy.security;

import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.service.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final UserService userService;

    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2User oauth2User = super.loadUser(userRequest);

        String email = oauth2User.getAttribute("email");
        String name = oauth2User.getAttribute("name");

        log.info("OAuth2 loadUser called - email: {}, name: {}", email, name);

        if (email == null) {
            log.error("Missing email attribute from OAuth2 provider");
            throw new OAuth2AuthenticationException("Missing required user attributes");
        }

        String provider = userRequest.getClientRegistration().getRegistrationId();
        log.info("Creating/updating user with provider: {}", provider);

        User user = userService.createOrUpdateUser(email, name != null ? name : "", provider);

        log.info("User created/updated successfully - id: {}, email: {}", user.getId(), user.getEmail());

        return new CustomOAuth2User(oauth2User, user);
    }
}
