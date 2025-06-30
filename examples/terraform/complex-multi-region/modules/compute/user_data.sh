#!/bin/bash
# User data script for application setup

# Update system
yum update -y

# Install required packages
yum install -y httpd amazon-cloudwatch-agent jq

# Configure CloudWatch agent
cat > /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json << EOF
{
  "metrics": {
    "namespace": "${project_name}-${environment}",
    "metrics_collected": {
      "cpu": {
        "measurement": [
          {
            "name": "cpu_usage_idle",
            "rename": "CPU_USAGE_IDLE",
            "unit": "Percent"
          },
          {
            "name": "cpu_usage_iowait",
            "rename": "CPU_USAGE_IOWAIT",
            "unit": "Percent"
          },
          "cpu_time_guest"
        ],
        "totalcpu": false,
        "metrics_collection_interval": 60
      },
      "disk": {
        "measurement": [
          {
            "name": "used_percent",
            "rename": "DISK_USED_PERCENT",
            "unit": "Percent"
          }
        ],
        "metrics_collection_interval": 60,
        "resources": [
          "*"
        ]
      },
      "mem": {
        "measurement": [
          {
            "name": "mem_used_percent",
            "rename": "MEM_USED_PERCENT",
            "unit": "Percent"
          }
        ],
        "metrics_collection_interval": 60
      }
    }
  }
}
EOF

# Start CloudWatch agent
/opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
  -a fetch-config \
  -m ec2 \
  -s \
  -c file:/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json

# Create application
cat > /var/www/html/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>${project_name} - Complex Multi-Region App</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1000px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .info { background: #e7f3ff; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .metrics { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; margin: 20px 0; }
        .metric { background: #f8f9fa; padding: 15px; border-radius: 5px; text-align: center; }
        .metric h3 { margin: 0; color: #666; font-size: 14px; }
        .metric p { margin: 10px 0 0 0; font-size: 24px; font-weight: bold; color: #333; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Ship CLI - Complex Multi-Region Example</h1>
        <div class="info">
            <h2>Environment: ${environment}</h2>
            <p>This application demonstrates a complex, production-ready architecture:</p>
            <ul>
                <li>Multi-region deployment (Primary + DR)</li>
                <li>Auto Scaling Groups with CloudWatch metrics</li>
                <li>RDS Multi-AZ database</li>
                <li>S3 cross-region replication</li>
                <li>VPC with public, private, and database subnets</li>
                <li>Comprehensive monitoring and alerting</li>
            </ul>
        </div>
        
        <div class="metrics">
            <div class="metric">
                <h3>Instance ID</h3>
                <p id="instance-id">Loading...</p>
            </div>
            <div class="metric">
                <h3>Availability Zone</h3>
                <p id="az">Loading...</p>
            </div>
            <div class="metric">
                <h3>Region</h3>
                <p id="region">Loading...</p>
            </div>
        </div>

        <div class="info">
            <h3>Database Connection</h3>
            <p>Endpoint: ${database_endpoint}</p>
            <p id="db-status">Status: Checking...</p>
        </div>
    </div>

    <script>
        // Fetch instance metadata
        fetch('http://169.254.169.254/latest/meta-data/instance-id')
            .then(response => response.text())
            .then(data => document.getElementById('instance-id').textContent = data)
            .catch(() => document.getElementById('instance-id').textContent = 'N/A');

        fetch('http://169.254.169.254/latest/meta-data/placement/availability-zone')
            .then(response => response.text())
            .then(data => {
                document.getElementById('az').textContent = data;
                document.getElementById('region').textContent = data.slice(0, -1);
            })
            .catch(() => {
                document.getElementById('az').textContent = 'N/A';
                document.getElementById('region').textContent = 'N/A';
            });

        // Simulate database check
        setTimeout(() => {
            document.getElementById('db-status').textContent = 'Status: Connected';
        }, 1000);
    </script>
</body>
</html>
EOF

# Create health check endpoint
cat > /var/www/html/health << 'EOF'
OK
EOF

# Start and enable Apache
systemctl start httpd
systemctl enable httpd

# Log the completion
echo "User data script completed at $(date)" >> /var/log/user-data.log