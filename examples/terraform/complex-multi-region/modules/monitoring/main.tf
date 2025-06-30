# SNS Topic for alerts
resource "aws_sns_topic" "alerts" {
  name = "${var.project_name}-alerts"

  tags = {
    Name        = "${var.project_name}-alerts"
    Environment = var.environment
  }
}

resource "aws_sns_topic_subscription" "email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# CloudWatch Dashboard
resource "aws_cloudwatch_dashboard" "main" {
  dashboard_name = "${var.project_name}-${var.environment}"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/ApplicationELB", "TargetResponseTime", "LoadBalancer", split("/", var.monitored_resources.primary_alb_arn)[2]],
            ["AWS/ApplicationELB", "TargetResponseTime", "LoadBalancer", split("/", var.monitored_resources.secondary_alb_arn)[2]]
          ]
          period = 300
          stat   = "Average"
          region = "us-east-1"
          title  = "ALB Response Time"
        }
      },
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/EC2", "CPUUtilization", "AutoScalingGroupName", var.monitored_resources.primary_asg_name],
            ["AWS/EC2", "CPUUtilization", "AutoScalingGroupName", var.monitored_resources.secondary_asg_name]
          ]
          period = 300
          stat   = "Average"
          region = "us-east-1"
          title  = "ASG CPU Utilization"
        }
      },
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", var.monitored_resources.primary_db_id],
            ["AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", var.monitored_resources.secondary_db_id]
          ]
          period = 300
          stat   = "Average"
          region = "us-east-1"
          title  = "RDS CPU Utilization"
        }
      },
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/RDS", "DatabaseConnections", "DBInstanceIdentifier", var.monitored_resources.primary_db_id],
            ["AWS/RDS", "DatabaseConnections", "DBInstanceIdentifier", var.monitored_resources.secondary_db_id]
          ]
          period = 300
          stat   = "Average"
          region = "us-east-1"
          title  = "Database Connections"
        }
      }
    ]
  })
}

# CloudWatch Log Groups for Application Logs
resource "aws_cloudwatch_log_group" "app_logs" {
  name              = "/aws/ec2/${var.project_name}"
  retention_in_days = 30

  tags = {
    Name        = "${var.project_name}-app-logs"
    Environment = var.environment
  }
}

# CloudWatch Alarms
resource "aws_cloudwatch_metric_alarm" "alb_response_time" {
  alarm_name          = "${var.project_name}-high-response-time"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "TargetResponseTime"
  namespace           = "AWS/ApplicationELB"
  period              = "300"
  statistic           = "Average"
  threshold           = "1"
  alarm_description   = "This metric monitors ALB response time"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    LoadBalancer = split("/", var.monitored_resources.primary_alb_arn)[2]
  }
}

resource "aws_cloudwatch_metric_alarm" "alb_unhealthy_hosts" {
  alarm_name          = "${var.project_name}-unhealthy-hosts"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "UnHealthyHostCount"
  namespace           = "AWS/ApplicationELB"
  period              = "300"
  statistic           = "Average"
  threshold           = "0"
  alarm_description   = "This metric monitors unhealthy hosts"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    LoadBalancer = split("/", var.monitored_resources.primary_alb_arn)[2]
  }
}

# Cost anomaly detection
resource "aws_ce_anomaly_monitor" "main" {
  name              = "${var.project_name}-cost-monitor"
  monitor_type      = "DIMENSIONAL"
  monitor_dimension = "SERVICE"
}

resource "aws_ce_anomaly_subscription" "main" {
  name      = "${var.project_name}-cost-alerts"
  threshold_expression {
    dimension {
      key           = "ANOMALY_TOTAL_IMPACT_ABSOLUTE"
      values        = ["100"]
      match_options = ["GREATER_THAN_OR_EQUAL"]
    }
  }

  frequency = "DAILY"

  monitor_arn_list = [
    aws_ce_anomaly_monitor.main.arn
  ]

  subscriber {
    type    = "EMAIL"
    address = var.alert_email
  }
}