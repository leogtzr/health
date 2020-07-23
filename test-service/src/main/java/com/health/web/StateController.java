package com.health.web;

import com.health.domain.ServiceState;
import com.health.service.ServiceStateChanger;
import lombok.extern.log4j.Log4j2;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@Log4j2
@RestController
@RequestMapping("monitoring")
public class StateController {

    private ServiceStateChanger serviceStateChanger;

    public StateController(ServiceStateChanger serviceStateChanger) {
        this.serviceStateChanger = serviceStateChanger;
    }

    @GetMapping("healthcheck")
    public ResponseEntity<ServiceState> status() {
        final ServiceState status = this.serviceStateChanger.status();
        log.info(status);
        switch (status.getState()) {
            case DOWN:
                return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(status);
            case UP:
                return ResponseEntity.ok(status);
            case RESTARTING:
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(status);
        }
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(status);
    }

    @GetMapping("down")
    public ResponseEntity<ServiceState> down() {
        this.serviceStateChanger.down();
        final ServiceState status = this.serviceStateChanger.status();
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(status);
    }

    @GetMapping("up")
    public ResponseEntity<ServiceState> up() {
        this.serviceStateChanger.up();
        final ServiceState status = this.serviceStateChanger.status();
        return ResponseEntity.ok(status);
    }

    @GetMapping("restarting")
    public ResponseEntity<ServiceState> restarting() {
        this.serviceStateChanger.restarting();
        final ServiceState status = this.serviceStateChanger.status();
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(status);
    }

}
