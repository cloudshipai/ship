#!/bin/bash
# User data script for web server setup

# Update system
yum update -y

# Install Apache web server
yum install -y httpd

# Create a simple index page
cat > /var/www/html/index.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>${project_name} - Ship CLI Example</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        h1 { color: #333; }
        .info { background: #f0f0f0; padding: 20px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to ${project_name}</h1>
        <div class="info">
            <h2>Ship CLI Example Web Application</h2>
            <p>This is a medium complexity Terraform example demonstrating:</p>
            <ul>
                <li>VPC with public and private subnets</li>
                <li>Application Load Balancer</li>
                <li>Auto Scaling Group</li>
                <li>Security Groups</li>
                <li>NAT Gateways for outbound traffic</li>
            </ul>
            <p>Instance ID: $(ec2-metadata --instance-id | cut -d " " -f 2)</p>
            <p>Availability Zone: $(ec2-metadata --availability-zone | cut -d " " -f 2)</p>
        </div>
    </div>
</body>
</html>
EOF

# Start and enable Apache
systemctl start httpd
systemctl enable httpd