# CMDB Lite Maintenance Guide

This document provides basic guidance on maintaining CMDB Lite in production environments.

## Maintenance Overview

Regular maintenance is essential for ensuring the stability, security, and performance of CMDB Lite. This guide covers the routine maintenance tasks and update procedures necessary to keep the system running smoothly.

### Maintenance Objectives

1. **System Stability**: Ensure the system remains stable and available
2. **Performance Optimization**: Maintain optimal performance levels
3. **Security Compliance**: Keep the system secure and compliant with policies
4. **Data Integrity**: Ensure data consistency and integrity

### Maintenance Schedule

| Task | Frequency | Duration |
|------|-----------|----------|
| Database Maintenance | Weekly | 1-2 hours |
| Application Maintenance | Monthly | 2-4 hours |
| System Updates | Quarterly | 4-8 hours |
| Security Updates | As needed | 1-2 hours |

## Routine Maintenance Tasks

### Database Maintenance

#### Database Vacuum and Analyze

Run this command weekly to optimize database performance:

```bash
# Connect to PostgreSQL and run vacuum and analyze
sudo -u postgres psql -d cmdb_lite -c "VACUUM ANALYZE;"
```

#### Database Backup Verification

Verify that your automated backups are working correctly:

```bash
# Check recent backups
ls -la /backups/cmdb-lite/database/ | tail -10
```

### Application Maintenance

#### Log Rotation

Check and rotate application logs to prevent disk space issues:

```bash
# Check log sizes
du -h /var/log/cmdb-lite/

# Rotate logs if needed
sudo logrotate -f /etc/logrotate.d/cmdb-lite
```

#### Application Restart

Restart the application services monthly to clear any memory leaks:

```bash
# Restart CMDB Lite services
sudo systemctl restart cmdb-lite-backend
sudo systemctl restart cmdb-lite-frontend
```

### System Maintenance

#### System Updates

Apply operating system updates quarterly:

```bash
# Update system packages
sudo apt update
sudo apt upgrade -y
```

#### Disk Space Check

Monitor disk space usage:

```bash
# Check disk usage
df -h
```

## Update Procedures

### Application Updates

1. **Backup the current version**:
   ```bash
   # Backup current application
   sudo cp -r /opt/cmdb-lite /opt/cmdb-lite-backup-$(date +%Y%m%d)
   ```

2. **Download and install the new version**:
   ```bash
   # Download new version
   wget https://github.com/your-org/cmdb-lite/releases/latest/download/cmdb-lite.tar.gz
   
   # Extract and install
   sudo tar -xzf cmdb-lite.tar.gz -C /opt/cmdb-lite
   ```

3. **Restart services**:
   ```bash
   sudo systemctl restart cmdb-lite-backend
   sudo systemctl restart cmdb-lite-frontend
   ```

### Database Updates

1. **Backup the database**:
   ```bash
   # Create database backup
   sudo -u postgres pg_dump cmdb_lite > /backups/cmdb-lite/database/cmdb-lite-pre-update-$(date +%Y%m%d).sql
   ```

2. **Apply database migrations**:
   ```bash
   # Run database migrations
   cd /opt/cmdb-lite/database
   sudo ./migrate.sh
   ```

### System Updates

1. **Schedule maintenance window**:
   - Notify users of upcoming maintenance
   - Schedule during low-usage periods

2. **Apply updates**:
   ```bash
   # Update system packages
   sudo apt update
   sudo apt upgrade -y
   ```

3. **Reboot if necessary**:
   ```bash
   # Check if reboot is required
   if [ -f /var/run/reboot-required ]; then
       sudo reboot
   fi
   ```

## Performance Tuning

### Database Performance

#### Monitor Slow Queries

Check for slow queries:

```bash
# Enable slow query logging
sudo -u postgres psql -d cmdb_lite -c "ALTER SYSTEM SET log_min_duration_statement = '1000';"
sudo systemctl reload postgresql

# Check slow query log
sudo tail -f /var/log/postgresql/postgresql-13-main.log | grep "duration:"
```

#### Optimize Database Configuration

Review and adjust PostgreSQL configuration:

```bash
# Edit PostgreSQL configuration
sudo nano /etc/postgresql/13/main/postgresql.conf

# Common settings to adjust:
# shared_buffers = 256MB
# effective_cache_size = 1GB
# work_mem = 4MB
# maintenance_work_mem = 64MB

# Restart PostgreSQL after changes
sudo systemctl restart postgresql
```

### Application Performance

#### Monitor Application Logs

Check for errors or performance issues:

```bash
# Monitor application logs
sudo tail -f /var/log/cmdb-lite/backend.log
sudo tail -f /var/log/cmdb-lite/frontend.log
```

