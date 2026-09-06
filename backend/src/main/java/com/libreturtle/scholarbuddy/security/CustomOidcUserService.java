package com.libreturtle.scholarbuddy.security;

import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.service.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.oauth2.client.oidc.userinfo.OidcUserRequest;
import org.springframework.security.oauth2.client.oidc.userinfo.OidcUserService;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class CustomOidcUserService extends OidcUserService {

    private final UserService userService;

    @Override
    public OidcUser loadUser(OidcUserRequest userRequest) throws OAuth2AuthenticationException {
        OidcUser oidcUser = super.loadUser(userRequest);

        String email = oidcUser.getAttribute("email");
        String name = oidcUser.getAttribute("name");

        log.info("OIDC loadUser called - email: {}, name: {}", email, name);

        if (email == null) {
            log.error("Missing email attribute from OIDC provider");
            throw new OAuth2AuthenticationException("Missing required user attributes");
        }

        String provider = userRequest.getClientRegistration().getRegistrationId();
        log.info("Creating/updating user with provider: {}", provider);

        User user = userService.createOrUpdateUser(email, name != null ? name : "", provider);

        log.info("User created/updated successfully - id: {}, email: {}", user.getId(), user.getEmail());

        return new CustomOidcUser(oidcUser, user);
    }
}
