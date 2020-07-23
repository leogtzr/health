package com.health.config;

import com.health.domain.ServiceState;
import com.health.domain.State;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class Config {

    @Bean
    public ServiceState serviceState() {
        final ServiceState serviceState = new ServiceState(State.UP);
        return serviceState;
    }

}
