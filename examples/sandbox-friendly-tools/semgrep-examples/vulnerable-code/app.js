// JavaScript application with security vulnerabilities for Semgrep testing

const express = require('express');
const mysql = require('mysql');
const fs = require('fs');
const crypto = require('crypto');

const app = express();
app.use(express.json());

// SQL Injection vulnerability
app.post('/login', (req, res) => {
    const { username, password } = req.body;
    
    // Vulnerable: Direct string concatenation in SQL query
    const query = `SELECT * FROM users WHERE username = '${username}' AND password = '${password}'`;
    
    const connection = mysql.createConnection({
        host: 'localhost',
        user: 'root',
        password: 'password',
        database: 'myapp'
    });
    
    connection.query(query, (error, results) => {
        if (error) throw error;
        res.json(results);
    });
});

// Command Injection vulnerability
app.post('/file-info', (req, res) => {
    const { filename } = req.body;
    
    // Vulnerable: Direct command execution with user input
    const { exec } = require('child_process');
    exec(`file ${filename}`, (error, stdout, stderr) => {
        if (error) {
            res.status(500).json({ error: error.message });
            return;
        }
        res.json({ output: stdout });
    });
});

// Path Traversal vulnerability
app.get('/file/:filename', (req, res) => {
    const filename = req.params.filename;
    
    // Vulnerable: No path validation
    const filePath = `./uploads/${filename}`;
    
    fs.readFile(filePath, 'utf8', (err, data) => {
        if (err) {
            res.status(404).json({ error: 'File not found' });
            return;
        }
        res.send(data);
    });
});

// Weak Cryptography
app.post('/encrypt', (req, res) => {
    const { data } = req.body;
    
    // Vulnerable: Using MD5 which is cryptographically broken
    const hash = crypto.createHash('md5').update(data).digest('hex');
    
    res.json({ hash });
});

// XSS vulnerability (reflected)
app.get('/search', (req, res) => {
    const query = req.query.q;
    
    // Vulnerable: Direct output without sanitization
    res.send(`<h1>Search Results for: ${query}</h1>`);
});

// Hardcoded secrets
const API_KEY = 'sk-1234567890abcdef1234567890abcdef';  // Vulnerable: Hardcoded API key
const DATABASE_PASSWORD = 'super_secret_password123';   // Vulnerable: Hardcoded password

// Insecure randomness
function generateToken() {
    // Vulnerable: Using Math.random() for security-sensitive operations
    return Math.random().toString(36).substr(2, 9);
}

// CORS misconfiguration
app.use((req, res, next) => {
    // Vulnerable: Overly permissive CORS
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Headers', '*');
    next();
});

app.listen(3000, () => {
    console.log('Server running on port 3000');
});