package com.health.service;

import com.health.domain.ServiceState;
import com.health.domain.State;
import org.springframework.stereotype.Service;

@Service
public class ServiceStateChanger {

    private ServiceState serviceState;

    public ServiceStateChanger(final ServiceState serviceState) {
        this.serviceState = serviceState;
    }

    public void down() {
        this.serviceState.setState(State.DOWN);
    }

    public void up() {
        this.serviceState.setState(State.UP);
    }

    public void restarting() {
        this.serviceState.setState(State.RESTARTING);
    }

    public ServiceState status() {
        return this.serviceState;
    }

}
