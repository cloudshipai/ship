# Example Dockerfile with known vulnerabilities for testing
FROM ubuntu:18.04

# Install outdated packages with known vulnerabilities
RUN apt-get update && apt-get install -y \
    curl=7.58.0-2ubuntu3.8 \
    openssl=1.1.1-1ubuntu2.1~18.04.5 \
    nginx=1.14.0-0ubuntu1.7

# Add a user with weak configuration
RUN useradd -m -s /bin/bash testuser
USER testuser

# Copy application files
COPY app.py /app/
WORKDIR /app

EXPOSE 8080
CMD ["python3", "app.py"]