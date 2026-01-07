#!/bin/sh
set -e

# Generate runtime config for frontend
cat > /app/frontend/build/config.js << EOF
window.RUNTIME_CONFIG = {
  API_URL: '${DOMAIN:+https://$DOMAIN}'
};
EOF

# Execute the main command
exec "$@"
