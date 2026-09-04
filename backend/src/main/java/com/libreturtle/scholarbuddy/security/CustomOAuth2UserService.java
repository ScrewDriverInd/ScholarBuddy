package com.libreturtle.scholarbuddy.security;

import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final UserService userService;

    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2User oauth2User = super.loadUser(userRequest);

        String sub = oauth2User.getAttribute("sub");
        String email = oauth2User.getAttribute("email");
        String name = oauth2User.getAttribute("name");

        if (sub == null || email == null) {
            throw new OAuth2AuthenticationException("Missing required user attributes");
        }

        UUID userId = UUID.fromString(sub);
        User user = userService.createOrUpdateUser(userId, email, name != null ? name : "");

        return new CustomOAuth2User(oauth2User, user);
    }
}
