# freeplan-keeper
🛡️ Your free-tier service guardian — keeps services like RabbitMQ (CloudAMQP), PostgreSQL (Supabase), Redis (Upstash), and MongoDB Atlas active by sending regular heartbeat signals to prevent suspension or deletion due to inactivity.

## Features

- Automated heartbeat signals to keep your free-tier services active
- Currently supports:
  - MongoDB Atlas
  - Supabase
  - More services coming soon...
- Service-specific enable/disable switches
- Multiple deployment options (local, cron, GitHub Actions)

## Setup

### Local Setup

1. Clone the repository
```bash
git clone https://github.com/xudongzhaodev/freeplan-keeper.git
cd freeplan-keeper
```

2. Configure your services
Copy `config.yaml.example` to `config.yaml` and update with your service credentials:
```yaml
mongodb:
  enabled: true
  uri: your-mongodb-connection-string
  database: your-database-name

supabase:
  enabled: true
  url: your-supabase-project-url
  api_key: your-supabase-api-key
```

3. Build and run
```bash
go build -o keeper ./cmd/keeper
./keeper
```

### GitHub Actions Setup

You can run this keeper automatically using GitHub Actions. This is the recommended way as it:
- Runs in the cloud (no need for a local server)
- Keeps sensitive information secure using GitHub Secrets
- Provides automatic scheduling and execution logs

To set up GitHub Actions:

1. Fork this repository

2. Set up GitHub Secrets
   Go to your repository's Settings -> Secrets and variables -> Actions and add the following secrets:

   - `MONGODB_ENABLED`: Set to 'true' or 'false'
   - `MONGODB_URI`: Your MongoDB connection string
   - `MONGODB_DATABASE`: Your database name
   - `SUPABASE_ENABLED`: Set to 'true' or 'false'
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_API_KEY`: Your Supabase API key

3. The workflow will automatically:
   - Run every hour
   - Run on push to main branch
   - Allow manual triggers via GitHub Actions UI

You can monitor the execution in the Actions tab of your repository.

### Alternative Deployment Options

#### Linux Cron
```bash
# Run every hour
0 * * * * /path/to/keeper
```

#### Windows Task Scheduler
```batch
# Create hourly task
schtasks /create /tn "ServiceKeeper" /tr "path\to\keeper.exe" /sc hourly
```

## Configuration

- MongoDB configuration:
  - `enabled`: Whether to activate MongoDB checks
  - `uri`: MongoDB connection string
  - `database`: Database name to ping

- Supabase configuration:
  - `enabled`: Whether to activate Supabase checks
  - `url`: Your Supabase project URL
  - `api_key`: Your Supabase API key

## Security Notes

- Never commit your `config.yaml` file (it's in .gitignore)
- When using GitHub Actions, always use Secrets for sensitive information
- The program will clean up the config file after each GitHub Actions run

## License

MIT License - see LICENSE file for details
