# freeplan-keeper
🛡️ Your free-tier service guardian — keeps services like RabbitMQ (CloudAMQP), PostgreSQL (Supabase), Redis (Upstash), and MongoDB Atlas active by sending regular heartbeat signals to prevent suspension or deletion due to inactivity.

## Features

- Automated heartbeat signals to keep your free-tier services active
- Currently supports:
  - MongoDB Atlas
  - Supabase
  - More services coming soon...
- Configurable check intervals
- Graceful shutdown handling

## Setup

1. Clone the repository
```bash
git clone https://github.com/xudongzhaodev/freeplan-keeper.git
cd freeplan-keeper
```

2. Configure your services
Copy `config.yaml.example` to `config.yaml` and update with your service credentials:
```yaml
mongodb:
  uri: your-mongodb-connection-string
  database: your-database-name

supabase:
  url: your-supabase-project-url
  api_key: your-supabase-api-key

check_interval: 30  # minutes
```

3. Build and run
```bash
go build -o keeper ./cmd/keeper
./keeper
```

## Configuration

- `check_interval`: Time between checks in minutes (default: 60)
- MongoDB configuration:
  - `uri`: MongoDB connection string
  - `database`: Database name to ping
- Supabase configuration:
  - `url`: Your Supabase project URL
  - `api_key`: Your Supabase API key

## License

MIT License - see LICENSE file for details
