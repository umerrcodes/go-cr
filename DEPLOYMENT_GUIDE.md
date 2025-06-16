# Deployment Guide for Vercel

## Prerequisites

1. **Vercel Account**: Sign up at [vercel.com](https://vercel.com)
2. **Vercel CLI**: Install globally
   ```bash
   npm install -g vercel
   ```
3. **Neon Database**: Your database is already set up!

## Step-by-Step Deployment

### 1. Login to Vercel
```bash
vercel login
```

### 2. Initialize Project
```bash
vercel
```
Follow the prompts:
- Set up and deploy? **Y**
- Which scope? **Your personal account**
- Link to existing project? **N**
- What's your project's name? **dummy-backend** (or any name)
- In which directory is your code located? **./** 

### 3. Set Environment Variables

After deployment, set these environment variables in your Vercel dashboard:

```env
DB_DSN=postgresql://neondb_owner:npg_jgLarnDwv19q@ep-twilight-shadow-a1yxhoxp-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
JWT_SECRET=your-super-secret-production-key-change-this
GIN_MODE=release
```

**Important:** Change the JWT_SECRET to a secure random string in production!

### 4. Deploy to Production
```bash
vercel --prod
```

## Environment Variables Setup via Dashboard

1. Go to [vercel.com/dashboard](https://vercel.com/dashboard)
2. Click on your project
3. Go to **Settings** → **Environment Variables**
4. Add each variable:
   - **Name**: `DB_DSN`
   - **Value**: `postgresql://neondb_owner:npg_jgLarnDwv19q@ep-twilight-shadow-a1yxhoxp-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require`
   - **Environment**: All (Production, Preview, Development)

Repeat for `JWT_SECRET` and `GIN_MODE`.

## Alternative: Environment Variables via CLI

```bash
# Set environment variables via CLI
vercel env add DB_DSN
# Paste your database URL when prompted

vercel env add JWT_SECRET
# Enter a secure JWT secret

vercel env add GIN_MODE
# Enter: release
```

## Testing Your Deployed API

Once deployed, you'll get a URL like: `https://dummy-backend-xyz.vercel.app`

### Test Health Endpoint
```bash
curl https://your-app.vercel.app/health
```

### Test Registration
```bash
curl -X POST https://your-app.vercel.app/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Test Tasks (with token from registration)
```bash
curl -X GET https://your-app.vercel.app/api/tasks \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## API Endpoints

Your deployed API will have these endpoints:

- `GET /health` - Health check
- `POST /api/auth/register` - User registration  
- `POST /api/auth/login` - User login
- `GET /api/tasks` - Get all tasks (auth required)
- `POST /api/tasks` - Create task (auth required)
- `GET /api/tasks/:id` - Get specific task (auth required)
- `PUT /api/tasks/:id` - Update task (auth required)
- `DELETE /api/tasks/:id` - Delete task (auth required)

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify your `DB_DSN` is correct
   - Check Neon database is active
   - Ensure SSL mode is required

2. **Import Errors**
   - Vercel automatically handles Go modules
   - Make sure `go.mod` is in the root directory

3. **JWT Token Issues**
   - Ensure `JWT_SECRET` is set
   - Token expires in 24 hours by default

4. **CORS Issues**
   - CORS middleware is already configured
   - Allows all origins (`*`) for testing

### Checking Logs

1. Go to your Vercel dashboard
2. Click on your project
3. Go to **Functions** tab
4. Click on any function to see logs

## Production Considerations

1. **Security**
   - Change JWT_SECRET to a strong random string
   - Consider implementing rate limiting
   - Add input validation

2. **Database**
   - Monitor Neon database usage
   - Consider connection pooling for high traffic

3. **Monitoring**
   - Set up Vercel Analytics
   - Monitor function execution times
   - Track error rates

## Next Steps

1. **Custom Domain**: Add your domain in Vercel settings
2. **CI/CD**: Connect to GitHub for automatic deployments
3. **Monitoring**: Set up alerts for downtime
4. **Scaling**: Monitor and optimize based on usage

Your API is now live and ready to use! 🚀 