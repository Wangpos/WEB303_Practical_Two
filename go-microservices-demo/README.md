# WEB303 Practical 2: API Gateway with Service Discovery

## Repository Information

## Project Overview

This project demonstrates a microservices architecture with service discovery using Consul and an API Gateway for request routing. The system consists of:

- **Consul**: Service registry and discovery system
- **Users Service**: Microservice handling user-related requests
- **Products Service**: Microservice handling product-related requests
- **API Gateway**: Entry point that routes requests to appropriate services

## Architecture

```
Client Request → API Gateway → Service Discovery (Consul) → Target Service
                      ↓
                 Route Request
                      ↓
              [Users Service] or [Products Service]
```

## Implementation Approach

### 1. Project Structure

```
go-microservices-demo/
├── api-gateway/
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── services/
    ├── products-service/
    │   ├── go.mod
    │   ├── go.sum
    │   └── main.go
    └── users-service/
        ├── go.mod
        ├── go.sum
        └── main.go
```

### 2. Service Implementation Details

#### Users Service (Port 8081)

- **Endpoints:**
  - `GET /users/{id}` - Get user details
  - `GET /health` - Health check endpoint for Consul
- **Registration:** Registers with Consul as "users-service"
- **Dependencies:** Chi router, Consul API client

#### Products Service (Port 8082)

- **Endpoints:**
  - `GET /products/{id}` - Get product details
  - `GET /health` - Health check endpoint for Consul
- **Registration:** Registers with Consul as "products-service"
- **Dependencies:** Chi router, Consul API client

#### API Gateway (Port 8080)

- **Functionality:** Routes requests based on URL patterns
- **Routing Logic:**
  - `/api/users/*` → users-service
  - `/api/products/*` → products-service
- **Service Discovery:** Queries Consul for healthy service instances
- **Load Balancing:** Uses first available healthy instance

### 3. Service Discovery Implementation

Each microservice registers itself with Consul on startup:

```go
registration := &consulapi.AgentServiceRegistration{
    ID:      fmt.Sprintf("%s-%s", serviceName, hostname),
    Name:    serviceName,
    Port:    servicePort,
    Address: "localhost",
    Check: &consulapi.AgentServiceCheck{
        HTTP:     fmt.Sprintf("http://localhost:%d/health", servicePort),
        Interval: "10s",
        Timeout:  "1s",
    },
}
```

The API Gateway discovers services by querying Consul's health API:

```go
services, _, err := consul.Health().Service(name, "", true, nil)
```

## Steps Taken

1. **Environment Setup**

   - Installed Go programming language
   - Installed Consul service registry
   - Set up project directory structure

2. **Service Development**

   - Created independent Go modules for each service
   - Implemented HTTP endpoints using Chi router
   - Added Consul registration logic with health checks
   - Fixed hostname resolution issues by using "localhost"

3. **API Gateway Development**

   - Implemented request routing logic
   - Added service discovery integration
   - Created reverse proxy functionality

4. **Testing and Validation**
   - Started Consul in development mode
   - Launched all microservices
   - Verified service registration in Consul
   - Tested end-to-end request flow

## Challenges Encountered

### 1. Hostname Resolution Issue

**Problem:** Services registered with hostname "wangs" but Consul couldn't resolve it for health checks.
**Solution:** Modified service registration to use "localhost" instead of hostname for both service address and health check URL.

### 2. Port Conflicts

**Problem:** Services failed to start due to ports already in use.
**Solution:** Implemented proper process cleanup using `lsof` and `kill` commands before restarting services.

### 3. Service Discovery Timing

**Problem:** API Gateway couldn't find services immediately after startup.
**Solution:** Ensured services had time to register and pass health checks before testing.

## Running the Application

### Prerequisites

- Go 1.18+
- Consul installed locally

### Startup Sequence

1. **Start Consul:**

```bash
consul agent -dev
```

2. **Start Users Service:**

```bash
cd services/users-service
go run .
```

3. **Start Products Service:**

```bash
cd services/products-service
go run .
```

4. **Start API Gateway:**

```bash
cd api-gateway
go run .
```

### Testing

1. **Test Users Service:**

```bash
curl http://localhost:8080/api/users/123
# Expected: Response from 'users-service': Details for user 123
```

2. **Test Products Service:**

```bash
curl http://localhost:8080/api/products/abc
# Expected: Response from 'products-service': Details for product abc
```

### Consul UI Access

- **URL:** http://localhost:8500
- **Services:** Should show both users-service and products-service as healthy

## Key Features Demonstrated

1. **Service Registration:** Automatic registration with Consul on startup
2. **Health Monitoring:** Consul performs periodic health checks
3. **Service Discovery:** API Gateway dynamically discovers service locations
4. **Request Routing:** Intelligent routing based on URL patterns
5. **Resilience:** System continues to work if individual services restart

## Learning Outcomes Achieved

- **Learning Outcome 2:** Implemented microservices using HTTP and JSON for inter-service communication
- **Learning Outcome 8:** Demonstrated observability through Consul's service monitoring and health checking

## Success Metrics

✅ Both services successfully register with Consul  
✅ Health checks pass for all services  
✅ API Gateway correctly routes requests  
✅ End-to-end request flow works seamlessly  
✅ Services can be restarted without system reconfiguration

## Screenshots

### 1. Consul UI showing registered services

_Screenshot should show both users-service and products-service with healthy status_

### 2. API Gateway terminal logs

_Screenshot should show gateway receiving requests and routing them successfully_

### 3. cURL/Postman requests and responses

_Screenshots showing successful API calls through the gateway_

## Conclusion

This implementation successfully demonstrates a microservices architecture with service discovery. The system allows for dynamic service registration, health monitoring, and request routing without hardcoded service locations. This pattern enables easy scaling, maintenance, and deployment of individual services in a distributed system.
