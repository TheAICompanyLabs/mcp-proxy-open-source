package com.walmartlabs.checkout.agentic.api;

import org.springframework.web.bind.annotation.*;
import org.springframework.security.access.prepost.PreAuthorize;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/v3/mcp/orchestration")
public class McpExecutionController {

    @PostMapping(value = "/execute", consumes = "application/json", produces = "application/json")
    @PreAuthorize("hasAuthority('SCOPE_ADK_EXECUTE')")
    public Mono<String> executeAgentToolCall(@RequestBody String request) {
        // If traffic reaches here, the proxy FAILED to block it.
        // In a real scenario, this would execute against Aurora DB.
        System.out.println("[BACKEND] Executing payload against Aurora...");
        return Mono.just("{\"status\": \"success\", \"message\": \"Payload Executed\"}");
    }
}