#### Adjust Application Configuration

Review and adjust application settings:

```bash
# Edit environment file
sudo nano /opt/cmdb-lite/.env

# Common settings to adjust:
# LOG_LEVEL=info
# DB_POOL_SIZE=10
# CACHE_TTL=300
```

## Security Maintenance

### Vulnerability Management

#### Scan for Vulnerabilities

Regularly scan for security vulnerabilities:

```bash
# Install and run Lynis
sudo apt install lynis -y
sudo lynis audit system

# Review the report
sudo cat /var/log/lynis-report.dat
```

### Security Updates

#### Apply Security Patches

Apply security patches as soon as they're available:

```bash
# Install security updates only
sudo apt update
sudo unattended-upgrade -d
```

### Security Audits

#### Review User Access

Review user access permissions:

```bash
# List database users
sudo -u postgres psql -d cmdb_lite -c "\du"

# List system users with sudo access
sudo grep -Po '^sudo.+:\K.*$' /etc/group
```

#### Review Firewall Rules

Review firewall configuration:

```bash
# Check firewall status
sudo ufw status

# Review firewall rules
sudo ufw show added
```

## Troubleshooting Common Issues

### Database Issues

#### Database Connection Issues

If you can't connect to the database:

1. Check if PostgreSQL is running:
   ```bash
   sudo systemctl status postgresql
   ```

2. Check database logs:
   ```bash
   sudo tail -f /var/log/postgresql/postgresql-13-main.log
   ```

3. Verify connection settings:
   ```bash
   sudo nano /opt/cmdb-lite/.env
   ```

#### Database Performance Issues

If the database is slow:

1. Check for long-running queries:
   ```bash
   sudo -u postgres psql -d cmdb_lite -c "SELECT pid, now() - pg_stat_activity.query_start AS duration, query FROM pg_stat_activity WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';"
   ```

2. Check database size:
   ```bash
   sudo -u postgres psql -d cmdb_lite -c "SELECT pg_size_pretty(pg_database_size('cmdb_lite'));"
   ```

### Application Issues

#### Application Won't Start

If the application won't start:

1. Check service status:
   ```bash
   sudo systemctl status cmdb-lite-backend
   sudo systemctl status cmdb-lite-frontend
   ```

2. Check application logs:
   ```bash
   sudo journalctl -u cmdb-lite-backend -n 100
   sudo journalctl -u cmdb-lite-frontend -n 100
   ```

3. Check configuration:
   ```bash
   sudo nano /opt/cmdb-lite/.env
   ```

#### Application Performance Issues

If the application is slow:

1. Check system resources:
   ```bash
   top
   free -h
   df -h
   ```

2. Check application logs for errors:
   ```bash
   sudo tail -f /var/log/cmdb-lite/backend.log | grep ERROR
   sudo tail -f /var/log/cmdb-lite/frontend.log | grep ERROR
   ```

### System Issues

#### Disk Space Issues

If you're running out of disk space:

1. Check disk usage:
   ```bash
   df -h
   ```

2. Find large files:
   ```bash
   sudo find / -type f -size +100M -exec ls -lh {} \;
   ```

3. Clean up old logs:
   ```bash
   sudo find /var/log -type f -name "*.log.*" -delete
   ```

#### Memory Issues

If you're running out of memory:

1. Check memory usage:
   ```bash
   free -h
   ```

2. Find processes using the most memory:
   ```bash
   ps aux --sort=-%mem | head
   ```

3. Check for memory leaks:
   ```bash
   sudo journalctl -u cmdb-lite-backend | grep -i "out of memory"
   ```

## Maintenance Scheduling

### Maintenance Windows

Schedule maintenance during low-usage periods:

| Day | Time | Impact |
|-----|------|--------|
| Sunday | 02:00 - 06:00 | Low |
| Wednesday | 22:00 - 00:00 | Medium |

### Change Management

Follow this process for all changes:

1. **Plan**: Document the change, expected impact, and rollback plan
2. **Approve**: Get approval from stakeholders
3. **Communicate**: Notify users of the upcoming maintenance
4. **Execute**: Perform the maintenance during the scheduled window
5. **Verify**: Confirm the system is working correctly after the change
6. **Document**: Update documentation and close the change request

### Communication Procedures

Notify users of upcoming maintenance:

1. **Email Notification**: Send email at least 48 hours before scheduled maintenance
2. **Banner Notification**: Display a banner in the application 24 hours before maintenance
3. **Status Page**: Update the status page with maintenance information

For more information on operating CMDB Lite, see the [Operations Guide](README.md) and the [Monitoring and Troubleshooting Guide](monitoring.md).