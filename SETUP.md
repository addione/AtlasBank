# Setup Guide for AtlasBank

## Handling Corporate VPN / Certificate Issues

If you're working behind a corporate VPN or proxy and encountering certificate verification errors when running `go mod tidy`, you have several options:

### Option 1: Use Direct Proxy (Recommended)
This bypasses the Go module proxy and downloads directly from source repositories:

```bash
GOPROXY=direct GOSUMDB=off go mod tidy
```

### Option 2: Set Environment Variables Permanently
Add these to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
export GOPROXY=direct
export GOSUMDB=off
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

### Option 3: Configure Git to Use HTTP Instead of HTTPS
For specific repositories:

```bash
git config --global url."http://github.com/".insteadOf "https://github.com/"
```

### Option 4: Install Corporate CA Certificates
If your company provides CA certificates, install them:

```bash
# Copy your corporate CA certificate to the system trust store
sudo cp /path/to/corporate-ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

### Option 5: Use GOINSECURE (Not Recommended for Production)
This disables certificate verification for specific modules:

```bash
GOINSECURE=* go mod tidy
```

**Warning**: This is insecure and should only be used in development environments.

## Building with Docker

The good news is that once you have the dependencies downloaded, Docker will handle the build process inside the container, which may not have the same certificate issues:

```bash
# This will build inside Docker, bypassing local Go issues
docker-compose build
```

## Alternative: Pre-download Dependencies

If you continue to have issues, you can:

1. Download dependencies on a machine without VPN restrictions
2. Commit the `vendor/` directory to your repository
3. Use vendored dependencies:

```bash
go mod vendor
go build -mod=vendor ./cmd/api
```

## Quick Start (After Dependencies are Downloaded)

Once `go mod tidy` completes successfully:

```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f app

# Test the application
curl http://localhost:8081/health
```

## Troubleshooting

### Still Getting Certificate Errors?

1. Check your proxy settings:
```bash
echo $HTTP_PROXY
echo $HTTPS_PROXY
```

2. Try disabling proxy temporarily:
```bash
unset HTTP_PROXY
unset HTTPS_PROXY
unset http_proxy
unset https_proxy
```

3. Use the Makefile with environment variables:
```bash
GOPROXY=direct GOSUMDB=off make tidy
```

### Docker Build Fails?

If Docker build fails due to Go module issues, you can modify the Dockerfile to use the same workaround:

```dockerfile
# In Dockerfile, before RUN go mod download
ENV GOPROXY=direct
ENV GOSUMDB=off
```

## Need Help?

If you continue to experience issues:
1. Contact your IT department for corporate CA certificates
2. Ask for proxy configuration details
3. Consider using a development VM outside the corporate network
