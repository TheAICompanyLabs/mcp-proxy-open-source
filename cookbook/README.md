# The MCP Proxy Policy Cookbook

While the MCP Proxy automatically discovers tools exposed by the underlying server, it defaults to a **fail-closed (Zero-Trust)** posture. To allow an AI agent to execute tools safely, you must provide explicit rules in a `policy.yaml` file in your execution directory.

Choose a scenario below and copy the YAML template to secure your environment.

* [📁 File System Security](./filesystem/policy.yaml)
* [🗄️ Database Security](./database/policy.yaml)
* [🌐 Web API Security (SSRF Protection)](./web-api/policy.yaml)
