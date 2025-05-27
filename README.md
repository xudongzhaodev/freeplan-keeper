# Freeplan Keeper

A service keeper application designed to prevent free-tier services from being suspended due to inactivity. Currently supports:
- MongoDB Atlas
- Supabase

## Features

- Configurable ping intervals
- Record cleanup with configurable limits
- Support for multiple services
- Detailed ping records with timestamp and source tracking
- Global instance identification

## Installation

1. Clone the repository:
```bash
git clone https://github.com/xudongzhaodev/freeplan-keeper.git
cd freeplan-keeper
```

2. Copy the example configuration:
```bash
cp config.example.yaml config.yaml
```

3. Edit `config.yaml` with your service credentials

4. Build the application:
```bash
go build -o keeper ./cmd/keeper
```

## Configuration

The configuration file (`config.yaml`) supports the following options:

```yaml
# Global identifier for this keeper instance
hostname: "local-keeper"

mongodb:
  enabled: false  # Set to true to enable MongoDB keeper
  uri: "mongodb+srv://<username>:<password>@<cluster>.mongodb.net/?retryWrites=true&w=majority"
  database: "keeper"  # The database to store keep-alive records
  keep_records_limit: 100  # Number of records to keep in keep_alive_reserved collection

supabase:
  enabled: true  # Set to true to enable Supabase keeper
  uri: "postgresql://<user>:<password>@<host>:<port>/<dbname>"  # Get this from Supabase connection string
  keep_records_limit: 100  # Number of records to keep in keep_alive_reserved table
```

### MongoDB Configuration

1. Get your MongoDB Atlas connection string from the Atlas dashboard
2. Replace `<username>`, `<password>`, and `<cluster>` with your credentials
3. Set `enabled: true` to activate the MongoDB keeper

### Supabase Configuration

1. Get your connection string from Supabase dashboard:
   - Go to Project Settings > Database
   - Find the Connection String section
   - Copy the Connection string and replace the password
2. Set `enabled: true` to activate the Supabase keeper

## Usage

Run the keeper:
```bash
./keeper
```

The keeper will:
1. Connect to enabled services
2. Create necessary tables/collections if they don't exist
3. Periodically ping each service
4. Maintain a record of pings with cleanup

## Record Structure

### MongoDB Collection: keep_alive_reserved
```json
{
  "ping_timestamp": "2024-03-19T10:00:00Z",
  "ping_source": "mongodb-keeper",
  "ping_details": {
    "hostname": "local-keeper",
    "version": "1.0"
  }
}
```

### Supabase Table: keep_alive_reserved
```sql
CREATE TABLE keep_alive_reserved (
  id BIGSERIAL PRIMARY KEY,
  ping_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ping_source VARCHAR(255),
  ping_details JSONB DEFAULT '{}'::jsonb
);
```

## License

MIT License
