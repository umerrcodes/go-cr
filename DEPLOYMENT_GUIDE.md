# Go Backend Deployment Guide 🚀

## 🤔 **Deployment Models Explained**

### **Traditional Server (Current Setup)**
```
✅ Your Local Development
- Always running process
- Persistent connections
- In-memory state
- Port-based routing
```

### **Serverless Functions (Vercel)**
```
🔄 Production Deployment
- Functions start on-demand
- No persistent state
- Stateless execution
- Event-driven routing
```

## 🔧 **Required Changes for Vercel**

### **1. Project Structure**
```
task-api/
├── api/
│   ├── health.go       # GET /api/health
│   ├── tasks.go        # GET /api/tasks
│   └── tasks/
│       └── [id].go     # GET/PUT/DELETE /api/tasks/[id]
├── go.mod
├── vercel.json
└── README.md
```

### **2. Code Changes**
- ❌ Remove `main()` function and server setup
- ❌ Remove global variables (database connections)
- ✅ Each endpoint becomes a separate function
- ✅ Database connection per request
- ✅ Use environment variables for config

### **3. Database Changes**
- ❌ Can't use local PostgreSQL
- ✅ Need cloud database (Supabase/PlanetScale)
- ✅ Connection string from environment variables

## 🗄️ **Database Options**

### **Supabase (Recommended)**
- PostgreSQL-compatible
- Free tier available
- Built-in API dashboard
- Easy connection strings

### **PlanetScale**
- MySQL-compatible
- Serverless database
- Branch-based development

### **Railway PostgreSQL**
- Traditional PostgreSQL
- Simple deployment
- Good for learning

## 🚀 **Step-by-Step Deployment**

### **Step 1: Set Up Cloud Database (Railway)**

1. **Sign up for Railway:**
   - Go to [railway.app](https://railway.app)
   - Sign up with GitHub
   - Verify your account

2. **Create PostgreSQL Database:**
   - Click "New Project"
   - Select "Provision PostgreSQL"
   - Wait for deployment (1-2 minutes)

3. **Get Database Connection:**
   - Click on your PostgreSQL service
   - Go to "Variables" tab
   - Copy the `DATABASE_URL` value
   - It looks like: `postgresql://postgres:password@containers-us-west-xxx.railway.app:5432/railway`

4. **Create Tables:**
   - Go to "Data" tab in Railway dashboard
   - Run this SQL:
   ```sql
   CREATE TABLE IF NOT EXISTS tasks (
       id SERIAL PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       description TEXT,
       completed BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
   );
   
   INSERT INTO tasks (title, description, completed) VALUES 
       ('Learn Go', 'Study Go programming language', FALSE),
       ('Build API', 'Create a REST API with Go', TRUE),
       ('Deploy to Vercel', 'Deploy Go API as serverless functions', FALSE);
   ```

### **Step 2: Deploy to Vercel**

1. **Install Vercel CLI:**
   ```bash
   npm install -g vercel
   ```

2. **Initialize Git (if not already):**
   ```bash
   git init
   git add .
   git commit -m "Initial serverless Go API"
   ```

3. **Deploy to Vercel:**
   ```bash
   vercel
   ```

4. **Configure Environment Variables:**
   - During deployment, Vercel will ask for environment variables
   - Add: `DATABASE_URL` = your Railway database URL
   - Or set it later in Vercel dashboard

5. **Test Your Deployed API:**
   ```bash
   # Replace YOUR_VERCEL_URL with your actual Vercel URL
   curl https://YOUR_VERCEL_URL.vercel.app/api/health
   curl https://YOUR_VERCEL_URL.vercel.app/api/tasks
   ```

### **Step 3: API Endpoints Structure**

Your deployed API will have these endpoints:

```
https://your-app.vercel.app/api/health          # GET - Health check
https://your-app.vercel.app/api/tasks           # GET - All tasks, POST - Create task
https://your-app.vercel.app/api/task?id=1       # GET/PUT/DELETE - Specific task
```

### **Step 4: Testing Commands**

```bash
# Health check
curl https://your-app.vercel.app/api/health

# Get all tasks
curl https://your-app.vercel.app/api/tasks

# Create new task
curl -X POST https://your-app.vercel.app/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Deployed Task", "description": "Created via Vercel API"}'

# Get specific task
curl https://your-app.vercel.app/api/task?id=1

# Update task
curl -X PUT https://your-app.vercel.app/api/task?id=1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Task", "description": "Updated via Vercel", "completed": true}'

# Delete task
curl -X DELETE https://your-app.vercel.app/api/task?id=1
```

## 🔧 **Key Differences: Local vs Serverless**

### **Local Development:**
- Single process, always running
- Global database connection pool
- In-memory state
- Port-based routing

### **Serverless (Vercel):**
- Function per request
- New database connection per function
- Stateless execution
- File-based routing

## 🛠️ **Troubleshooting**

### **Common Issues:**

1. **Database Connection Failed:**
   - Check `DATABASE_URL` environment variable
   - Ensure Railway database is running
   - Verify SSL settings (`sslmode=require` for cloud)

2. **Function Timeout:**
   - Serverless functions have time limits
   - Optimize database queries
   - Use connection pooling settings

3. **CORS Errors:**
   - Already handled in our functions
   - Check browser console for specific errors

### **Environment Variables:**
```bash
# Required for production
DATABASE_URL=postgresql://user:pass@host:port/db

# Optional for development
DB_HOST=localhost
DB_PORT=5432
DB_USER=umer
DB_NAME=taskapi
``` 