package com.libreturtle.scholarbuddy.validation;

import jakarta.validation.Constraint;
import jakarta.validation.Payload;

import java.lang.annotation.*;

@Target({ElementType.FIELD, ElementType.PARAMETER})
@Retention(RetentionPolicy.RUNTIME)
@Constraint(validatedBy = UrlValidator.class)
@Documented
public @interface ValidUrl {
    String message() default "link must be an absolute URL";
    Class<?>[] groups() default {};
    Class<? extends Payload>[] payload() default {};
}
