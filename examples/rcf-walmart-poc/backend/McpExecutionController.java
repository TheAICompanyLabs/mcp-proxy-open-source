package com.walmartlabs.checkout.agentic.api;

import org.springframework.web.bind.annotation.*;
import org.springframework.security.access.prepost.PreAuthorize;
import reactor.core.publisher.Mono;
import com.walmartlabs.checkout.agentic.service.McpToolRegistry;
import com.walmartlabs.checkout.agentic.kafka.EventPublisher;
import com.walmartlabs.mcp.protocol.McpToolRequest;
import com.walmartlabs.mcp.protocol.McpResponse;

@RestController
@RequestMapping("/api/v3/mcp/orchestration")
public class McpExecutionController {

    private final McpToolRegistry toolRegistry;
    private final EventPublisher auditStream;

    public McpExecutionController(McpToolRegistry toolRegistry, EventPublisher auditStream) {
        this.toolRegistry = toolRegistry;
        this.auditStream = auditStream;
    }

    /**
     * ADK Execution Endpoint
     * Vulnerability Note: WSO2/OAuth terminates at the perimeter. 
     * If RCF is bypassed, TOCTOU payload mutations will execute against Aurora here.
     */
    @PostMapping(value = "/execute", consumes = "application/json", produces = "application/json")
    @PreAuthorize("hasAuthority('SCOPE_ADK_EXECUTE')")
    public Mono<McpResponse> executeAgentToolCall(@RequestBody McpToolRequest request) {
        
        return toolRegistry.dispatch(request)
            .doOnSuccess(res -> auditStream.publish("mcp-checkout-success", request.getAgentId()))
            .doOnError(err -> auditStream.publishDeadLetter("mcp-dlq", err.getMessage()));
    }
}
