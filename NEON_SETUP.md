# Neon PostgreSQL Setup Guide 🐘

## 🚀 **Step 1: Create Neon Database**

1. **Sign up for Neon:**
   - Go to [neon.tech](https://neon.tech)
   - Sign up with GitHub
   - Verify your email

2. **Create a New Project:**
   - Click "Create Project"
   - Choose a project name (e.g., "task-api")
   - Select region (choose closest to your users)
   - Click "Create Project"

3. **Get Connection Details:**
   After creation, you'll see connection details:
   ```
   Host: ep-xxx-xxx.us-east-2.aws.neon.tech
   Database: neondb
   Username: your-username
   Password: your-password
   ```

4. **Copy Connection String:**
   Neon provides a ready-to-use connection string:
   ```
   postgresql://username:password@ep-xxx-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require
   ```

## 🛠️ **Step 2: Create Database Schema**

1. **Access Neon SQL Editor:**
   - In your Neon dashboard, click "SQL Editor"
   - Or use any PostgreSQL client with the connection string

2. **Run Schema Creation:**
   ```sql
   -- Create tasks table
   CREATE TABLE IF NOT EXISTS tasks (
       id SERIAL PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       description TEXT,
       completed BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
   );

   -- Insert sample data
   INSERT INTO tasks (title, description, completed) VALUES 
       ('Learn Go', 'Study Go programming language', FALSE),
       ('Build API', 'Create a REST API with Go', TRUE),
       ('Deploy to Vercel', 'Deploy Go API as serverless functions', FALSE),
       ('Use Neon Database', 'Integrate with Neon PostgreSQL', FALSE);

   -- Verify data
   SELECT * FROM tasks;
   ```

## 🔗 **Step 3: Connection Details**

Your Neon connection string will look like:
```
postgresql://username:password@ep-xxx-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require
```

**Important Notes:**
- ✅ Always use `sslmode=require` for Neon
- ✅ Connection string includes all details needed
- ✅ Neon handles connection pooling automatically
- ✅ Free tier includes 512MB storage and 1 million queries/month

## 🌐 **Step 4: Environment Variables**

For your deployment, you'll need:
```bash
DATABASE_URL=postgresql://username:password@ep-xxx-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require
```

That's it! Neon provides everything in one connection string. 
postgresql://neondb_owner:npg_jgLarnDwv19q@ep-twilight-shadow-a1yxhoxp-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require