# Nebulae

Nebulae is a powerful tool designed to parse code repositories and generate various visualization charts, helping developers better understand and manage complex software systems.

## Key Features

1. **Dependency Graph**: Displays the dependencies between different modules, classes, or packages in the project, helping developers understand code structure and organization.

2. **Network Call Topology**: Visually represents the network call relationships between various services or components in the system, aiding in understanding system architecture and communication patterns.

3. **Middleware Association Diagram**: Shows the relationships between the application and various middleware (such as databases, message queues, caches, etc.), facilitating understanding of the system's technology stack and integration points.

## How It Works

Nebulae generates these visualization charts through the following steps:

1. Parses the code repository, analyzing source code files.
2. Extracts key information, such as configuration files, Maven Pom, etc.
3. Builds an internal data model representing various entities and their relationships.
4. Uses a graphics rendering engine to convert the data model into visual charts.

## Advantages

- Improves code maintainability: By visually presenting code structure, it helps developers more easily understand and maintain complex systems.
- Optimizes system architecture: Identifies potential architectural issues, such as circular dependencies or excessive coupling.
- Aids technical decision-making: Provides technical teams with an intuitive overview of the system, assisting in architectural design and optimization decisions.
- Promotes team collaboration: Offers team members a shared view of the system, facilitating communication and collaboration.

Nebulae is a powerful tool that can help development teams better understand, manage, and optimize their software systems.

## Command

### Grpc protoc

```
protoc -I=proto \
    --go_out=proto --go_opt=paths=source_relative \
    --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative \
    proto/v1/nebulae.proto
```